package persistence

import (
	"checker-api/internal/domain/entities"
	"checker-api/internal/ports/output/repository"
	"checker-api/pkg/kvs"
	"context"
	"fmt"
	"strings"
)

// DynamoDBComparisonRepository implements the ComparisonRepository port using AWS DynamoDB
// Uses the KVS client for DynamoDB operations
type DynamoDBComparisonRepository struct {
	kvsClient *kvs.Client
}

// NewDynamoDBComparisonRepository creates a new instance of the DynamoDB repository
// If kvsClient is nil, it will create a new KVS client using environment configuration
func NewDynamoDBComparisonRepository(kvsClient *kvs.Client) (repository.ComparisonRepository, error) {
	if kvsClient == nil {
		ctx := context.Background()
		client, err := kvs.NewClient(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create KVS client: %w", err)
		}
		kvsClient = client
	}

	return &DynamoDBComparisonRepository{
		kvsClient: kvsClient,
	}, nil
}

// Save stores a comparison in DynamoDB
// Uses the comparison ID as the key
func (r *DynamoDBComparisonRepository) Save(ctx context.Context, comparison *entities.Comparison) error {
	if comparison == nil {
		return fmt.Errorf("comparison cannot be nil")
	}

	if comparison.ID == "" {
		return fmt.Errorf("comparison ID cannot be empty")
	}

	// Use KVS Put to store the comparison
	// The KVS client will serialize the comparison to JSON
	if err := r.kvsClient.Put(ctx, comparison.ID, comparison); err != nil {
		return fmt.Errorf("failed to save comparison: %w", err)
	}

	return nil
}

// FindByID retrieves a comparison by its ID from DynamoDB
// Returns an error if the comparison is not found
func (r *DynamoDBComparisonRepository) FindByID(ctx context.Context, id string) (*entities.Comparison, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	var comparison entities.Comparison

	// Use KVS Get to retrieve the comparison
	// The KVS client will deserialize the JSON to the Comparison entity
	if err := r.kvsClient.Get(ctx, id, &comparison); err != nil {
		// Check if error is "item not found"
		if strings.Contains(err.Error(), "item not found") {
			return nil, fmt.Errorf("comparison not found: %s", id)
		}
		return nil, fmt.Errorf("failed to retrieve comparison: %w", err)
	}

	return &comparison, nil
}

// Exists checks if a comparison exists for a given ID in DynamoDB
func (r *DynamoDBComparisonRepository) Exists(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("id cannot be empty")
	}

	// Use KVS Exists to check if the key exists
	exists, err := r.kvsClient.Exists(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check comparison existence: %w", err)
	}

	return exists, nil
}
