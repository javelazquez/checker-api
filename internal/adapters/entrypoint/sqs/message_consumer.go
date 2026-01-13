package sqs

import (
	"context"
	"fmt"
	"log"

	"checker-api/internal/ports/input/entrypoint"
	sqsclient "checker-api/pkg/sqs"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSMessageConsumer implements the MessageConsumer port for consuming messages from AWS SQS
type SQSMessageConsumer struct {
	sqsClient *sqsclient.Client
}

// NewSQSMessageConsumer creates a new instance of the SQS message consumer
// If sqsClient is nil, it will create a new SQS client using environment configuration
// Returns nil if SQS is not configured (SQS_QUEUE_URL is not set)
func NewSQSMessageConsumer(sqsClient *sqsclient.Client) (entrypoint.MessageConsumer, error) {
	if sqsClient == nil {
		ctx := context.Background()
		client, err := sqsclient.NewClient(ctx, nil)
		if err != nil {
			// If SQS is not configured, return nil (consumer will not start)
			// This allows the application to run without SQS if not needed
			return nil, nil
		}
		sqsClient = client
	}

	// Check if queue URL is configured
	if sqsClient.GetConfig().QueueURL == "" {
		return nil, nil
	}

	return &SQSMessageConsumer{
		sqsClient: sqsClient,
	}, nil
}

// Consume starts consuming messages from the SQS queue
// It runs in a continuous loop until the context is cancelled
// For each message received:
//   - Calls the handler with the message body
//   - If handler succeeds, deletes the message from the queue
//   - If handler fails, the message will become visible again after visibility timeout
func (c *SQSMessageConsumer) Consume(ctx context.Context, handler entrypoint.MessageHandler) error {
	if c.sqsClient == nil {
		return fmt.Errorf("SQS client is not configured")
	}

	log.Println("Starting SQS message consumption...")

	// Continuous loop to receive and process messages
	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			log.Println("SQS consumer stopped: context cancelled")
			return ctx.Err()
		default:
		}

		// Receive messages from SQS
		messages, err := c.sqsClient.ReceiveMessages(ctx)
		if err != nil {
			log.Printf("Error receiving messages from SQS: %v", err)
			// Continue to next iteration to retry
			continue
		}

		// Process each message
		for _, message := range messages {
			// Check if context is cancelled before processing
			select {
			case <-ctx.Done():
				log.Println("SQS consumer stopped: context cancelled")
				return ctx.Err()
			default:
			}

			// Process message
			if err := c.processMessage(ctx, message, handler); err != nil {
				messageID := "unknown"
				if message.MessageId != nil && *message.MessageId != "" {
					messageID = *message.MessageId
				}
				log.Printf("Error processing message (ID: %s): %v", messageID, err)
				// Don't delete the message, let it become visible again after visibility timeout
				// This allows for retry logic
				continue
			}

			// Delete message from queue after successful processing
			if message.ReceiptHandle != nil {
				if err := c.sqsClient.DeleteMessage(ctx, *message.ReceiptHandle); err != nil {
					log.Printf("Error deleting message from SQS: %v", err)
					// Continue processing other messages even if delete fails
				}
			}
		}
	}
}

// processMessage processes a single SQS message
func (c *SQSMessageConsumer) processMessage(ctx context.Context, message types.Message, handler entrypoint.MessageHandler) error {
	if message.Body == nil {
		return fmt.Errorf("message body is nil")
	}

	// Convert message body to byte slice
	rawMessage := []byte(*message.Body)

	// Call the handler
	if err := handler(ctx, rawMessage); err != nil {
		return fmt.Errorf("handler failed: %w", err)
	}

	return nil
}
