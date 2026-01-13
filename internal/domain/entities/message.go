package entities

// Message represents a message received from SQS
type Message struct {
	Source      Response
	Destination Response
}

// Response represents an API response
type Response struct {
	Body    []byte
	Headers map[string]string
	Status  int
}

// Comparison represents the result of comparing two responses
type Comparison struct {
	ID          string
	Source      Response
	Destination Response
	Differences []Difference
	IsEqual     bool
}

// Difference represents a specific difference found between two responses
type Difference struct {
	Type        DifferenceType
	Path        string
	SourceValue interface{}
	DestValue   interface{}
	Description string
}

// DifferenceType defines the type of difference found
type DifferenceType string

const (
	DifferenceTypeStatus   DifferenceType = "status"
	DifferenceTypeHeader   DifferenceType = "header"
	DifferenceTypeBody     DifferenceType = "body"
	DifferenceTypeMissing  DifferenceType = "missing"
	DifferenceTypeExtra    DifferenceType = "extra"
)
