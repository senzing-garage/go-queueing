package sqs_test

import (
	"testing"

	"github.com/senzing-garage/go-queueing/queues"
	"github.com/senzing-garage/go-queueing/queues/sqs"
)

func TestStartManagedProducer(test *testing.T) {
	tests := []struct {
		name            string
		jsonOutput      bool
		logLevel        string
		numberOfWorkers int
		recordchan      <-chan queues.Record
		urlString       string
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
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
