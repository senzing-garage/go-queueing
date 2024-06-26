package rabbitmq

import (
	"context"
	"fmt"
	"runtime"

	"github.com/senzing-garage/go-queueing/queues"
	"github.com/sourcegraph/conc/pool"
)

var clientPool chan *Client

// ----------------------------------------------------------------------------
// Job implementation
// ----------------------------------------------------------------------------

// define a structure that will implement the Job interface
// type RabbitProducerJob struct {
// 	id          int
// 	newClientFn func() (*Client, error)
// 	record      queues.Record
// }

// ----------------------------------------------------------------------------

// read a record from the record channel and push it to the RabbitMQ queue
func processRecord(ctx context.Context, record queues.Record, newClientFn func() (*Client, error)) (err error) {
	client := <-clientPool
	err = client.Push(ctx, record)
	if err != nil {
		// on error, create a new RabbitMQ client
		err = fmt.Errorf("error pushing record, creating new client %w", err)
		// put a new client in the pool, dropping the current one
		newClient, newClientErr := newClientFn()
		if newClientErr != nil {
			err = fmt.Errorf("error creating new client %w %w", newClientErr, err)
		} else {
			clientPool <- newClient
		}
		// make sure to close the old client
		closeErr := client.Close()
		if closeErr != nil {
			err = fmt.Errorf("error closing client %w %w", closeErr, err)
		}
		return
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
func StartManagedProducer(ctx context.Context, urlString string, numberOfWorkers int, recordchan <-chan queues.Record, logLevel string) {

	// default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)
	if err := SetLogLevel(ctx, logLevel); err != nil {
		log(3003, logLevel, err)
	}
	logger = getLogger()

	log(2013, numberOfWorkers)
	clientPool = make(chan *Client, numberOfWorkers)
	newClientFn := func() (*Client, error) { return NewClient(urlString) }

	// populate an initial client pool
	go func() { _ = createClients(ctx, numberOfWorkers, newClientFn) }()

	p := pool.New().WithMaxGoroutines(numberOfWorkers)
	for record := range recordchan {
		record := record
		p.Go(func() {
			err := processRecord(ctx, record, newClientFn)
			if err != nil {
				log(4006, record.GetMessageID(), err)
			}
		})
	}

	// Wait for all the records in the record channel to be processed
	p.Wait()

	// clean up after ourselves
	cancel()
	log(2014)
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

// create a number of clients and put them into the client queue
func createClients(ctx context.Context, numOfClients int, newClientFn func() (*Client, error)) error {
	_ = ctx
	countOfClientsCreated := 0
	var errorStack error
	for i := 0; i < numOfClients; i++ {
		client, err := newClientFn()
		if err != nil {
			errorStack = fmt.Errorf("error creating new client %w", err)
		} else {
			countOfClientsCreated++
			clientPool <- client
		}
	}

	return errorStack
}
