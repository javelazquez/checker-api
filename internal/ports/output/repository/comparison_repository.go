package repository

import (
	"context"
	"checker-api/internal/domain/entities"
)

// ComparisonRepository defines the output port for persisting comparisons
// Defines the contract for storing and retrieving comparisons in DynamoDB
// This is the contract that persistence adapters must implement
type ComparisonRepository interface {
	// Save stores a comparison in the repository
	Save(ctx context.Context, comparison *entities.Comparison) error

	// FindByID retrieves a comparison by its ID
	FindByID(ctx context.Context, id string) (*entities.Comparison, error)

	// Exists checks if a comparison exists for a given ID
	Exists(ctx context.Context, id string) (bool, error)
}
