package model

// JSONBody is a custom type that can unmarshal JSON objects/arrays/values into []byte
// and marshal []byte back to JSON
type JSONBody []byte

// UnmarshalJSON implements json.Unmarshaler
// Accepts any JSON value (object, array, string, number, boolean, null) and converts it to []byte
func (jb *JSONBody) UnmarshalJSON(data []byte) error {
	// Store the raw JSON bytes
	*jb = make([]byte, len(data))
	copy(*jb, data)
	return nil
}

// MarshalJSON implements json.Marshaler
// Returns the stored JSON bytes as-is
func (jb JSONBody) MarshalJSON() ([]byte, error) {
	if jb == nil {
		return []byte("null"), nil
	}
	return jb, nil
}

// Bytes returns the body as []byte
func (jb JSONBody) Bytes() []byte {
	return []byte(jb)
}

// CompareRequest represents the HTTP request for comparing responses
// @Description Request payload for comparing two API responses
type CompareRequest struct {
	// Source response to compare
	Source Response `json:"source"`
	// Destination response to compare against
	Destination Response `json:"destination"`
}

// Response represents an API response in the HTTP request
// @Description API response structure containing status, headers, and body
type Response struct {
	// Response body as JSON (can be object, array, or any JSON value)
	Body JSONBody `json:"body"`
	// HTTP headers
	Headers map[string]string `json:"headers"`
	// HTTP status code
	Status int `json:"status"`
}
