package rabbitmq

import (
	"context"
	"runtime"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-queueing/queues"
	"github.com/sourcegraph/conc/pool"
)

var clientPool chan *ClientRabbitMQ

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Starts a number of workers that push Records in the record channel to
// the given queue.
// - Workers restart when they are killed or die.
// - respond to standard system signals.
func StartManagedProducer(
	ctx context.Context,
	urlString string,
	numberOfWorkers int,
	recordchan <-chan queues.Record,
	logLevel string,
	jsonOutput bool,
) {
	_ = jsonOutput

	// default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)

	logger := createLogger()

	err := logger.SetLogLevel(logLevel)
	if err != nil {
		panic(err)
	}

	logger.Log(2013, numberOfWorkers)
	clientPool = make(chan *ClientRabbitMQ, numberOfWorkers)
	newClientFn := func() (*ClientRabbitMQ, error) { return NewClient(urlString) }

	// populate an initial client pool
	go func() { _ = createClients(ctx, numberOfWorkers, newClientFn) }()

	workerPool := pool.New().WithMaxGoroutines(numberOfWorkers)
	for record := range recordchan {
		workerPool.Go(func() {
			err := processRecord(ctx, record, newClientFn)
			if err != nil {
				logger.Log(4006, record.GetMessageID(), err)
			}
		})
	}

	// Wait for all the records in the record channel to be processed
	workerPool.Wait()

	// clean up after ourselves
	cancel()
	logger.Log(2014)
	close(clientPool)
	// drain the client pool, closing rabbit mq connections
	for len(clientPool) > 0 {
		client, ok := <-clientPool
		if ok && client != nil {
			// swallow any errors on cleanup
			_ = client.Close()
		}
	}
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Read a record from the record channel and push it to the RabbitMQ queue.
func processRecord(ctx context.Context, record queues.Record, newClientFn func() (*ClientRabbitMQ, error)) error {
	var err error

	client := <-clientPool

	err = client.Push(ctx, record)
	if err != nil {
		// on error, create a new RabbitMQ client
		err = wraperror.Errorf(err, "error pushing record, creating new client")
		// put a new client in the pool, dropping the current one
		newClient, errNewClientFn := newClientFn()
		if errNewClientFn != nil {
			err = wraperror.Errorf(err, "error creating new client: %s", errNewClientFn.Error())
		} else {
			clientPool <- newClient
		}
		// make sure to close the old client
		errClose := client.Close()
		if errClose != nil {
			err = wraperror.Errorf(err, "error closing client: %s", errClose.Error())
		}

		return err
	}
	// return the client to the pool when done
	clientPool <- client

	return err
}

// Create a number of clients and put them into the client queue.
func createClients(ctx context.Context, numOfClients int, newClientFn func() (*ClientRabbitMQ, error)) error {
	_ = ctx
	countOfClientsCreated := 0

	var errorStack error

	for i := range numOfClients {
		_ = i

		client, err := newClientFn()
		if err != nil {
			errorStack = wraperror.Errorf(err, "error creating new client")
		} else {
			countOfClientsCreated++
			clientPool <- client
		}
	}

	return errorStack
}
