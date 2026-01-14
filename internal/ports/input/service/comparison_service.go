package service

import (
	"context"
	"checker-api/internal/domain/entities"
)

// ComparisonService defines the input port for comparison operations
// This is the contract that input adapters (SQS, HTTP, etc.) must use
// Difference: ComparisonService includes context.Context (application interface),
// while domain ports do not include it
type ComparisonService interface {
	// ProcessMessage processes a received message:
	// 1. Parses and validates the message
	// 2. Compares the responses
	// 3. Calculates the differences
	// 4. Persists the result
	ProcessMessage(ctx context.Context, rawMessage []byte) (*entities.Comparison, error)

	// GetByID retrieves a comparison by its ID from the repository
	GetByID(ctx context.Context, id string) (*entities.Comparison, error)
}
