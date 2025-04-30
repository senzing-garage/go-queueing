package rabbitmq_test

import (
	"reflect"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-queueing/queues"
	"github.com/senzing-garage/go-queueing/queues/rabbitmq"
	"github.com/stretchr/testify/require"
)

func TestNewClient(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		expected    *rabbitmq.ClientRabbitMQ
		expectedErr bool
		urlString   string
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			actual, err := rabbitmq.NewClient(testCase.urlString)
			require.NoError(test, err)
			require.Equal(test, testCase.expected, actual)

			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("NewClient() = %v, expected: %v", actual, testCase.expected)
			}
		})
	}
}

func TestInit(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name      string
		client    *rabbitmq.ClientRabbitMQ
		expected  *rabbitmq.ClientRabbitMQ
		urlString string
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			actual := rabbitmq.Init(testCase.client, testCase.urlString)
			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("Init() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_Push(test *testing.T) {
	test.Parallel()
	ctx := test.Context()

	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		expectedErr bool
		record      queues.Record
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			err := testCase.client.Push(ctx, testCase.record)
			require.NoError(test, err)
		})
	}
}

func TestClient_UnsafePush(test *testing.T) {
	test.Parallel()
	ctx := test.Context()

	type args struct {
		record queues.Record
	}

	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		args        args
		expectedErr error
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			err := testCase.client.UnsafePush(ctx, testCase.args.record)
			require.Equal(test, testCase.expectedErr, err)
		})
	}
}

func TestClient_Consume(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		expected    <-chan amqp.Delivery
		expectedErr error
		prefetch    int
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			actual, actualErr := testCase.client.Consume(testCase.prefetch)
			require.Equal(test, testCase.expectedErr, actualErr)

			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("Client.Consume() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_Close(test *testing.T) {
	test.Parallel()

	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		expectedErr error
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			actualErr := testCase.client.Close()
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}
