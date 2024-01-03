package rabbitmq

import (
	"context"
	"fmt"
	"runtime"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-common/record"
	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/sourcegraph/conc/pool"
)

var jobPool chan *RabbitConsumerJob

// ----------------------------------------------------------------------------
// Job implementation
// ----------------------------------------------------------------------------

// define a structure that will implement the Job interface
type RabbitConsumerJob struct {
	delivery  amqp.Delivery
	engine    *g2api.G2engine
	id        int
	usedCount int
	withInfo  bool
}

// ----------------------------------------------------------------------------

// Job interface implementation:
// Execute() is run once for each Job
func (j *RabbitConsumerJob) Execute(ctx context.Context) error {
	// increment the number of times this job struct was used and return to the pool
	defer func() {
		j.usedCount++
		jobPool <- j
	}()
	// fmt.Printf("Received a message- msgId: %s, msgCnt: %d, ConsumerTag: %s\n", id, j.delivery.MessageCount, j.delivery.ConsumerTag)
	record, newRecordErr := record.NewRecord(string(j.delivery.Body))
	if newRecordErr == nil {
		loadID := "Load"
		if j.withInfo {
			var flags int64 = 0
			_, withInfoErr := (*j.engine).AddRecordWithInfo(ctx, record.DataSource, record.Id, record.Json, loadID, flags)
			if withInfoErr != nil {
				return fmt.Errorf("add record error, record id: %s, message id: %s, %v", j.delivery.MessageId, record.Id, withInfoErr)
			}
			//TODO:  what do we do with the "withInfo" data here?
			// fmt.Printf("Record added: %s:%s:%s:%s\n", j.delivery.MessageId, loadID, record.DataSource, record.Id)
			// fmt.Printf("WithInfo: %s\n", withInfo)
		} else {
			addRecordErr := (*j.engine).AddRecord(ctx, record.DataSource, record.Id, record.Json, loadID)
			if addRecordErr != nil {
				return fmt.Errorf("add record error, record id: %s, message id: %s, %v", record.Id, j.delivery.MessageId, addRecordErr)
			}
		}

		// when we successfully process a delivery, acknowledge it.
		return j.delivery.Ack(false)
	} else {
		// when we get an invalid delivery, negatively acknowledge and send to the dead letter queue
		err := j.delivery.Nack(false, false)
		return fmt.Errorf("invalid deliver from RabbitMQ, message id: %s, %v", j.delivery.MessageId, err)
	}
}

// ----------------------------------------------------------------------------

// Whenever Execute() returns an error or panics, this is called
func (j *RabbitConsumerJob) OnError(err error) {
	if j.delivery.Redelivered {
		//swallow any error, it'll timeout and be redelivered
		_ = j.delivery.Nack(false, false)
	} else {
		//swallow any error, it'll timeout and be redelivered
		_ = j.delivery.Nack(false, true)
	}
}

// ----------------------------------------------------------------------------
// -- add records in RabbitMQ to Senzing
// ----------------------------------------------------------------------------

// Starts a number of workers that read Records from the given queue and add
// them to Senzing.
// - Workers restart when they are killed or die.
// - respond to standard system signals.
func StartManagedConsumer(ctx context.Context, urlString string, numberOfWorkers int, g2engine *g2api.G2engine, withInfo bool, logLevel string, jsonOutput bool) error {

	//default to the max number of OS threads
	if numberOfWorkers <= 0 {
		numberOfWorkers = runtime.GOMAXPROCS(0)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := SetLogLevel(ctx, logLevel); err != nil {
		log(3003, logLevel, err)
	}
	logger = getLogger()

	log(2012, numberOfWorkers)

	// setup jobs that will be used to process RabbitMQ deliveries
	jobPool = make(chan *RabbitConsumerJob, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		jobPool <- &RabbitConsumerJob{
			engine:    g2engine,
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
			log(2010, jobCount)
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
		log(2011, job.id, job.usedCount)
	}

	return nil
}
