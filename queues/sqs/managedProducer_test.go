package sqs

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-queueing/queues"
)

func Test_processRecordBatch(test *testing.T) {
	type args struct {
		ctx         context.Context
		recordchan  <-chan queues.Record
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
		test.Run(tt.name, func(test *testing.T) {
			if err := processRecordBatch(tt.args.ctx, tt.args.recordchan, tt.args.newClientFn); (err != nil) != tt.wantErr {
				test.Errorf("processRecordBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartManagedProducer(test *testing.T) {
	type args struct {
		ctx             context.Context
		urlString       string
		numberOfWorkers int
		recordchan      <-chan queues.Record
		logLevel        string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			StartManagedProducer(tt.args.ctx, tt.args.urlString, tt.args.numberOfWorkers, tt.args.recordchan, tt.args.logLevel)
		})
	}
}

func Test_createClients(test *testing.T) {
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
		test.Run(tt.name, func(test *testing.T) {
			if err := createClients(tt.args.ctx, tt.args.numOfClients, tt.args.newClientFn); (err != nil) != tt.wantErr {
				test.Errorf("createClients() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
