package model

// CompareResponse represents the HTTP response of a comparison
// @Description Response containing the comparison result with differences
type CompareResponse struct {
	// Unique identifier for the comparison
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Whether the responses are equal
	IsEqual bool `json:"is_equal" example:"false"`
	// List of differences found between the responses
	Differences []Difference `json:"differences"`
}

// Difference represents a difference in the HTTP response
// @Description Details about a specific difference found between responses
type Difference struct {
	// Type of difference (status, header, body, missing, extra)
	Type string `json:"type" example:"status"`
	// Path to the difference (e.g., "body.user.name", "headers.Content-Type")
	Path string `json:"path" example:"status"`
	// Value in the source response (can be any type: string, number, boolean, etc.)
	SourceValue interface{} `json:"source_value"`
	// Value in the destination response (can be any type: string, number, boolean, etc.)
	DestValue interface{} `json:"dest_value"`
	// Human-readable description of the difference
	Description string `json:"description" example:"Status codes differ: 200 vs 404"`
}

// ErrorResponse represents an HTTP error response
// @Description Error response structure
type ErrorResponse struct {
	// Error type or code
	Error string `json:"error" example:"Invalid request body"`
	// Detailed error message
	Message string `json:"message,omitempty" example:"Failed to parse JSON"`
}
