package model

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
	// Response body as bytes (base64 encoded in JSON)
	Body []byte `json:"body"`
	// HTTP headers
	Headers map[string]string `json:"headers"`
	// HTTP status code
	Status int `json:"status"`
}
