package domain

import (
	"checker-api/internal/domain/entities"
)

// ResponseComparator defines the domain contract for comparing two responses
// This is the main domain port that will be implemented by the application layer
// IMPORTANT: Does not include context.Context because the domain should not depend on frameworks
type ResponseComparator interface {
	// Compare compares two responses and determines if they are equal
	// Returns true if equal, false if different
	Compare(source, destination entities.Response) (bool, error)
}

// DifferenceCalculator defines the domain contract for calculating differences between responses
// This is the main domain port that will be implemented by the application layer
// IMPORTANT: Does not include context.Context because the domain should not depend on frameworks
type DifferenceCalculator interface {
	// CalculateDifferences identifies all differences between two responses
	// Returns a list of found differences
	CalculateDifferences(source, destination entities.Response) ([]entities.Difference, error)
}
