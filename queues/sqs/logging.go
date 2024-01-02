package sqs

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// logging variables.
var logger logging.LoggingInterface
var ouputJSON bool

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func getLogger() logging.LoggingInterface {
	var err error = nil
	if logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		logger, err = logging.NewSenzingToolsLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

// Log message.
func log(messageNumber int, details ...interface{}) {
	if ouputJSON {
		getLogger().Log(messageNumber, details...)
	} else {
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...))
	}
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = getLogger().SetLogLevel(logLevelName)
	return err
}
