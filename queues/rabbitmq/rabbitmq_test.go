package rabbitmq

import (
	"context"
	"reflect"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-queueing/queues"
)

func TestNewClient(test *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			got, err := NewClient(tt.args.urlString)
			if (err != nil) != tt.wantErr {
				test.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				test.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(test *testing.T) {
	type args struct {
		client    *Client
		urlString string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if got := Init(tt.args.client, tt.args.urlString); !reflect.DeepEqual(got, tt.want) {
				test.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleReconnect(test *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name   string
		client *Client
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			tt.client.handleReconnect(tt.args.addr)
		})
	}
}

func TestClient_connect(test *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *amqp.Connection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			got, err := tt.client.connect(tt.args.addr)
			if (err != nil) != tt.wantErr {
				test.Errorf("Client.connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				test.Errorf("Client.connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleReInit(test *testing.T) {
	type args struct {
		conn *amqp.Connection
	}
	tests := []struct {
		name   string
		client *Client
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if got := tt.client.handleReInit(tt.args.conn); got != tt.want {
				test.Errorf("Client.handleReInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_init(test *testing.T) {
	type args struct {
		conn *amqp.Connection
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := tt.client.init(tt.args.conn); (err != nil) != tt.wantErr {
				test.Errorf("Client.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_changeConnection(test *testing.T) {
	type args struct {
		connection *amqp.Connection
	}
	tests := []struct {
		name   string
		client *Client
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			tt.client.changeConnection(tt.args.connection)
		})
	}
}

func TestClient_changeChannel(test *testing.T) {
	type args struct {
		channel *amqp.Channel
	}
	tests := []struct {
		name   string
		client *Client
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			tt.client.changeChannel(tt.args.channel)
		})
	}
}

func TestClient_progressiveDelay(test *testing.T) {
	type args struct {
		delay time.Duration
	}
	tests := []struct {
		name   string
		client *Client
		args   args
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if got := tt.client.progressiveDelay(tt.args.delay); !reflect.DeepEqual(got, tt.want) {
				test.Errorf("Client.progressiveDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Push(test *testing.T) {
	ctx := context.TODO()
	type args struct {
		record queues.Record
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := tt.client.Push(ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UnsafePush(test *testing.T) {
	ctx := context.TODO()
	type args struct {
		record queues.Record
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := tt.client.UnsafePush(ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.UnsafePush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Consume(test *testing.T) {
	type args struct {
		prefetch int
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    <-chan amqp.Delivery
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			got, err := tt.client.Consume(tt.args.prefetch)
			if (err != nil) != tt.wantErr {
				test.Errorf("Client.Consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				test.Errorf("Client.Consume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Close(test *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := tt.client.Close(); (err != nil) != tt.wantErr {
				test.Errorf("Client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
