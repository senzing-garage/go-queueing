package sqs_test

import (
	"testing"

	"github.com/senzing-garage/go-queueing/queues"
	"github.com/senzing-garage/go-queueing/queues/sqs"
)

func TestStartManagedProducer(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name            string
		jsonOutput      bool
		logLevel        string
		numberOfWorkers int
		recordchan      <-chan queues.Record
		urlString       string
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			sqs.StartManagedProducer(
				ctx,
				testCase.urlString,
				testCase.numberOfWorkers,
				testCase.recordchan,
				testCase.logLevel,
				testCase.jsonOutput,
			)
		})
	}
}
