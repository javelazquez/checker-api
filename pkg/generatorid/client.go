package generatorid

import (
	"github.com/google/uuid"
)

// Client is the client for generating unique IDs using UUIDs
type Client struct {
	// Future configuration if needed
}

// NewClient creates a new instance of the ID generator client
func NewClient() *Client {
	return &Client{}
}

// Generate generates a unique ID using UUID v4 (random)
// Returns a string with the generated UUID
func (c *Client) Generate() string {
	return uuid.New().String()
}

// GenerateWithPrefix generates a unique ID with an optional prefix
// Useful for adding context to the ID (e.g., "comparison-{uuid}")
func (c *Client) GenerateWithPrefix(prefix string) string {
	if prefix == "" {
		return c.Generate()
	}
	return prefix + "-" + uuid.New().String()
}
