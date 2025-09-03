package sqs

import (
	"context"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/roncewind/go-util/util"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/sourcegraph/conc/pool"
)

// Define a structure that will implement the ConsumerJobSqs interface.
type ConsumerJobSqs struct {
	client    *ClientSqs
	engine    senzing.SzEngine
	id        int
	message   types.Message
	startTime time.Time
	usedCount int
	withInfo  bool
}

var jobPool chan ConsumerJobSqs

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

// Job interface implementation:
// Execute() is run once for each Job.
func (job *ConsumerJobSqs) Execute(ctx context.Context, visibilitySeconds int32) error {
	var err error
	// Increment the number of times this job struct was used and return to the pool.
	defer func() {
		job.usedCount++
		jobPool <- *job
	}()

	record, newRecordErr := record.NewRecord(*job.message.Body)
	if newRecordErr != nil {
		// fmt.Println(time.Now(), "ERROR: Invalid delivery from SQS. msg id:", *j.message.MessageId)
		// When we get an invalid delivery, send to the dead letter queue.
		err := job.client.PushDeadRecord(ctx, job.message)
		if err != nil {
			return wraperror.Errorf(
				err,
				"unable to push message to the dead letter queue, record id: %s, message id: %s",
				*job.message.MessageId,
				record.ID,
			)
		}
	}

	err = job.bob(ctx, record, visibilitySeconds)

	return err
}

// Whenever Execute() returns an error or panics, this is called.
func (job *ConsumerJobSqs) OnError(ctx context.Context, err error) {
	_ = err
	// fmt.Println("ERROR: Worker", err)
	// fmt.Println("ERROR: Failed to add record. msg id:", *j.message.MessageId)
	_ = job.client.PushDeadRecord(ctx, job.message)
	// if err != nil {
	// 	fmt.Println("ERROR: Pushing message to the dead letter queue. msg id:", *j.message.MessageId, "error:", err)
	// }
	_ = err
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Starts a number of workers that read Records from the given queue and add
// them to Senzing.
// - Workers restart when they are killed or die.
// - respond to standard system signals.
func StartManagedConsumer(
	ctx context.Context,
	urlString string,
	numberOfWorkers int,
	szEngine senzing.SzEngine,
	withInfo bool,
	visibilitySeconds int32,
	logLevel string,
	jsonOutput bool,
) error {
	if szEngine == nil {
		return wraperror.Errorf(errPackage, "the Sz Engine is not set, unable to start the managed consumer")
	}

	// default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := createLogger()

	err := logger.SetLogLevel(logLevel)
	if err != nil {
		panic(err)
	}

	client, err := NewClient(ctx, urlString, logLevel, jsonOutput)
	if err != nil {
		return wraperror.Errorf(err, "unable to get a new SQS client")
	}

	defer client.Close()

	logger.Log(2012, numberOfWorkers)

	// setup jobs that will be used to process SQS deliveries
	jobPool = make(chan ConsumerJobSqs, numberOfWorkers)
	for i := range numberOfWorkers {
		jobPool <- ConsumerJobSqs{
			client:    client,
			engine:    szEngine,
			id:        i,
			usedCount: 0,
			withInfo:  withInfo,
		}
	}

	messages, err := client.Consume(ctx, visibilitySeconds)
	if err != nil {
		logger.Log(4019, err)

		return wraperror.Errorf(err, "unable to get a new SQS message channel")
	}

	workerPool := pool.New().WithMaxGoroutines(numberOfWorkers)
	jobCount := 0
	// for message := range messages {
	for message := range util.OrDone(ctx, messages) {
		job := <-jobPool
		job.message = message
		job.startTime = time.Now()

		workerPool.Go(func() {
			err := job.Execute(ctx, visibilitySeconds)
			if err != nil {
				job.OnError(ctx, err)
			}
		})

		jobCount++
		if jobCount%10000 == 0 {
			logger.Log(2010, jobCount)
		}
	}

	// Wait for all the records in the record channel to be processed
	workerPool.Wait()

	// clean up after ourselves
	close(jobPool)
	// drain the job pool
	var job ConsumerJobSqs

	ok := true
	for ok {
		job, ok = <-jobPool
		logger.Log(2011, job.id, job.usedCount)
	}

	return nil
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (job *ConsumerJobSqs) bob(ctx context.Context, record *record.Record, visibilitySeconds int32) error {
	var err error

	visibilityContext, visibilityCancel := context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(time.Duration(visibilitySeconds-1) * time.Second)

		for {
			select {
			case <-ctx.Done():
				// swallow the error, we're done
				_ = job.client.SetMessageVisibility(visibilityContext, job.message, 0)

				return
			case <-visibilityContext.Done():
				return
			case <-ticker.C:
				setVisibilityError := job.client.SetMessageVisibility(visibilityContext, job.message, visibilitySeconds)
				if setVisibilityError != nil {
					// when there's an error setting visibility, let the message requeue
					return
				}
			}
		}
	}()

	flags := senzing.SzWithoutInfo
	if job.withInfo {
		flags = senzing.SzWithInfo
	}

	result, err := job.engine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)

	visibilityCancel()

	if err != nil {
		return wraperror.Errorf(
			err,
			"add record error, record id: %s, message id: %s, result: %s",
			*job.message.MessageId,
			record.ID,
			result,
		)
	}
	// IMPROVE:  what do we do with the "withInfo" data here?

	// when we successfully process a message, delete it.
	// as long as there was no error delete the message from the queue
	err = job.client.RemoveMessage(ctx, job.message)
	if err != nil {
		return wraperror.Errorf(
			err,
			"record not removed from queue, record id: %s, message id: %s",
			*job.message.MessageId,
			record.ID,
		)
	}

	return err
}
