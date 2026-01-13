package entrypoint

import "context"

// MessageConsumer defines the input port for consuming messages from SQS
// Defines the contract for receiving messages from the queue
// This is the contract that input adapters (SQS) must use
type MessageConsumer interface {
	// Consume starts consuming messages from the queue
	// The handler will be called for each received message
	Consume(ctx context.Context, handler MessageHandler) error
}

// MessageHandler processes a received message
type MessageHandler func(ctx context.Context, rawMessage []byte) error
