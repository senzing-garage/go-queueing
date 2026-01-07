# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-queueing provides queue adapters for RabbitMQ and AWS SQS, designed for use with Senzing entity resolution. It includes managed producers and consumers that handle connection pooling, retries with progressive backoff, and worker orchestration.

## Build and Development Commands

```bash
# Install development dependencies (one-time)
make dependencies-for-development

# Update Go dependencies
make dependencies

# Build for current platform
make build

# Build for all platforms
make build-all

# Run linting (golangci-lint + govulncheck + cspell)
make lint

# Run tests
make test

# Run a single test
go test -v -run TestFunctionName ./path/to/package

# Run tests with coverage
make coverage

# Check test coverage against threshold
make check-coverage

# Apply lint fixers
make fix
```

## Architecture

### Core Interface

`queues/main.go` defines the `Record` interface that all queue implementations use:

- `GetMessage() string` - Returns the message body
- `GetMessageID() string` - Returns the unique message identifier

### Queue Implementations

Each queue adapter (`rabbitmq/`, `sqs/`) follows the same structure:

**Client** (`rabbitmq.go`, `sqs.go`):

- Handles connection lifecycle with automatic reconnection
- Provides `Push()`, `Consume()`, and `Close()` methods
- Uses progressive delay for retries (increases delay on failures, resets on success)

**ManagedProducer** (`managedProducer.go`):

- Starts worker pool to push records from a channel to the queue
- Maintains a client pool for concurrent publishing
- Workers automatically replace failed clients

**ManagedConsumer** (`managedConsumer.go`):

- Starts worker pool to consume messages and process them via Senzing SDK
- Uses job pool pattern - reuses job structs to minimize allocations
- Handles acknowledgment/negative-acknowledgment based on processing success
- Dead letter queue support for failed messages

### Key Dependencies

- `github.com/sourcegraph/conc/pool` - Concurrent worker pool management
- `github.com/rabbitmq/amqp091-go` - RabbitMQ client
- `github.com/aws/aws-sdk-go-v2` - AWS SQS client
- `github.com/senzing-garage/sz-sdk-go` - Senzing SDK integration

### URL Connection Strings

**RabbitMQ**: `amqp://user:pass@host:port/vhost?exchange=name&queue-name=name&routing-key=key`

**SQS**: Either a full SQS URL or `sqs://...?queue-name=name`

## Code Conventions

- Component ID for logging: `6481`
- Log levels follow ranges: TRACE (0-999), DEBUG (1000-1999), INFO (2000-2999), WARN (3000-3999), ERROR (4000-4999)
- Maximum line length: 120 characters
- Tests use `test.Parallel()` and stretchr/testify
- Error wrapping uses `github.com/senzing-garage/go-helpers/wraperror`
