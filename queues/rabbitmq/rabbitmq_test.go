package rabbitmq_test

import (
	"context"
	"reflect"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-queueing/queues"
	"github.com/senzing-garage/go-queueing/queues/rabbitmq"
	"github.com/stretchr/testify/require"
)

func TestNewClient(test *testing.T) {
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
			actual := rabbitmq.Init(testCase.client, testCase.urlString)
			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("Init() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_Push(test *testing.T) {
	ctx := context.TODO()
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
			err := testCase.client.Push(ctx, testCase.record)
			require.NoError(test, err)
		})
	}
}

func TestClient_UnsafePush(test *testing.T) {
	ctx := context.TODO()
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
			err := testCase.client.UnsafePush(ctx, testCase.args.record)
			require.Equal(test, testCase.expectedErr, err)
		})
	}
}

func TestClient_Consume(test *testing.T) {
	type args struct {
		prefetch int
	}
	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		args        args
		expected    <-chan amqp.Delivery
		expectedErr error
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			actual, actualErr := testCase.client.Consume(testCase.args.prefetch)
			require.Equal(test, testCase.expectedErr, actualErr)
			if !reflect.DeepEqual(actual, testCase.expected) {
				test.Errorf("Client.Consume() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}

func TestClient_Close(test *testing.T) {
	tests := []struct {
		name        string
		client      *rabbitmq.ClientRabbitMQ
		expectedErr error
	}{
		// IMPROVE: Add test cases.
	}
	for _, testCase := range tests {
		test.Run(testCase.name, func(test *testing.T) {
			actualErr := testCase.client.Close()
			require.Equal(test, testCase.expectedErr, actualErr)
		})
	}
}
