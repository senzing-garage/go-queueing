package sqs

import (
	"context"
	"testing"

	"github.com/senzing-garage/sz-sdk-go/senzing"
)

func TestSQSJob_Execute(test *testing.T) {
	type args struct {
		ctx               context.Context
		visibilitySeconds int32
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := tt.j.Execute(tt.args.ctx, tt.args.visibilitySeconds); (err != nil) != tt.wantErr {
				test.Errorf("SQSJob.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQSJob_OnError(test *testing.T) {
	ctx := context.TODO()
	type args struct {
		err error
	}
	tests := []struct {
		name string
		j    *Job
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			tt.j.OnError(ctx, tt.args.err)
		})
	}
}

func TestStartManagedConsumer(test *testing.T) {
	_ = test
	type args struct {
		ctx               context.Context
		urlString         string
		numberOfWorkers   int
		szEngine          senzing.SzEngine
		withInfo          bool
		visibilitySeconds int32
		logLevel          string
		jsonOutput        bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			if err := StartManagedConsumer(tt.args.ctx, tt.args.urlString, tt.args.numberOfWorkers, tt.args.szEngine, tt.args.withInfo, tt.args.visibilitySeconds, tt.args.logLevel, tt.args.jsonOutput); (err != nil) != tt.wantErr {
				t.Errorf("StartManagedConsumer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
