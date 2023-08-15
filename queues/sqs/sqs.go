package sqs

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/roncewind/go-util/util"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-queueing/queues"
)

type Client struct {
	DeadLetterQueueURL string
	JSONOutput         bool
	LogLevel           string
	QueueName          string
	QueueURL           *string
	// desired / default delay durations
	MaxDelay       time.Duration
	ReconnectDelay time.Duration
	ResendDelay    time.Duration
	RoutingKey     string

	isReady bool
	logger  logging.LoggingInterface
	// current delay durations
	reconnectDelay time.Duration
	resendDelay    time.Duration

	// region       string //: configure region??
	sqsClient    *sqs.Client
	sqsDLQClient *sqs.Client
	sqsURL       *sqs.GetQueueUrlOutput
}

// ----------------------------------------------------------------------------

// New creates a single SQS client
func NewClient(ctx context.Context, urlString string, logLevel string, jsonOutput bool) (*Client, error) {
	client := Client{
		LogLevel:       logLevel,
		JSONOutput:     jsonOutput,
		MaxDelay:       10 * time.Minute,
		ReconnectDelay: 2 * time.Second,
		ResendDelay:    1 * time.Second,
	}
	client.SetLogLevel(context.Background(), logLevel)
	client.logger = client.getLogger()
	// load the default aws config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return &client, err
	}
	client.sqsClient = sqs.NewFromConfig(cfg)

	err = client.getQueueURL(ctx, urlString)
	if err != nil {
		return &client, err
	}

	client.reconnectDelay = client.ReconnectDelay
	client.resendDelay = client.ResendDelay
	client.getRedrivePolicy(ctx)
	client.isReady = true
	client.log(2005, *client.QueueURL, client.QueueName, client.DeadLetterQueueURL)
	return &client, nil
}

// ----------------------------------------------------------------------------

// Interrogates the given input URL and retrieves the SQS URL or uses it as it.
func (client *Client) getQueueURL(ctx context.Context, urlString string) error {
	u, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("unable to parse SQS URL string %w", err)
	}
	queryMap, _ := url.ParseQuery(u.RawQuery)
	if len(queryMap["queue-name"]) < 1 {
		client.QueueURL = &urlString
	} else {
		client.QueueName = queryMap["queue-name"][0]
		// Get the URL for the queue
		input := &sqs.GetQueueUrlInput{
			QueueName: &client.QueueName,
		}
		sqsURL, err := client.sqsClient.GetQueueUrl(ctx, input)
		if err != nil {
			return fmt.Errorf("unable to retrieve SQS URL from: %s, error: %w", urlString, err)
		}
		client.sqsURL = sqsURL
		client.QueueURL = sqsURL.QueueUrl
	}
	return nil
}

// ----------------------------------------------------------------------------

// Pulls the SQS queues Redrive Policy from AWS and sets the dead letter
// queue URL in the client
func (client *Client) getRedrivePolicy(ctx context.Context) {
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(*client.QueueURL),
		AttributeNames: []types.QueueAttributeName{
			// types.QueueAttributeNameAll,
			types.QueueAttributeNameRedrivePolicy,
		},
	}
	queueAttributes, err := client.sqsClient.GetQueueAttributes(ctx, params)
	if err != nil {
		client.log(4006, err)
		return
	}
	redrive := queueAttributes.Attributes[string(types.QueueAttributeNameRedrivePolicy)]
	var redrivePolicy redrivePolicy
	err = json.Unmarshal([]byte(redrive), &redrivePolicy)
	if err != nil {
		client.log(4007, err)
		return
	}
	fields := strings.Split(redrivePolicy.DeadLetterTargetArn, ":")
	client.DeadLetterQueueURL = fmt.Sprintf("https://queue.amazonaws.com/%s/%s", fields[4], fields[5])
	if len(client.DeadLetterQueueURL) > 0 {
		// load the default aws config
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			client.sqsDLQClient = nil
		} else {
			client.sqsDLQClient = sqs.NewFromConfig(cfg)
		}
	}
}

// ----------------------------------------------------------------------------

// internal struct to parse the AWS Redrive Policy attribute JSON
type redrivePolicy struct {
	DeadLetterTargetArn string `json:"deadLetterTargetArn"`
	MaxReceiveCount     int    `json:"maxReceiveCount"`
}

// ----------------------------------------------------------------------------

