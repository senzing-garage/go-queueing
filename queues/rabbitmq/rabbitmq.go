package rabbitmq

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-queueing/queues"
)

type ClientRabbitMQ struct {
	ExchangeName string
	QueueName    string
	// desired / default delay durations
	ReconnectDelay time.Duration
	ReInitDelay    time.Duration
	ResendDelay    time.Duration
	RoutingKey     string

	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	isReady         bool
	logger          logging.Logging
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	notifyReady     chan interface{}
	// current delay durations
	reconnectDelay time.Duration
	reInitDelay    time.Duration
	resendDelay    time.Duration
}

const (
	OptionCallerSkip      = 4
	ReconnectDelaySeconds = 2
	ReInitDelaySeconds    = 2
	ResendDelaySeconds    = 1
)

// ----------------------------------------------------------------------------
// Constructors
// ----------------------------------------------------------------------------

// New creates a single RabbitMQ client that will automatically
// attempt to connect to the server.  Reconnection delays are set to defaults.
func NewClient(urlString string) (*ClientRabbitMQ, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, wraperror.Errorf(err, "unable to parse RabbitMQ URL string")
	}

	queryMap, _ := url.ParseQuery(u.RawQuery)
	if len(queryMap["exchange"]) < 1 || len(queryMap["queue-name"]) < 1 {
		return nil, wraperror.Errorf(errForPackage, "please define an exchange and queue-name as query parameters")
	}

	routingKey := queryMap["queue-name"][0]
	if len(queryMap["routing-key"]) > 0 {
		routingKey = queryMap["routing-key"][0]
	}

	client := ClientRabbitMQ{
		ExchangeName:   queryMap["exchange"][0],
		QueueName:      queryMap["queue-name"][0],
		ReconnectDelay: ReconnectDelaySeconds * time.Second,
		ReInitDelay:    ReInitDelaySeconds * time.Second,
		ResendDelay:    ResendDelaySeconds * time.Second,
		RoutingKey:     routingKey,

		done:        make(chan bool),
		notifyReady: make(chan interface{}),
	}
	client.reconnectDelay = client.ReconnectDelay
	client.reInitDelay = client.ReInitDelay
	client.resendDelay = client.ResendDelay

	go client.handleReconnect(urlString)

	return &client, nil
}

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

// Close will cleanly shutdown the channel and connection.
func (client *ClientRabbitMQ) Close() error {
	if !client.isReady {
		return wraperror.Errorf(errForPackage, "already closed: not connected to the server")
	}

	close(client.done)
	close(client.notifyReady)

	// IMPROVE:  connection.Close() closes underlying channels so do we need this?
	err := client.channel.Close()
	if err != nil {
		client.log(4004, err)
	}

	err = client.connection.Close()
	if err != nil {
		client.log(4005, err)
	}

	client.isReady = false

	return nil
}

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
//   - prefetch is the number of deliveries for the rabbit mq client to prefetch
//     from the server.  a good rule of thumb is to prefetch as many as there
//     will be clients processing the deliveries.
func (client *ClientRabbitMQ) Consume(prefetch int) (<-chan amqp.Delivery, error) {
	if !client.isReady {
		// wait for client to be ready
		<-client.notifyReady
	}

	if prefetch < 0 {
		prefetch = 1
	}

	if err := client.channel.Qos(
		prefetch, // prefetch count (should set to the number of load goroutines?)
		0,        // prefetch size
		false,    // global
	); err != nil {
		return nil, wraperror.Errorf(err, "unable to set the RabbitMQ channel quality of service")
	}

	channel, err := client.channel.Consume(
		client.QueueName, // queue
		"",               // consumer uid, blank auto-generates a uid
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local (not supported by rabbitmq)
		false,            // no-wait, false seems safest
		nil,              // args
	)

	return channel, wraperror.Errorf(err, wraperror.NoMessage)
}

// Init initializes a single RabbitMQ client that will automatically
// attempt to connect to the server.
func Init(client *ClientRabbitMQ, urlString string) *ClientRabbitMQ {
	// set up defaults if none provided
	if client.ReconnectDelay <= 0 {
		client.ReconnectDelay = ReconnectDelaySeconds * time.Second
	}

	if client.ReInitDelay <= 0 {
		client.ReInitDelay = ReInitDelaySeconds * time.Second
	}

	if client.ResendDelay <= 0 {
		client.ResendDelay = ResendDelaySeconds * time.Second
	}

	// set up internals
	client.done = make(chan bool)
	client.notifyReady = make(chan interface{})

	client.reconnectDelay = client.ReconnectDelay
	client.reInitDelay = client.ReInitDelay
	client.resendDelay = client.ResendDelay

	go client.handleReconnect(urlString)

	return client
}

