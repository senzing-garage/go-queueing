package sqs

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/senzing-garage/go-queueing/queues"
	"github.com/sourcegraph/conc/pool"
)

var clientPool chan *Client

// ----------------------------------------------------------------------------
// DELETE ME??  This is the old way of doing things, one record at a time.

// read a record from the record channel and push it to the RabbitMQ queue
// func processRecord(ctx context.Context, record queues.Record, newClientFn func() (*Client, error)) (err error) {
// 	client := <-clientPool
// 	err = client.Push(ctx, record)
// 	if err != nil {
// 		// on error, create a new SQS client
// 		err = fmt.Errorf("error pushing record batch %w", err)
// 		//put a new client in the pool, dropping the current one
// 		newClient, newClientErr := newClientFn()
// 		if newClientErr != nil {
// 			err = fmt.Errorf("error creating a new client %w", newClientErr)
// 		} else {
// 			clientPool <- newClient
// 		}
// 		// make sure to close the old client
// 		return fmt.Errorf("error creating a new client %w %w", client.Close(), err)
// 	}
// 	// return the client to the pool when done
// 	clientPool <- client
// 	return
// }

// ----------------------------------------------------------------------------

// read a record from the record channel and push it to the SQS queue
func processRecordBatch(ctx context.Context, recordchan <-chan queues.Record, newClientFn func() (*Client, error)) (err error) {
	client := <-clientPool
	err = client.PushBatch(ctx, recordchan)
	if err != nil {
		// on error, create a new SQS client
		err = fmt.Errorf("error pushing record batch %w", err)
		// put a new client in the pool, dropping the current one
		newClient, newClientErr := newClientFn()
		if newClientErr != nil {
			err = fmt.Errorf("error creating a new client %w", newClientErr)
		} else {
			clientPool <- newClient
		}
		// make sure to close the old client
		return fmt.Errorf("error creating a new client %w %w", client.Close(), err)
	}
	// return the client to the pool when done
	clientPool <- client
	return
}

// ----------------------------------------------------------------------------

// Starts a number of workers that push Records in the record channel to
// the given queue.
// - Workers restart when they are killed or die.
// - respond to standard system signals.
func StartManagedProducer(ctx context.Context, urlString string, numberOfWorkers int, recordchan <-chan queues.Record, logLevel string, jsonOutput bool) {

	// default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}
	fmt.Println(time.Now(), "Number of producer workers:", numberOfWorkers)

	ctx, cancel := context.WithCancel(ctx)
	if err := SetLogLevel(ctx, logLevel); err != nil {
		log(3003, logLevel, err)
	}

	clientPool = make(chan *Client, numberOfWorkers)
	newClientFn := func() (*Client, error) { return NewClient(ctx, urlString, logLevel, jsonOutput) }

	// populate an initial client pool
	go func() { _ = createClients(ctx, numberOfWorkers, newClientFn) }()

	p := pool.New().WithMaxGoroutines(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		p.Go(func() {
			i := i
			err := processRecordBatch(ctx, recordchan, newClientFn)
			if err != nil {
				fmt.Println("Worker[", i, "] error:", err)
			}
		})
	}

	// Wait for all the records in the record channel to be processed
	p.Wait()

	// clean up after ourselves
	cancel()
	fmt.Println(time.Now(), "Clean up job queue and client pool.")
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

// create a number of clients and put them into the client queue
func createClients(ctx context.Context, numOfClients int, newClientFn func() (*Client, error)) error {
	_ = ctx
	countOfClientsCreated := 0
	var errorStack error
	for i := 0; i < numOfClients; i++ {
		client, err := newClientFn()
		if err != nil {
			errorStack = fmt.Errorf("error creating a new client %w", err)
		} else {
			countOfClientsCreated++
			clientPool <- client
		}
	}
	fmt.Println(time.Now(), countOfClientsCreated, "SQS clients created,", numOfClients, "requested")
	return errorStack
}
