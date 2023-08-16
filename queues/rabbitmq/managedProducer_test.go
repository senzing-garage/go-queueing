package rabbitmq

import (
	"context"
	"testing"

	"github.com/senzing/go-queueing/queues"
)

func Test_processRecord(t *testing.T) {
	type args struct {
		ctx         context.Context
		record      queues.Record
		newClientFn func() (*Client, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processRecord(tt.args.ctx, tt.args.record, tt.args.newClientFn); (err != nil) != tt.wantErr {
				t.Errorf("processRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartManagedProducer(t *testing.T) {
	type args struct {
		ctx             context.Context
		urlString       string
		numberOfWorkers int
		recordchan      <-chan queues.Record
		logLevel        string
		jsonOutput      bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartManagedProducer(tt.args.ctx, tt.args.urlString, tt.args.numberOfWorkers, tt.args.recordchan, tt.args.logLevel, tt.args.jsonOutput)
		})
	}
}

func Test_createClients(t *testing.T) {
	type args struct {
		ctx          context.Context
		numOfClients int
		newClientFn  func() (*Client, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createClients(tt.args.ctx, tt.args.numOfClients, tt.args.newClientFn); (err != nil) != tt.wantErr {
				t.Errorf("createClients() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
