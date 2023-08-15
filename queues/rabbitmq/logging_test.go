package rabbitmq

import (
	"context"
	"reflect"
	"testing"

	"github.com/senzing/go-logging/logging"
)

func Test_getLogger(t *testing.T) {
	tests := []struct {
		name string
		want logging.LoggingInterface
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_log(t *testing.T) {
	type args struct {
		messageNumber int
		details       []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log(tt.args.messageNumber, tt.args.details...)
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	type args struct {
		ctx          context.Context
		logLevelName string
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
			if err := SetLogLevel(tt.args.ctx, tt.args.logLevelName); (err != nil) != tt.wantErr {
				t.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