// send a message to a queue.
func (client *Client) sendDeadRecord(ctx context.Context, record types.Message) (err error) {

	if client.sqsDLQClient == nil {
		return errors.New("no dead letter queue found")
	}

	// Send a message with attributes to the given queue
	messageInput := &sqs.SendMessageInput{
		DelaySeconds:      0,
		MessageAttributes: record.MessageAttributes,
		MessageBody:       aws.String(*record.Body),
		QueueUrl:          &client.DeadLetterQueueURL,
	}

	resp, err := client.sqsDLQClient.SendMessage(ctx, messageInput)
	if err != nil {
		client.log(4008, err)
		return
	}

	client.log(2006, *resp.MessageId)

	return nil
}

// ----------------------------------------------------------------------------

// send a message to a queue.
func (client *Client) sendRecord(ctx context.Context, record queues.Record) (err error) {

	// Send a message with attributes to the given queue
	messageInput := &sqs.SendMessageInput{
		DelaySeconds: 0,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(record.GetMessageID()),
			},
		},
		MessageBody: aws.String(record.GetMessage()),
		QueueUrl:    client.QueueURL,
	}

	resp, err := client.sqsClient.SendMessage(ctx, messageInput)
	if err != nil {
		client.log(4009, err)
		return
	}

	client.log(2006, *resp.MessageId)

	return nil
}

// ----------------------------------------------------------------------------

// send a message to a queue.
func (client *Client) sendRecordBatch(ctx context.Context, records []queues.Record) (err error) {
	var messages []types.SendMessageBatchRequestEntry = make([]types.SendMessageBatchRequestEntry, len(records))
	// r := rand.New(rand.NewSource(time.Now().Unix()))
	// id := r.Intn(10000)
	r, _ := rand.Int(rand.Reader, big.NewInt(10000))
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	id := r.Int64()
	i := int64(0)
	for _, record := range records {
		if record != nil {
			messages[i] = types.SendMessageBatchRequestEntry{
				DelaySeconds: 0,
				Id:           aws.String(fmt.Sprintf("%d", id+i)),
				MessageAttributes: map[string]types.MessageAttributeValue{
					"MessageID": {
						DataType:    aws.String("String"),
						StringValue: aws.String(record.GetMessageID()),
					},
				},
				MessageBody: aws.String(record.GetMessage()), //?  aws.String(string(utils.Base64Encode([]byte(body)))),
			}
			i++
		}
	}
	// bail of we have no messages to send
	if i <= 0 {
		return
	}
	// Send a message with attributes to the given queue
	messageInput := &sqs.SendMessageBatchInput{
		Entries:  messages[0:i],
		QueueUrl: client.QueueURL,
	}

	resp, err := client.sqsClient.SendMessageBatch(ctx, messageInput)
	if err != nil {
		client.log(4010, err)
	}
	if resp != nil {
		if len(resp.Failed) > 0 {
			for _, fail := range resp.Failed {
				client.log(4011, fail.Id, fail.Message)
			}
		}
		client.log(2007, len(resp.Successful))
	}

	return
}

// ----------------------------------------------------------------------------

// progressively increase the retry delay
func (client *Client) progressiveDelay(delay time.Duration) time.Duration {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(delay/time.Second)))
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	newDelay := delay + time.Duration(r.Int64())*time.Second
	if newDelay > client.MaxDelay {
		return client.MaxDelay
	}
	return newDelay
}

// ----------------------------------------------------------------------------

// PushDeadRecord will push an erroneous record onto the DLQ.
// TODO: work on resend with delay...
func (client *Client) PushDeadRecord(ctx context.Context, record types.Message) error {
	return client.sendDeadRecord(ctx, record)
}

// ----------------------------------------------------------------------------

// Push will push data onto the queue and wait for a response.
// TODO: work on resend with delay...
func (client *Client) Push(ctx context.Context, record queues.Record) error {

	if !client.isReady {
		return fmt.Errorf("the SQS client is not ready")
	}

	for {
		err := client.sendRecord(ctx, record)
		if err != nil {
			client.log(3001, client.resendDelay, record.GetMessageID(), err)
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled: %v", ctx.Err())
			case <-time.After(client.resendDelay):
				//:  resend forever???
				client.resendDelay = client.progressiveDelay(client.resendDelay)
			}
			continue
		} else {
			//reset the resend delay
			client.resendDelay = client.ResendDelay
			return nil
		}
	}
}

// ----------------------------------------------------------------------------

