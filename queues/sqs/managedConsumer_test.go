package sqs_test

import (
	"testing"

	"github.com/senzing-garage/go-queueing/queues/sqs"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
)

func TestSQSJob_Execute(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name              string
		expectedErr       error
		job               *sqs.ConsumerJobSqs
		visibilitySeconds int32
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			expectedErr := testCase.job.Execute(ctx, testCase.visibilitySeconds)
			require.Equal(test, testCase.expectedErr, expectedErr)
		})
	}
}

func TestSQSJob_OnError(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		expectedErr error
		job         *sqs.ConsumerJobSqs
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			testCase.job.OnError(ctx, testCase.expectedErr)
		})
	}
}

func TestStartManagedConsumer(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name              string
		expectedErr       bool
		jsonOutput        bool
		logLevel          string
		numberOfWorkers   int
		szEngine          senzing.SzEngine
		urlString         string
		visibilitySeconds int32
		withInfo          bool
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actualErr := sqs.StartManagedConsumer(
				ctx,
				testCase.urlString,
				testCase.numberOfWorkers,
				testCase.szEngine,
				testCase.withInfo,
				testCase.visibilitySeconds,
				testCase.logLevel,
				testCase.jsonOutput,
			)
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}
