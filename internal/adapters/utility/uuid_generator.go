package utility

import (
	outpututility "checker-api/internal/ports/output/utility"
	"checker-api/pkg/generatorid"
)

// UUIDIDGenerator implements the IDGenerator port using UUIDs
// Uses the generatorid client from the pkg package
type UUIDIDGenerator struct {
	client *generatorid.Client
}

// NewUUIDIDGenerator creates a new instance of the UUID ID generator
func NewUUIDIDGenerator() outpututility.IDGenerator {
	return &UUIDIDGenerator{
		client: generatorid.NewClient(),
	}
}

// Generate generates a unique ID using UUID v4 (random)
func (g *UUIDIDGenerator) Generate() string {
	return g.client.Generate()
}
