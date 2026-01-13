package utility

// IDGenerator generates unique IDs for comparisons
// Defines the contract for generating unique identifiers
// This is the contract that adapters must implement
type IDGenerator interface {
	// Generate generates a unique ID
	Generate() string
}
