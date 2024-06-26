package sqs

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/roncewind/go-util/util"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/sourcegraph/conc/pool"
)

var jobPool chan Job

// ----------------------------------------------------------------------------
// Job implementation
// ----------------------------------------------------------------------------

// define a structure that will implement the Job interface
type Job struct {
	client    *Client
	engine    sz.SzEngine
	id        int
	message   types.Message
	startTime time.Time
	usedCount int
	withInfo  bool
}

// ----------------------------------------------------------------------------

// Job interface implementation:
// Execute() is run once for each Job
func (j *Job) Execute(ctx context.Context, visibilitySeconds int32) error {
	// increment the number of times this job struct was used and return to the pool
	defer func() {
		j.usedCount++
		jobPool <- *j
	}()
	record, newRecordErr := record.NewRecord(string(*j.message.Body))
	if newRecordErr == nil {
		visibilityContext, visibilityCancel := context.WithCancel(ctx)
		go func() {
			ticker := time.NewTicker(time.Duration(visibilitySeconds-1) * time.Second)
			for {
				select {
				case <-ctx.Done():
					// swallow the error, we're done
					_ = j.client.SetMessageVisibility(visibilityContext, j.message, 0)
					return
				case <-visibilityContext.Done():
					return
				case <-ticker.C:
					setVisibilityError := j.client.SetMessageVisibility(visibilityContext, j.message, visibilitySeconds)
					if setVisibilityError != nil {
						// when there's an error setting visibility, let the message requeue
						return
					}
				}
			}
		}()
		flags := sz.SZ_WITHOUT_INFO
		if j.withInfo {
			flags = sz.SZ_WITH_INFO
		}
		result, err := j.engine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		visibilityCancel()
		if err != nil {
			return fmt.Errorf("add record error, record id: %s, message id: %s, result: %s, %w", *j.message.MessageId, record.ID, result, err)
		}
		// TODO:  what do we do with the "withInfo" data here?

		// when we successfully process a message, delete it.
		// as long as there was no error delete the message from the queue
		err = j.client.RemoveMessage(ctx, j.message)
		if err != nil {
			return fmt.Errorf("record not removed from queue, record id: %s, message id: %s, %w", *j.message.MessageId, record.ID, err)
		}
	} else {
		// fmt.Println(time.Now(), "ERROR: Invalid delivery from SQS. msg id:", *j.message.MessageId)
		// when we get an invalid delivery, send to the dead letter queue
		err := j.client.PushDeadRecord(ctx, j.message)
		if err != nil {
			return fmt.Errorf("unable to push message to the dead letter queue, record id: %s, message id: %s, %w", *j.message.MessageId, record.ID, err)
		}
	}
	return nil
}

// ----------------------------------------------------------------------------

// Whenever Execute() returns an error or panics, this is called
func (j *Job) OnError(ctx context.Context, err error) {
	_ = err
	// fmt.Println("ERROR: Worker error:", err)
	// fmt.Println("ERROR: Failed to add record. msg id:", *j.message.MessageId)
	_ = j.client.PushDeadRecord(ctx, j.message)
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
func StartManagedConsumer(ctx context.Context, urlString string, numberOfWorkers int, szEngine sz.SzEngine, withInfo bool, visibilitySeconds int32, logLevel string) error {

	if szEngine == nil {
		return errors.New("the G2 Engine is not set, unable to start the managed consumer")
	}

	// default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := SetLogLevel(ctx, logLevel); err != nil {
		log(3003, logLevel, err)
	}

	client, err := NewClient(ctx, urlString)
	if err != nil {
		return fmt.Errorf("unable to get a new SQS client, %w", err)
	}
	defer client.Close()
	log(2012, numberOfWorkers)

	// setup jobs that will be used to process SQS deliveries
	jobPool = make(chan Job, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		jobPool <- Job{
			client:    client,
			engine:    szEngine,
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
				job.OnError(ctx, err)
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
	var job Job
	ok := true
	for ok {
		job, ok = <-jobPool
		log(2011, job.id, job.usedCount)
	}

	return nil
}
