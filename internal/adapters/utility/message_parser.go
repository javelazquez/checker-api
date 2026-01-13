package utility

import (
	"checker-api/internal/domain/entities"
	outpututility "checker-api/internal/ports/output/utility"
	"encoding/json"
)

// JSONMessageParser implements the MessageParser port for parsing JSON messages
// TODO: Add more robust validations and detailed error handling
type JSONMessageParser struct{}

// NewJSONMessageParser creates a new instance of the JSON message parser
func NewJSONMessageParser() outpututility.MessageParser {
	return &JSONMessageParser{}
}

// Parse converts a raw message into a domain Message
// TODO: Improve implementation:
// 1. Validate JSON structure
// 2. Validate that source and destination exist
// 3. Validate data types
// 4. Add error logging
// 5. Handle different message formats if necessary
func (p *JSONMessageParser) Parse(rawMessage []byte) (*entities.Message, error) {
	var message struct {
		Source      entities.Response `json:"source"`
		Destination entities.Response `json:"destination"`
	}

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		// TODO: Return more descriptive error
		return nil, err
	}

	return &entities.Message{
		Source:      message.Source,
		Destination: message.Destination,
	}, nil
}

// Validate validates that a parsed message meets the requirements
// TODO: Implement validations:
// 1. Validate that source and destination are not empty
// 2. Validate that required fields are present
// 3. Validate data formats (valid status codes, etc.)
// 4. Return descriptive errors for each failed validation
func (p *JSONMessageParser) Validate(message *entities.Message) error {
	// TODO: Implement message validations
	return nil
}
