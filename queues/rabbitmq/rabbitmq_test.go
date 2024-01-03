package rabbitmq

import (
	"reflect"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/senzing-garage/go-queueing/queues"
)

func TestNewClient(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.urlString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.client, tt.args.urlString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleReconnect(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			tt.client.handleReconnect(tt.args.addr)
		})
	}
}

func TestClient_connect(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.connect(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleReInit(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.handleReInit(tt.args.conn); got != tt.want {
				t.Errorf("Client.handleReInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_init(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.init(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("Client.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_changeConnection(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			tt.client.changeConnection(tt.args.connection)
		})
	}
}

func TestClient_changeChannel(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			tt.client.changeChannel(tt.args.channel)
		})
	}
}

func TestClient_progressiveDelay(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.progressiveDelay(tt.args.delay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.progressiveDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Push(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Push(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("Client.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UnsafePush(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.UnsafePush(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("Client.UnsafePush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Consume(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Consume(tt.args.prefetch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Consume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
