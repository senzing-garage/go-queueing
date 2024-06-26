package rabbitmq

//lint:file-ignore U1000 Ignore all unused code, it's a test file

import (
	"bufio"
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/senzing-garage/go-logging/logging"
)

func Test_getLogger(test *testing.T) {
	tests := []struct {
		name string
		want logging.Logging
	}{
		{name: "Test non-nil logger", want: getLogger()},
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if got := getLogger(); !reflect.DeepEqual(got, tt.want) {
				test.Errorf("getLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_log(test *testing.T) {

	scanner, cleanUpStdout := mockStdout(test)
	defer cleanUpStdout()

	type args struct {
		messageNumber int
		details       []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test log", args: args{messageNumber: 2001, details: []interface{}{"RabbitMQ"}}, want: "RabbitMQ"},
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			log(tt.args.messageNumber, tt.args.details...)
			got := ""
			for i := 0; i < 1; i++ {
				scanner.Scan()
				got += scanner.Text()
			}
			if !strings.Contains(got, tt.want) {
				test.Errorf("getLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_log_JSON_output(test *testing.T) {

	type args struct {
		messageNumber int
		details       []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test log JSON output", args: args{messageNumber: 2001, details: []interface{}{"RabbitMQ"}}, want: "RabbitMQ"},
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			_ = test
			ouputJSON = true
			//TODO: intercept logging and test that it is JSON
			// log(tt.args.messageNumber, tt.args.details...)
			// if !strings.Contains(got, tt.want) {
			// 	t.Errorf("getLogger() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestSetLogLevel(test *testing.T) {
	type args struct {
		ctx          context.Context
		logLevelName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test SetLogLevel", args: args{ctx: context.Background(), logLevelName: "DEBUG"}, want: "DEBUG", wantErr: false},
	}
	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			if err := SetLogLevel(tt.args.ctx, tt.args.logLevelName); (err != nil) != tt.wantErr {
				test.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if level := getLogger().GetLogLevel(); level != tt.want {
				test.Errorf("SetLogLevel() error = %v, wantErr %v", level, tt.wantErr)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

// capture stdout for testing
func mockStdout(test *testing.T) (buffer *bufio.Scanner, cleanUp func()) {
	test.Helper()
	origStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		test.Fatalf("couldn't get os Pipe: %v", err)
	}
	os.Stdout = writer

	return bufio.NewScanner(reader),
		func() {
			// clean-up
			os.Stdout = origStdout
		}
}

// capture stderr for testing
func mockStderr(test *testing.T) (buffer *bufio.Scanner, sync func(), cleanUp func()) {
	test.Helper()
	origStderr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		test.Fatalf("couldn't get os Pipe: %v", err)
	}
	os.Stderr = writer

	return bufio.NewScanner(reader),
		func() {
			err = writer.Sync()
			if err != nil {
				panic(err)
			}
		},
		func() {
			// clean-up
			os.Stderr = origStderr
		}
}
