package sqs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Client is the client for interacting with AWS SQS
type Client struct {
	sqsClient *sqs.Client
	config    *Config
}

// NewClient creates a new instance of the SQS client
// Uses default AWS configuration (environment variables, credential files, etc.)
// If cfg is nil, loads configuration from environment variables using envconfig
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	var err error
	if cfg == nil {
		cfg, err = NewConfig()
		if err != nil {
			return nil, err
		}
	}

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Support custom endpoint (e.g., LocalStack)
	// Use BaseEndpoint option when creating the client (non-deprecated approach)
	var sqsClientOptions []func(*sqs.Options)
	if endpointURL := os.Getenv("AWS_ENDPOINT_URL"); endpointURL != "" {
		sqsClientOptions = append(sqsClientOptions, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
		})
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(awsCfg, sqsClientOptions...)

	return &Client{
		sqsClient: sqsClient,
		config:    cfg,
	}, nil
}

// ReceiveMessages receives messages from the SQS queue
// Returns the received messages and any error
func (c *Client) ReceiveMessages(ctx context.Context) ([]types.Message, error) {
	if c.config.QueueURL == "" {
		return nil, fmt.Errorf("queue URL is required")
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.config.QueueURL),
		MaxNumberOfMessages: c.config.MaxNumberOfMessages,
		WaitTimeSeconds:     c.config.WaitTimeSeconds,
		VisibilityTimeout:   c.config.VisibilityTimeout,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
		MessageAttributeNames: []string{
			"All",
		},
	}

	result, err := c.sqsClient.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	return result.Messages, nil
}

// DeleteMessage deletes a message from the SQS queue after successfully processing it
func (c *Client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	if c.config.QueueURL == "" {
		return fmt.Errorf("queue URL is required")
	}

	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.config.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.sqsClient.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// SendMessage sends a message to the SQS queue
func (c *Client) SendMessage(ctx context.Context, messageBody string) error {
	if c.config.QueueURL == "" {
		return fmt.Errorf("queue URL is required")
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.config.QueueURL),
		MessageBody: aws.String(messageBody),
	}

	_, err := c.sqsClient.SendMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// GetQueueURL gets the queue URL by its name
func (c *Client) GetQueueURL(ctx context.Context, queueName string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := c.sqsClient.GetQueueUrl(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get queue URL: %w", err)
	}

	return *result.QueueUrl, nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}
