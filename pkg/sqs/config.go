package sqs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config contains the configuration for the SQS client
// Environment variables are automatically loaded using the "SQS" prefix
type Config struct {
	// Region is the AWS region where the SQS queue is located
	// Environment variable: SQS_REGION (default: "us-east-1")
	Region string `envconfig:"REGION" default:"us-east-1"`

	// QueueURL is the SQS queue URL
	// Environment variable: SQS_QUEUE_URL (required)
	QueueURL string `envconfig:"QUEUE_URL" required:"true"`

	// MaxNumberOfMessages is the maximum number of messages to receive per request (1-10)
	// Environment variable: SQS_MAX_NUMBER_OF_MESSAGES (default: 10)
	MaxNumberOfMessages int32 `envconfig:"MAX_NUMBER_OF_MESSAGES" default:"10"`

	// WaitTimeSeconds is the wait time for long polling (0-20 seconds)
	// Environment variable: SQS_WAIT_TIME_SECONDS (default: 20)
	WaitTimeSeconds int32 `envconfig:"WAIT_TIME_SECONDS" default:"20"`

	// VisibilityTimeout is the message visibility timeout in seconds
	// Environment variable: SQS_VISIBILITY_TIMEOUT (default: 30)
	VisibilityTimeout int32 `envconfig:"VISIBILITY_TIMEOUT" default:"30"`
}

// NewConfig creates a new configuration reading values from environment variables
// Uses the "SQS" prefix for environment variables
// Expected environment variables:
//   - SQS_REGION: AWS region (default: "us-east-1")
//   - SQS_QUEUE_URL: SQS queue URL (required)
//   - SQS_MAX_NUMBER_OF_MESSAGES: Maximum number of messages per request (default: 10)
//   - SQS_WAIT_TIME_SECONDS: Wait time for long polling (default: 20)
//   - SQS_VISIBILITY_TIMEOUT: Message visibility timeout (default: 30)
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("sqs", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load SQS config from environment: %w", err)
	}

	// Validate and adjust values
	if cfg.MaxNumberOfMessages < 1 || cfg.MaxNumberOfMessages > 10 {
		cfg.MaxNumberOfMessages = 10
	}
	if cfg.WaitTimeSeconds < 0 || cfg.WaitTimeSeconds > 20 {
		cfg.WaitTimeSeconds = 20
	}
	if cfg.VisibilityTimeout < 1 {
		cfg.VisibilityTimeout = 30
	}

	return &cfg, nil
}
