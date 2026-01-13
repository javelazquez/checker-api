package kvs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client is the client for interacting with AWS DynamoDB as a Key-Value Store
type Client struct {
	dynamoClient *dynamodb.Client
	config       *Config
}

// NewClient creates a new instance of the DynamoDB KVS client
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
	var dynamoClientOptions []func(*dynamodb.Options)
	if endpointURL := os.Getenv("AWS_ENDPOINT_URL"); endpointURL != "" {
		dynamoClientOptions = append(dynamoClientOptions, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
		})
	}

	// Create DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(awsCfg, dynamoClientOptions...)

	return &Client{
		dynamoClient: dynamoClient,
		config:       cfg,
	}, nil
}

// Put stores a value in DynamoDB with the given key
// The value will be serialized to JSON and stored as a string attribute
func (c *Client) Put(ctx context.Context, key string, value interface{}) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Serialize value to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Create item with key and value
	item := map[string]types.AttributeValue{
		c.config.KeyAttributeName: &types.AttributeValueMemberS{
			Value: key,
		},
		"value": &types.AttributeValueMemberS{
			Value: string(valueJSON),
		},
	}

	// Put item in DynamoDB
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.config.TableName),
		Item:      item,
	}

	_, err = c.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %w", err)
	}

	return nil
}

// Get retrieves a value from DynamoDB by key
// The value will be deserialized from JSON into the provided target
func (c *Client) Get(ctx context.Context, key string, target interface{}) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Get item from DynamoDB
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.config.TableName),
		Key: map[string]types.AttributeValue{
			c.config.KeyAttributeName: &types.AttributeValueMemberS{
				Value: key,
			},
		},
	}

	result, err := c.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	// Check if item exists
	if result.Item == nil {
		return fmt.Errorf("item not found: %s", key)
	}

	// Extract value attribute
	valueAttr, ok := result.Item["value"]
	if !ok {
		return fmt.Errorf("value attribute not found in item")
	}

	valueStr, ok := valueAttr.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("value attribute is not a string")
	}

	// Deserialize JSON value
	if err := json.Unmarshal([]byte(valueStr.Value), target); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from DynamoDB by key
func (c *Client) Delete(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Delete item from DynamoDB
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.config.TableName),
		Key: map[string]types.AttributeValue{
			c.config.KeyAttributeName: &types.AttributeValueMemberS{
				Value: key,
			},
		},
	}

	_, err := c.dynamoClient.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %w", err)
	}

	return nil
}

// Exists checks if a key exists in DynamoDB
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}

	// Get item from DynamoDB
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.config.TableName),
		Key: map[string]types.AttributeValue{
			c.config.KeyAttributeName: &types.AttributeValueMemberS{
				Value: key,
			},
		},
		// Only return the key attribute to minimize data transfer
		ProjectionExpression: aws.String(c.config.KeyAttributeName),
	}

	result, err := c.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return false, fmt.Errorf("failed to check item existence in DynamoDB: %w", err)
	}

	return result.Item != nil, nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}
