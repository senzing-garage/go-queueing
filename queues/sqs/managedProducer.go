package sqs

import (
	"context"
	"fmt"
	"runtime"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-queueing/queues"
	"github.com/sourcegraph/conc/pool"
)

var clientPool chan *ClientSqs

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

	clientPool = make(chan *ClientSqs, numberOfWorkers)
	newClientFn := func() (*ClientSqs, error) { return NewClient(ctx, urlString, logLevel, jsonOutput) }

	// populate an initial client pool
	go func() { _ = createClients(ctx, numberOfWorkers, newClientFn) }()

	workerPool := pool.New().WithMaxGoroutines(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		workerPool.Go(func() {
			i := i

			err := processRecordBatch(ctx, recordchan, newClientFn)
			if err != nil {
				fmt.Println("Worker[", i, "] error:", err) //nolint
			}
		})
	}

	// Wait for all the records in the record channel to be processed
	workerPool.Wait()

	// clean up after ourselves
	cancel()
	close(clientPool)
	// drain the client pool, closing sqs connections
	for len(clientPool) > 0 {
		client, ok := <-clientPool
		if ok && client != nil {
			// swallow the error, we're shutting down anyway
			_ = client.Close()
		}
	}
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Create a number of clients and put them into the client queue.
func createClients(ctx context.Context, numOfClients int, newClientFn func() (*ClientSqs, error)) error {
	_ = ctx
	countOfClientsCreated := 0

	var errorStack error

	for i := range numOfClients {
		_ = i

		client, err := newClientFn()
		if err != nil {
			errorStack = wraperror.Errorf(err, "error creating a new client")
		} else {
			countOfClientsCreated++
			clientPool <- client
		}
	}

	return errorStack
}

// ----------------------------------------------------------------------------
// DELETE ME??  This is the old way of doing things, one record at a time.

// read a record from the record channel and push it to the RabbitMQ queue
// func processRecord(ctx context.Context, record queues.Record, newClientFn func() (*Client, error)) (err error) {
// 	client := <-clientPool
// 	err = client.Push(ctx, record)
// 	if err != nil {
// 		// on error, create a new SQS client
// 		err = wraperror.Errorf(err, "error pushing record batch")
// 		//put a new client in the pool, dropping the current one
// 		newClient, newClientErr := newClientFn()
// 		if newClientErr != nil {
// 			err = wraperror.Errorf(err, "error creating a new client %w", newClientErr)
// 		} else {
// 			clientPool <- newClient
// 		}
// 		// make sure to close the old client
// 		return wraperror.Errorf(err, "error creating a new client %w", client.Close())
// 	}
// 	// return the client to the pool when done
// 	clientPool <- client
// 	return
// }

// ----------------------------------------------------------------------------

// Read a record from the record channel and push it to the SQS queue.
func processRecordBatch(
	ctx context.Context,
	recordchan <-chan queues.Record,
	newClientFn func() (*ClientSqs, error),
) error {
	var err error

	client := <-clientPool

	err = client.PushBatch(ctx, recordchan)
	if err != nil {
		// on error, create a new SQS client
		err = wraperror.Errorf(err, "error pushing record batch")
		// put a new client in the pool, dropping the current one
		newClient, newClientErr := newClientFn()
		if newClientErr != nil {
			err = wraperror.Errorf(err, "error creating a new client %s", newClientErr.Error())
		} else {
			clientPool <- newClient
		}
		// make sure to close the old client
		return wraperror.Errorf(err, "error creating a new client %s", client.Close().Error())
	}
	// return the client to the pool when done
	clientPool <- client

	return err
}
