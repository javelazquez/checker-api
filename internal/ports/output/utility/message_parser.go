package utility

import "checker-api/internal/domain/entities"

// MessageParser defines the output port for parsing and validating messages
// Defines the contract for converting raw messages into domain structures
// This is the contract that adapters must implement
type MessageParser interface {
	// Parse converts a raw message into a domain Message
	// Returns error if the message is not valid or cannot be parsed
	Parse(rawMessage []byte) (*entities.Message, error)

	// Validate validates that a parsed message meets the requirements
	// Returns error if the message is not valid
	Validate(message *entities.Message) error
}
