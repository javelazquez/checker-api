package utility

import (
	"checker-api/internal/domain/entities"
	outpututility "checker-api/internal/ports/output/utility"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyMessage indicates that the message is empty
	ErrEmptyMessage = errors.New("message cannot be empty")
	// ErrInvalidJSON indicates that the JSON format is invalid
	ErrInvalidJSON = errors.New("invalid JSON format")
	// ErrMissingSource indicates that the source field is missing
	ErrMissingSource = errors.New("field 'source' is required")
	// ErrMissingDestination indicates that the destination field is missing
	ErrMissingDestination = errors.New("field 'destination' is required")
	// ErrInvalidStatusCode indicates that the HTTP status code is invalid
	ErrInvalidStatusCode = errors.New("HTTP status code must be in the range 100-599")
	// ErrInvalidBodyJSON indicates that the body field is not valid JSON
	ErrInvalidBodyJSON = errors.New("body must be valid JSON")
)

// JSONMessageParser implements the MessageParser port for parsing JSON messages
type JSONMessageParser struct{}

// NewJSONMessageParser creates a new instance of the JSON message parser
func NewJSONMessageParser() outpututility.MessageParser {
	return &JSONMessageParser{}
}

// Parse converts a raw message into a domain Message
// Validates JSON structure and returns descriptive errors
func (p *JSONMessageParser) Parse(rawMessage []byte) (*entities.Message, error) {
	// Validate that the message is not empty
	if len(rawMessage) == 0 {
		return nil, ErrEmptyMessage
	}

	// Validate that it's valid JSON (must be an object, not an array or other type)
	trimmed := strings.TrimSpace(string(rawMessage))
	if trimmed == "" || !strings.HasPrefix(trimmed, "{") {
		return nil, fmt.Errorf("%w: expected a JSON object", ErrInvalidJSON)
	}

	var jsonData map[string]json.RawMessage
	if err := json.Unmarshal(rawMessage, &jsonData); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	// Validate that source and destination fields exist
	if _, exists := jsonData["source"]; !exists {
		return nil, ErrMissingSource
	}

	if _, exists := jsonData["destination"]; !exists {
		return nil, ErrMissingDestination
	}

	// Deserialize the complete structure
	// Use json.RawMessage for Body to handle JSON objects/arrays/values properly
	var message struct {
		Source struct {
			Body    json.RawMessage   `json:"body"`
			Headers map[string]string `json:"headers"`
			Status  int               `json:"status"`
		} `json:"source"`
		Destination struct {
			Body    json.RawMessage   `json:"body"`
			Headers map[string]string `json:"headers"`
			Status  int               `json:"status"`
		} `json:"destination"`
	}

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		return nil, fmt.Errorf("%w: error deserializing message structure: %v", ErrInvalidJSON, err)
	}

	return &entities.Message{
		Source: entities.Response{
			Body:    message.Source.Body,
			Headers: message.Source.Headers,
			Status:  message.Source.Status,
		},
		Destination: entities.Response{
			Body:    message.Destination.Body,
			Headers: message.Destination.Headers,
			Status:  message.Destination.Status,
		},
	}, nil
}

// Validate validates that a parsed message meets the requirements
// Validates that source and destination are not empty, status codes are valid,
// and required fields are present
func (p *JSONMessageParser) Validate(message *entities.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	// Validate source
	if err := p.validateResponse(&message.Source, "source"); err != nil {
		return err
	}

	// Validate destination
	if err := p.validateResponse(&message.Destination, "destination"); err != nil {
		return err
	}

	return nil
}

// validateResponse validates that a response has the required and valid fields
func (p *JSONMessageParser) validateResponse(response *entities.Response, fieldName string) error {
	if response == nil {
		return fmt.Errorf("field '%s' cannot be nil", fieldName)
	}

	// Validate status code (valid HTTP codes are in the range 100-599)
	if response.Status < 100 || response.Status > 599 {
		return fmt.Errorf("%w in field '%s': status code %d", ErrInvalidStatusCode, fieldName, response.Status)
	}

	// Validate that headers is not nil (can be an empty map but not nil)
	if response.Headers == nil {
		return fmt.Errorf("field 'headers' in '%s' cannot be nil (can be an empty object)", fieldName)
	}

	// Body can be empty, but must be valid JSON (not nil)
	if response.Body == nil {
		return fmt.Errorf("field 'body' in '%s' cannot be nil (can be an empty JSON object or array)", fieldName)
	}

	// Validate that body is valid JSON
	if len(response.Body) > 0 {
		var jsonValue interface{}
		if err := json.Unmarshal(response.Body, &jsonValue); err != nil {
			return fmt.Errorf("%w in field '%s': %v", ErrInvalidBodyJSON, fieldName, err)
		}
	}

	return nil
}
