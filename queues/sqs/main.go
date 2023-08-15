package sqs

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
	2005: Prefix + "SQS client is setup! QueueURL: %v, Queue name: %v, DeadLetterQueueURL: %v",
	2006: Prefix + "AWS response Message ID: %s",
	2007: Prefix + "Successfully sent %v records to the queue",
	2008: Prefix + "SQS Client delete message: %v",
	2009: Prefix + "SQS Client set message visibility, MessageID: %v",
	2010: Prefix + "Jobs added to job queue:: %v",
	2011: Prefix + "Job ID: %v, Job count: %d",
	2012: Prefix + "Number of consumer workers: %d",

	// WARN 	3000-3999 	Unexpected situations, but processing was successful
	// ERROR 	4000-4999 	Unexpected situations, processing was not successful
	4006: Prefix + "Unable to retrieve queue redrive policy, error: %w",
	4007: Prefix + "Error unmarshalling redrive policy, error: %w",
	4008: Prefix + "Error sending to the dead record queue, error: %w",
	4009: Prefix + "Error sending the record, error: %w",
	4010: Prefix + "Error sending the record batch, error: %w",
	4011: Prefix + "Error sending record in batch, MessageID: %v, error: %v",
	4012: Prefix + "Error sending the last record batch, error: %w",
	4013: Prefix + "Error receiving records, error: %v",
	4014: Prefix + "No records found.",
	4015: Prefix + "Error deleting records, error: %v",
	4016: Prefix + "Error changing message visibility, MessageID: %v, error: %v",
	4017: Prefix + "Error getting delivery channel, error: %v",

	// FATAL 	5000-5999 	The process needs to shutdown
	// PANIC 	6000-6999 	The underlying system is at issue
	//
	//	8000-8999 	Reserved for observer messages
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}
