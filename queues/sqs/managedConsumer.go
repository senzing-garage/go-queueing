package sqs

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/roncewind/go-util/util"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	"github.com/senzing-garage/go-common/record"
	"github.com/sourcegraph/conc/pool"
)

var jobPool chan SQSJob

// ----------------------------------------------------------------------------
// Job implementation
// ----------------------------------------------------------------------------

// define a structure that will implement the Job interface
type SQSJob struct {
	client    *Client
	engine    g2api.G2engine
	id        int
	message   types.Message
	startTime time.Time
	usedCount int
	withInfo  bool
}

// ----------------------------------------------------------------------------

// Job interface implementation:
// Execute() is run once for each Job
func (j *SQSJob) Execute(ctx context.Context, visibilitySeconds int32) error {
	// increment the number of times this job struct was used and return to the pool
	defer func() {
		j.usedCount++
		jobPool <- *j
	}()
	record, newRecordErr := record.NewRecord(string(*j.message.Body))
	if newRecordErr == nil {
		loadID := "Load"
		visibilityContext, visibilityCancel := context.WithCancel(ctx)
		go func() {
			ticker := time.NewTicker(time.Duration(visibilitySeconds-1) * time.Second)
			for {
				select {
				case <-ctx.Done():
					//swallow the error, we're done
					_ = j.client.SetMessageVisibility(visibilityContext, j.message, 0)
					return
				case <-visibilityContext.Done():
					return
				case <-ticker.C:
					setVisibilityError := j.client.SetMessageVisibility(visibilityContext, j.message, visibilitySeconds)
					if setVisibilityError != nil {
						//when there's an error setting visibility, let the message requeue
						return
					}
				}
			}
		}()
		if j.withInfo {
			var flags int64 = 0
			_, withInfoErr := j.engine.AddRecordWithInfo(ctx, record.DataSource, record.Id, record.Json, loadID, flags)
			visibilityCancel()
			if withInfoErr != nil {
				return fmt.Errorf("add record error, record id: %s, message id: %s, %w", *j.message.MessageId, record.Id, withInfoErr)
			}
			//TODO:  what do we do with the "withInfo" data here?
		} else {
			addRecordErr := j.engine.AddRecord(ctx, record.DataSource, record.Id, record.Json, loadID)
			visibilityCancel()
			if addRecordErr != nil {
				return fmt.Errorf("add record error, record id: %s, message id: %s, %w", *j.message.MessageId, record.Id, addRecordErr)
			}
		}

		// when we successfully process a message, delete it.
		//as long as there was no error delete the message from the queue
		err := j.client.RemoveMessage(ctx, j.message)
		if err != nil {
			return fmt.Errorf("record not removed from queue, record id: %s, message id: %s, %w", *j.message.MessageId, record.Id, err)
		}
	} else {
		// fmt.Println(time.Now(), "ERROR: Invalid delivery from SQS. msg id:", *j.message.MessageId)
		// when we get an invalid delivery, send to the dead letter queue
		err := j.client.PushDeadRecord(ctx, j.message)
		if err != nil {
			return fmt.Errorf("unable to push message to the dead letter queue, record id: %s, message id: %s, %w", *j.message.MessageId, record.Id, err)
		}
	}
	return nil
}

// ----------------------------------------------------------------------------

// Whenever Execute() returns an error or panics, this is called
func (j *SQSJob) OnError(err error) {
	// fmt.Println("ERROR: Worker error:", err)
	// fmt.Println("ERROR: Failed to add record. msg id:", *j.message.MessageId)
	_ = j.client.PushDeadRecord(context.Background(), j.message)
	// if err != nil {
	// 	fmt.Println("ERROR: Pushing message to the dead letter queue. msg id:", *j.message.MessageId, "error:", err)
	// }
}

// ----------------------------------------------------------------------------
// -- add records in SQS to Senzing
// ----------------------------------------------------------------------------

// Starts a number of workers that read Records from the given queue and add
// them to Senzing.
// - Workers restart when they are killed or die.
// - respond to standard system signals.
func StartManagedConsumer(ctx context.Context, urlString string, numberOfWorkers int, g2engine g2api.G2engine, withInfo bool, visibilitySeconds int32, logLevel string, jsonOutput bool) error {

	if g2engine == nil {
		return errors.New("the G2 Engine is not set, unable to start the managed consumer")
	}

	//default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := SetLogLevel(ctx, logLevel); err != nil {
		log(3003, logLevel, err)
	}

	client, err := NewClient(ctx, urlString, logLevel, jsonOutput)
	if err != nil {
		return fmt.Errorf("unable to get a new SQS client, %w", err)
	}
	defer client.Close()
	log(2012, numberOfWorkers)

	// setup jobs that will be used to process SQS deliveries
	jobPool = make(chan SQSJob, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		jobPool <- SQSJob{
			client:    client,
			engine:    g2engine,
			id:        i,
			usedCount: 0,
			withInfo:  withInfo,
		}
	}

	messages, err := client.Consume(ctx, visibilitySeconds)
	if err != nil {
		log(4019, err)
		return fmt.Errorf("unable to get a new SQS message channel %w", err)
	}

	p := pool.New().WithMaxGoroutines(numberOfWorkers)
	jobCount := 0
	// for message := range messages {
	for message := range util.OrDone(ctx, messages) {
		job := <-jobPool
		job.message = message
		job.startTime = time.Now()
		p.Go(func() {
			err := job.Execute(ctx, visibilitySeconds)
			if err != nil {
				job.OnError(err)
			}
		})

		jobCount++
		if jobCount%10000 == 0 {
			log(2010, jobCount)
		}
	}

	// Wait for all the records in the record channel to be processed
	p.Wait()

	// clean up after ourselves
	close(jobPool)
	// drain the job pool
	var job SQSJob
	ok := true
	for ok {
		job, ok = <-jobPool
		log(2011, job.id, job.usedCount)
	}

	return nil
}
