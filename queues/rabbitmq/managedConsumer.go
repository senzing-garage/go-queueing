package rabbitmq

import (
	"context"
	"fmt"
	"runtime"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/sourcegraph/conc/pool"
)

// define a structure that will implement the Job interface
type RabbitConsumerJob struct {
	delivery  amqp.Delivery
	engine    *senzing.SzEngine
	id        int
	usedCount int
	withInfo  bool
}

var jobPool chan *RabbitConsumerJob

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

// Job interface implementation:
// Execute() is run once for each Job
func (job *RabbitConsumerJob) Execute(ctx context.Context) error {
	// increment the number of times this job struct was used and return to the pool
	defer func() {
		job.usedCount++
		jobPool <- job
	}()
	// fmt.Printf("Received a message- msgId: %s, msgCnt: %d, ConsumerTag: %s\n", id, j.delivery.MessageCount, j.delivery.ConsumerTag)
	record, newRecordErr := record.NewRecord(string(job.delivery.Body))
	if newRecordErr == nil {
		flags := senzing.SzWithoutInfo
		if job.withInfo {
			flags = senzing.SzWithInfo
		}
		result, err := (*job.engine).AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		if err != nil {
			return fmt.Errorf(
				"add record error, record id: %s, message id: %s, result: %s, %w",
				job.delivery.MessageId,
				record.ID,
				result,
				err,
			)
		}

		// when we successfully process a delivery, acknowledge it.
		return job.delivery.Ack(false)
	}
	// when we get an invalid delivery, negatively acknowledge and send to the dead letter queue
	err := job.delivery.Nack(false, false)
	return fmt.Errorf("invalid deliver from RabbitMQ, message id: %s, %w", job.delivery.MessageId, err)
}

// Whenever Execute() returns an error or panics, this is called
func (job *RabbitConsumerJob) OnError(err error) {
	_ = err
	if job.delivery.Redelivered {
		// swallow any error, it'll timeout and be redelivered
		_ = job.delivery.Nack(false, false)
	} else {
		// swallow any error, it'll timeout and be redelivered
		_ = job.delivery.Nack(false, true)
	}
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
	szEngine *senzing.SzEngine,
	withInfo bool,
	logLevel string,
	jsonOutput bool,
) error {
	_ = jsonOutput

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

	logger.Log(2012, numberOfWorkers)

	// setup jobs that will be used to process RabbitMQ deliveries
	jobPool = make(chan *RabbitConsumerJob, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		jobPool <- &RabbitConsumerJob{
			engine:    szEngine,
			id:        i,
			usedCount: 0,
			withInfo:  withInfo,
		}
	}

	client, err := NewClient(urlString)
	if err != nil {
		return fmt.Errorf("unable to get a new RabbitMQ client %w", err)
	}
	defer client.Close()

	deliveries, err := client.Consume(numberOfWorkers)
	if err != nil {
		return fmt.Errorf("unable to get a new RabbitMQ delivery channel %w", err)
	}

	p := pool.New().WithMaxGoroutines(numberOfWorkers)
	jobCount := 0
	for delivery := range deliveries {
		job := <-jobPool
		job.delivery = delivery
		p.Go(func() {
			err := job.Execute(ctx)
			if err != nil {
				job.OnError(err)
			}
		})

		jobCount++
		if jobCount%10000 == 0 {
			logger.Log(2010, jobCount)
		}
	}

	// Wait for all the records in the record channel to be processed
	p.Wait()

	// clean up after ourselves
	close(jobPool)
	// drain the job pool
	var job *RabbitConsumerJob
	ok := true
	for ok {
		job, ok = <-jobPool
		logger.Log(2011, job.id, job.usedCount)
	}

	return nil
}
