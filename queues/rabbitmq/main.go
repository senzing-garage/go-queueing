package rabbitmq

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// go-queueing is 6481:  https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md
const ComponentID = 6481

// Log message prefix.
const Prefix = "rabbitmq: "

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IDMessages = map[int]string{
	// Level 	Range 		Use 							Comments
	// TRACE 	0000-0999 	Entry/Exit tracing 				May contain sensitive data.
	// DEBUG 	1000-1999 	Values seen during processing 	May contain sensitive data.
	// INFO 	2000-2999 	Process steps achieved
	2001: Prefix + "Connecting to RabbitMQ at %s",
	2002: Prefix + "Connected to RabbitMQ at %s",
	2003: Prefix + "Connection closed. Reconnecting...",
	2004: Prefix + "Channel closed. Re-running init...",
	2005: Prefix + "RabbitMQ client setup!",

	// WARN 	3000-3999 	Unexpected situations, but processing was successful
	3001: Prefix + "Push failed. Retrying in %v... MessageId: %v, error: %w",
	3002: Prefix + "Push did not confirm. Retrying in %v... MessageId: %v, error: %w",

	// ERROR 	4000-4999 	Unexpected situations, processing was not successful
	4001: Prefix + "Failed to connect. Retrying in %v, error: %w",
	4002: Prefix + "Failed to initialize. Retrying in %v, error: %w",
	4003: Prefix + "Attempting to init RabbitMQ client while shutting down. error: %w",
	4004: Prefix + "channel close error: %w",
	4005: Prefix + "connection close error: %w",
	// FATAL 	5000-5999 	The process needs to shutdown
	// PANIC 	6000-6999 	The underlying system is at issue
	//
	//	8000-8999 	Reserved for observer messages
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}
