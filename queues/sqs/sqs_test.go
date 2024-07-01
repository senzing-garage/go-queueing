package sqs

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/senzing-garage/go-queueing/queues"
)

func TestNewClient(test *testing.T) {
	type args struct {
		ctx        context.Context
		urlString  string
		logLevel   string
		jsonOutput bool
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
			got, err := NewClient(tt.args.ctx, tt.args.urlString, tt.args.logLevel, tt.args.jsonOutput)
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

func TestClient_getQueueURL(test *testing.T) {
	type args struct {
		ctx       context.Context
		urlString string
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
			if err := tt.client.getQueueURL(tt.args.ctx, tt.args.urlString); (err != nil) != tt.wantErr {
				test.Errorf("Client.getQueueURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_getRedrivePolicy(test *testing.T) {
	type args struct {
		ctx context.Context
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
			tt.client.getRedrivePolicy(tt.args.ctx)
		})
	}
}

func TestClient_sendDeadRecord(test *testing.T) {
	type args struct {
		ctx    context.Context
		record types.Message
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
			if err := tt.client.sendDeadRecord(tt.args.ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.sendDeadRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_sendRecord(test *testing.T) {
	type args struct {
		ctx    context.Context
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
			if err := tt.client.sendRecord(tt.args.ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.sendRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_sendRecordBatch(test *testing.T) {
	type args struct {
		ctx     context.Context
		records []queues.Record
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
			if err := tt.client.sendRecordBatch(tt.args.ctx, tt.args.records); (err != nil) != tt.wantErr {
				test.Errorf("Client.sendRecordBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
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

func TestClient_PushDeadRecord(test *testing.T) {
	type args struct {
		ctx    context.Context
		record types.Message
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
			if err := tt.client.PushDeadRecord(tt.args.ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.PushDeadRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Push(test *testing.T) {
	type args struct {
		ctx    context.Context
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
			if err := tt.client.Push(tt.args.ctx, tt.args.record); (err != nil) != tt.wantErr {
				test.Errorf("Client.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_PushBatch(test *testing.T) {
	type args struct {
		ctx        context.Context
		recordchan <-chan queues.Record
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
			if err := tt.client.PushBatch(tt.args.ctx, tt.args.recordchan); (err != nil) != tt.wantErr {
				test.Errorf("Client.PushBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_receiveMessage(test *testing.T) {
	type args struct {
		ctx               context.Context
		visibilitySeconds int32
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *sqs.ReceiveMessageOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			got, err := tt.client.receiveMessage(tt.args.ctx, tt.args.visibilitySeconds)
			if (err != nil) != tt.wantErr {
				test.Errorf("Client.receiveMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				test.Errorf("Client.receiveMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Consume(test *testing.T) {
	type args struct {
		ctx               context.Context
		visibilitySeconds int32
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    <-chan types.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			got, err := tt.client.Consume(tt.args.ctx, tt.args.visibilitySeconds)
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

func TestClient_RemoveMessage(test *testing.T) {
	type args struct {
		ctx context.Context
		msg types.Message
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
			if err := tt.client.RemoveMessage(tt.args.ctx, tt.args.msg); (err != nil) != tt.wantErr {
				test.Errorf("Client.RemoveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SetMessageVisibility(test *testing.T) {
	type args struct {
		ctx     context.Context
		msg     types.Message
		seconds int32
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
			if err := tt.client.SetMessageVisibility(tt.args.ctx, tt.args.msg, tt.args.seconds); (err != nil) != tt.wantErr {
				test.Errorf("Client.SetMessageVisibility() error = %v, wantErr %v", err, tt.wantErr)
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
