package rabbitmq

import (
	"context"
	"testing"

	"github.com/senzing-garage/g2-sdk-go/g2api"
)

func TestRabbitConsumerJob_Execute(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *RabbitConsumerJob
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Execute(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RabbitConsumerJob.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRabbitConsumerJob_OnError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		j    *RabbitConsumerJob
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.OnError(tt.args.err)
		})
	}
}

func TestStartManagedConsumer(t *testing.T) {
	type args struct {
		ctx             context.Context
		urlString       string
		numberOfWorkers int
		g2engine        *g2api.G2engine
		withInfo        bool
		logLevel        string
		jsonOutput      bool
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
			if err := StartManagedConsumer(tt.args.ctx, tt.args.urlString, tt.args.numberOfWorkers, tt.args.g2engine, tt.args.withInfo, tt.args.logLevel, tt.args.jsonOutput); (err != nil) != tt.wantErr {
				t.Errorf("StartManagedConsumer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
