package kvs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config contains the configuration for the DynamoDB KVS client
// Environment variables are automatically loaded using the "KVS" prefix
type Config struct {
	// Region is the AWS region where the DynamoDB table is located
	// Environment variable: KVS_REGION (default: "us-east-1")
	Region string `envconfig:"REGION" default:"us-east-1"`

	// TableName is the DynamoDB table name
	// Environment variable: KVS_TABLE_NAME (required)
	TableName string `envconfig:"TABLE_NAME" required:"true"`

	// KeyAttributeName is the name of the primary key attribute (default: "id")
	// Environment variable: KVS_KEY_ATTRIBUTE_NAME (default: "id")
	KeyAttributeName string `envconfig:"KEY_ATTRIBUTE_NAME" default:"id"`
}

// NewConfig creates a new configuration reading values from environment variables
// Uses the "KVS" prefix for environment variables
// Expected environment variables:
//   - KVS_REGION: AWS region (default: "us-east-1")
//   - KVS_TABLE_NAME: DynamoDB table name (required)
//   - KVS_KEY_ATTRIBUTE_NAME: Primary key attribute name (default: "id")
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("kvs", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load KVS config from environment: %w", err)
	}

	// Validate configuration
	if cfg.TableName == "" {
		return nil, fmt.Errorf("table name is required")
	}

	if cfg.KeyAttributeName == "" {
		cfg.KeyAttributeName = "id"
	}

	return &cfg, nil
}