// Push will push data onto the queue and wait for a response.
// TODO: work on resend with delay????
func (client *Client) PushBatch(ctx context.Context, recordchan <-chan queues.Record) error {

	if !client.isReady {
		return fmt.Errorf("the SQS client is not ready")
	}
	i := 0
	records := make([]queues.Record, 10)
	for record := range util.OrDone(ctx, recordchan) {
		records[i] = record
		fmt.Println("batch push record:!", record.GetMessageID())
		i++
		if i >= 10 {
			err := client.sendRecordBatch(ctx, records)
			if err != nil {
				client.log(4010, err)
			}
			i = 0
			records = make([]queues.Record, 10)
		}
	}
	//handle a last partial batch
	if i > 0 {
		err := client.sendRecordBatch(ctx, records)
		if err != nil {
			client.log(4012, err)
		}
	}
	return nil
}

// ----------------------------------------------------------------------------

// receive a message from a queue.
func (client *Client) receiveMessage(ctx context.Context, visibilitySeconds int32) (*sqs.ReceiveMessageOutput, error) {

	// Receive a message with attributes to the given queue
	receiveInput := &sqs.ReceiveMessageInput{
		QueueUrl:              client.QueueURL,
		MessageAttributeNames: []string{"All"},
		MaxNumberOfMessages:   10,
		VisibilityTimeout:     visibilitySeconds,
	}

	msg, err := client.sqsClient.ReceiveMessage(ctx, receiveInput)
	if err != nil {
		client.log(4013, err)
		return nil, fmt.Errorf("error receiving records %w", err)
	}
	if msg.Messages == nil || len(msg.Messages) <= 0 {
		client.log(4014)
		return nil, fmt.Errorf("no messages")
	}

	return msg, nil
}

// ----------------------------------------------------------------------------

// Consume will continuously put queue messages on the channel.
func (client *Client) Consume(ctx context.Context, visibilitySeconds int32) (<-chan types.Message, error) {
	if !client.isReady {
		return nil, fmt.Errorf("the SQS client is not ready")
	}
	outChan := make(chan types.Message, 10)
	go func() {
		messageCount := 0
		for {
			output, err := client.receiveMessage(ctx, visibilitySeconds)

			if err != nil {
				time.Sleep(client.reconnectDelay)
				client.reconnectDelay = client.progressiveDelay(client.reconnectDelay)
			} else {
				select {
				case <-ctx.Done():
					fmt.Println("DEBUG: final messageCount:", messageCount)
					return
				default:
					for _, m := range output.Messages {
						outChan <- m
						messageCount++
						fmt.Println("DEBUG: read message from SQS. messageCount:", messageCount)
					}
					// reset the reconnectDelay
					client.reconnectDelay = client.ReconnectDelay
				}
			}
		}
	}()
	return outChan, nil
}

// ----------------------------------------------------------------------------

// Remove a message from the SQS queue
func (client *Client) RemoveMessage(ctx context.Context, msg types.Message) error {
	deleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      client.QueueURL,
		ReceiptHandle: msg.ReceiptHandle,
	}

	client.log(2008, msg.MessageId)

	_, err := client.sqsClient.DeleteMessage(ctx, deleteMessageInput)
	if err != nil {
		client.log(4015, err)
		return fmt.Errorf("error deleting records, %w", err)
	}
	return nil
}

// ----------------------------------------------------------------------------

// Remove a message from the SQS queue
func (client *Client) SetMessageVisibility(ctx context.Context, msg types.Message, seconds int32) error {

	if seconds < 0 {
		seconds = 0
	}

	if seconds > 12*60*60 { //12 hour max
		seconds = 12 * 60 * 60
	}

	setVisibilityInput := &sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     msg.ReceiptHandle,
		QueueUrl:          client.QueueURL,
		VisibilityTimeout: seconds,
	}

	client.log(2009, msg.MessageId)
	_, err := client.sqsClient.ChangeMessageVisibility(ctx, setVisibilityInput)
	if err != nil {
		client.log(4016, msg.MessageId, err)
		return err
	}
	return nil
}

// ----------------------------------------------------------------------------

// Close will cleanly shutdown the channel and connection.
func (client *Client) Close() error {
	return nil
}

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (c *Client) getLogger() logging.LoggingInterface {
	var err error = nil
	if c.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		c.logger, err = logging.NewSenzingToolsLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return c.logger
}

// Log message.
func (c *Client) log(messageNumber int, details ...interface{}) {
	if c.JSONOutput {
		c.getLogger().Log(messageNumber, details...)
	} else {
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...))
	}
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (c *Client) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = c.getLogger().SetLogLevel(logLevelName)
	return err
}
