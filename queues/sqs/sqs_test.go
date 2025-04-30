package sqs_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/senzing-garage/go-queueing/queues"
	"github.com/senzing-garage/go-queueing/queues/sqs"
	"github.com/stretchr/testify/require"
)

func TestNewClient(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		expected    *sqs.ClientSqs
		expectedErr error
		jsonOutput  bool
		logLevel    string
		urlString   string
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actual, actualErr := sqs.NewClient(
				ctx,
				testCase.urlString,
				testCase.logLevel,
				testCase.jsonOutput,
			)
			require.Equal(test, testCase.expectedErr, actualErr)

			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("NewClient() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_PushDeadRecord(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *sqs.ClientSqs
		expectedErr error
		record      types.Message
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			expectedErr := testCase.client.PushDeadRecord(ctx, testCase.record)
			require.Equal(test, testCase.expectedErr, expectedErr)
		})
	}
}

func TestClient_Push(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *sqs.ClientSqs
		expectedErr error
		record      queues.Record
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actualErr := testCase.client.Push(ctx, testCase.record)
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}

func TestClient_PushBatch(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name          string
		client        *sqs.ClientSqs
		expectedErr   error
		recordChannel <-chan queues.Record
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actualErr := testCase.client.PushBatch(ctx, testCase.recordChannel)
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}

func TestClient_Consume(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name              string
		client            *sqs.ClientSqs
		expected          <-chan types.Message
		expectedErr       error
		visibilitySeconds int32
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actual, actualErr := testCase.client.Consume(ctx, testCase.visibilitySeconds)
			require.Equal(test, testCase.expectedErr, actualErr)

			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("Client.Consume() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_RemoveMessage(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *sqs.ClientSqs
		expectedErr error
		message     types.Message
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			actualErr := testCase.client.RemoveMessage(ctx, testCase.message)
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}

func TestClient_SetMessageVisibility(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *sqs.ClientSqs
		expectedErr bool
		message     types.Message
		seconds     int32
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()
			ctx := test.Context()
			err := testCase.client.SetMessageVisibility(ctx, testCase.message, testCase.seconds)
			require.Equal(test, testCase.expectedErr, err)
		})
	}
}

func TestClient_Close(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *sqs.ClientSqs
		expectedErr error
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			expectedErr := testCase.client.Close()
			require.Equal(test, testCase.expectedErr, expectedErr)
		})
	}
}