// Push will push data onto the queue and wait for a confirm.
// If no confirm is received by the resendTimeout,
// it re-sends messages until a confirm is received.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (client *ClientRabbitMQ) Push(ctx context.Context, record queues.Record) error {
	if !client.isReady {
		// wait for client to be ready
		<-client.notifyReady
	}

	for {
		err := client.UnsafePush(ctx, record)
		if err != nil {
			client.log(3001, client.resendDelay, record.GetMessageID(), err)
			select {
			case <-client.done:
				return wraperror.Errorf(errForPackage, "client is shutting down")
			case <-time.After(client.resendDelay):
				client.resendDelay = client.progressiveDelay(client.resendDelay)
			}

			continue
		}
		select {
		case confirm := <-client.notifyConfirm:
			if confirm.Ack {
				// reset resend delay
				client.resendDelay = client.ResendDelay

				return nil
			}
		case <-time.After(client.resendDelay):
			client.resendDelay = client.progressiveDelay(client.resendDelay)
		}
		client.log(3002, client.resendDelay, record.GetMessageID())
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for if the server will
// receive the message.
func (client *ClientRabbitMQ) UnsafePush(ctx context.Context, record queues.Record) error {
	if !client.isReady {
		// wait for client to be ready
		<-client.notifyReady
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := client.channel.PublishWithContext(
		ctx,                 // context
		client.ExchangeName, // exchange name
		client.RoutingKey,   // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{ //nolint
			Body:         []byte(record.GetMessage()),
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			MessageId:    record.GetMessageID(),
		},
	)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (client *ClientRabbitMQ) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (client *ClientRabbitMQ) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

// Connect will create a new AMQP connection.
func (client *ClientRabbitMQ) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, wraperror.Errorf(err, "unable to dial RabbitMQ address: %v", addr)
	}

	client.changeConnection(conn)
	client.log(2002, addr)

	return conn, nil
}

func (client *ClientRabbitMQ) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = createLogger()
	}

	return client.logger
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (client *ClientRabbitMQ) handleReconnect(addr string) {
	for {
		client.isReady = false
		client.log(2001, addr)

		conn, err := client.connect(addr)
		if err != nil {
			client.log(4001, client.reconnectDelay, err)

			select {
			case <-client.done:
				return
			case <-time.After(client.reconnectDelay):
				client.reconnectDelay = client.progressiveDelay(client.reconnectDelay)
			}

			continue
		}

		// reset client.reconnectDealy
		client.reconnectDelay = client.ReconnectDelay
		if done := client.handleReInit(conn); done {
			break
		}
	}
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels.
func (client *ClientRabbitMQ) handleReInit(conn *amqp.Connection) bool {
	for {
		client.isReady = false

		err := client.init(conn)
		if err != nil {
			client.log(4002, client.reInitDelay, err)

			select {
			case <-client.done:
				return true
			case <-time.After(client.reInitDelay):
				client.reInitDelay = client.progressiveDelay(client.reInitDelay)
			}

			continue
		}

		// reset the init delay
		client.reInitDelay = client.ReInitDelay

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.log(2003)

			return false
		case <-client.notifyChanClose:
			client.log(2004)
		}
	}
}

// init will initialize channel, declare the exchange, and declare the queue.
func (client *ClientRabbitMQ) init(conn *amqp.Connection) error {
	defer func() {
		if r := recover(); r != nil {
			client.log(4003, r)
		}
	}()

	channel, err := conn.Channel()
	if err != nil {
		return wraperror.Errorf(err, "unable to initialize the RabbitMQ channel")
	}

	err = channel.Confirm(false)
	if err != nil {
		return wraperror.Errorf(err, "Confirm")
	}

	err = channel.ExchangeDeclare(
		client.ExchangeName, // name
		"direct",            // type
		false,               // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return wraperror.Errorf(err, "unable to declare the RabbitMQ exchange")
	}

	var queue amqp.Queue

	queue, err = channel.QueueDeclare(
		client.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return wraperror.Errorf(err, "unable to initialize the RabbitMQ queue")
	}

	err = channel.QueueBind(
		queue.Name,          // queue name
		client.RoutingKey,   // routing key
		client.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return wraperror.Errorf(err, "unable to bind the RabbitMQ queue")
	}

	client.changeChannel(channel)
	client.isReady = true
	client.notifyReady <- struct{}{}
	client.log(2005)

	return nil
}

// Log message.
func (client *ClientRabbitMQ) log(messageNumber int, details ...interface{}) {
	client.getLogger().Log(messageNumber, details...)
}

// Progressively increase the retry delay.
func (client *ClientRabbitMQ) progressiveDelay(delay time.Duration) time.Duration {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(delay/time.Second)))
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	return delay + time.Duration(r.Int64())*time.Second
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func createLogger() logging.Logging {
	options := []interface{}{
		&logging.OptionCallerSkip{Value: OptionCallerSkip},
	}

	result, err := logging.NewSenzingLogger(ComponentID, IDMessages, options...)
	if err != nil {
		panic(err)
	}

	return result
}
