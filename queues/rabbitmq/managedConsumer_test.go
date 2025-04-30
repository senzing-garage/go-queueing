package rabbitmq_test

import (
	"testing"

	"github.com/senzing-garage/go-queueing/queues/rabbitmq"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
)

func TestRabbitConsumerJob_Execute(test *testing.T) {
	tests := []struct {
		name        string
		expectedErr bool
		job         *rabbitmq.RabbitConsumerJob
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			ctx := test.Context()
			err := testCase.job.Execute(ctx)
			require.NoError(test, err)
		})
	}
}

func TestRabbitConsumerJob_OnError(test *testing.T) {
	tests := []struct {
		name string
		err  error
		job  *rabbitmq.RabbitConsumerJob
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			_ = test
			testCase.job.OnError(testCase.err)
		})
	}
}

func TestStartManagedConsumer(test *testing.T) {
	tests := []struct {
		name            string
		expectedErr     bool
		jsonOutput      bool
		logLevel        string
		numberOfWorkers int
		szEngine        *senzing.SzEngine
		urlString       string
		withInfo        bool
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			ctx := test.Context()
			err := rabbitmq.StartManagedConsumer(
				ctx,
				testCase.urlString,
				testCase.numberOfWorkers,
				testCase.szEngine,
				testCase.withInfo,
				testCase.logLevel,
				testCase.jsonOutput,
			)
			require.NoError(test, err)
		})
	}
}
