package usecases

import (
	"encoding/json"
	"fmt"
	"strings"

	"checker-api/internal/domain"
	"checker-api/internal/domain/entities"
)

// CompareResponsesUseCase implements the use case for comparing responses
// Implements domain ports: ResponseComparator and DifferenceCalculator
type CompareResponsesUseCase struct{}

// NewCompareResponsesUseCase creates a new instance of the use case
func NewCompareResponsesUseCase() (domain.ResponseComparator, domain.DifferenceCalculator) {
	uc := &CompareResponsesUseCase{}
	return uc, uc
}

// Compare compares two responses and determines if they are equal
// Implements domain.ResponseComparator
// Compares status codes, headers (case-insensitive), and JSON bodies (ignoring field order)
func (uc *CompareResponsesUseCase) Compare(source, destination entities.Response) (bool, error) {
	// Compare status codes
	if source.Status != destination.Status {
		return false, nil
	}

	// Compare headers (case-insensitive)
	if !uc.compareHeaders(source.Headers, destination.Headers) {
		return false, nil
	}

	// Compare JSON bodies
	equal, err := uc.compareJSONBodies(source.Body, destination.Body)
	if err != nil {
		return false, fmt.Errorf("error comparing JSON bodies: %w", err)
	}

	return equal, nil
}

// CalculateDifferences identifies all differences between two responses
// Implements domain.DifferenceCalculator
// Returns a detailed list of all differences found in status, headers, and JSON body
func (uc *CompareResponsesUseCase) CalculateDifferences(source, destination entities.Response) ([]entities.Difference, error) {
	var differences []entities.Difference

	// Compare status codes
	if source.Status != destination.Status {
		differences = append(differences, entities.Difference{
			Type:        entities.DifferenceTypeStatus,
			Path:        "status",
			SourceValue: source.Status,
			DestValue:   destination.Status,
			Description: fmt.Sprintf("Status code differs: source=%d, destination=%d", source.Status, destination.Status),
		})
	}

	// Compare headers
	headerDiffs := uc.calculateHeaderDifferences(source.Headers, destination.Headers)
	differences = append(differences, headerDiffs...)

	// Compare JSON bodies
	bodyDiffs, err := uc.calculateJSONBodyDifferences(source.Body, destination.Body)
	if err != nil {
		return nil, fmt.Errorf("error calculating JSON body differences: %w", err)
	}
	differences = append(differences, bodyDiffs...)

	return differences, nil
}

// compareHeaders compares two header maps case-insensitively
func (uc *CompareResponsesUseCase) compareHeaders(source, destination map[string]string) bool {
	if len(source) != len(destination) {
		return false
	}

	// Normalize headers to lowercase keys for comparison
	normalizedSource := uc.normalizeHeaders(source)
	normalizedDest := uc.normalizeHeaders(destination)

	for key, sourceValue := range normalizedSource {
		destValue, exists := normalizedDest[key]
		if !exists || sourceValue != destValue {
			return false
		}
	}

	return true
}

// normalizeHeaders converts all header keys to lowercase
func (uc *CompareResponsesUseCase) normalizeHeaders(headers map[string]string) map[string]string {
	normalized := make(map[string]string, len(headers))
	for k, v := range headers {
		normalized[strings.ToLower(k)] = v
	}
	return normalized
}

// compareJSONBodies compares two JSON byte arrays, ignoring field order
func (uc *CompareResponsesUseCase) compareJSONBodies(sourceBody, destBody []byte) (bool, error) {
	// Handle empty bodies
	if len(sourceBody) == 0 && len(destBody) == 0 {
		return true, nil
	}
	if len(sourceBody) == 0 || len(destBody) == 0 {
		return false, nil
	}

	// Parse JSON
	var sourceJSON, destJSON interface{}
	if err := json.Unmarshal(sourceBody, &sourceJSON); err != nil {
		return false, fmt.Errorf("invalid JSON in source body: %w", err)
	}
	if err := json.Unmarshal(destBody, &destJSON); err != nil {
		return false, fmt.Errorf("invalid JSON in destination body: %w", err)
	}

	// Deep compare
	return uc.deepEqualJSON(sourceJSON, destJSON), nil
}

// deepEqualJSON performs a deep comparison of two JSON values
func (uc *CompareResponsesUseCase) deepEqualJSON(a, b interface{}) bool {
	// Use JSON marshaling to normalize and compare (this handles field order differences)
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}

	// Unmarshal both to normalized format
	var aNorm, bNorm interface{}
	if err := json.Unmarshal(aBytes, &aNorm); err != nil {
		return false
	}
	if err := json.Unmarshal(bBytes, &bNorm); err != nil {
		return false
	}

	return uc.deepEqual(aNorm, bNorm)
}

// deepEqual performs recursive deep comparison
func (uc *CompareResponsesUseCase) deepEqual(a, b interface{}) bool {
	// Type check
	if reflectType(a) != reflectType(b) {
		return false
	}

	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for key, aValue := range aVal {
			bValue, exists := bVal[key]
			if !exists || !uc.deepEqual(aValue, bValue) {
				return false
			}
		}
		return true

	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !uc.deepEqual(aVal[i], bVal[i]) {
				return false
			}
		}
		return true

	default:
		return a == b
	}
}

// reflectType returns a string representation of the type for comparison
func reflectType(v interface{}) string {
	switch v.(type) {
	case nil:
		return "nil"
	case bool:
		return "bool"
	case float64:
		return "float64"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return fmt.Sprintf("%T", v)
	}
}

// calculateHeaderDifferences calculates all differences between two header maps
func (uc *CompareResponsesUseCase) calculateHeaderDifferences(source, destination map[string]string) []entities.Difference {
	var differences []entities.Difference

	normalizedSource := uc.normalizeHeaders(source)
	normalizedDest := uc.normalizeHeaders(destination)

	// Find missing and different headers in destination
	for key, sourceValue := range normalizedSource {
		destValue, exists := normalizedDest[key]
		if !exists {
			differences = append(differences, entities.Difference{
				Type:        entities.DifferenceTypeMissing,
				Path:        fmt.Sprintf("headers.%s", key),
				SourceValue: sourceValue,
				DestValue:   nil,
				Description: fmt.Sprintf("Header '%s' is missing in destination", key),
			})
		} else if sourceValue != destValue {
			differences = append(differences, entities.Difference{
				Type:        entities.DifferenceTypeHeader,
				Path:        fmt.Sprintf("headers.%s", key),
				SourceValue: sourceValue,
				DestValue:   destValue,
				Description: fmt.Sprintf("Header '%s' value differs: source='%s', destination='%s'", key, sourceValue, destValue),
			})
		}
	}

	// Find extra headers in destination
	for key, destValue := range normalizedDest {
		if _, exists := normalizedSource[key]; !exists {
			differences = append(differences, entities.Difference{
				Type:        entities.DifferenceTypeExtra,
				Path:        fmt.Sprintf("headers.%s", key),
				SourceValue: nil,
				DestValue:   destValue,
				Description: fmt.Sprintf("Header '%s' is extra in destination", key),
			})
		}
	}

	return differences
}

// calculateJSONBodyDifferences calculates all differences between two JSON bodies
func (uc *CompareResponsesUseCase) calculateJSONBodyDifferences(sourceBody, destBody []byte) ([]entities.Difference, error) {
	var differences []entities.Difference

	// Handle empty bodies
	if len(sourceBody) == 0 && len(destBody) == 0 {
		return differences, nil
	}
	if len(sourceBody) == 0 {
		return []entities.Difference{{
			Type:        entities.DifferenceTypeBody,
			Path:        "body",
			SourceValue: nil,
			DestValue:   string(destBody),
			Description: "Source body is empty but destination is not",
		}}, nil
	}
	if len(destBody) == 0 {
		return []entities.Difference{{
			Type:        entities.DifferenceTypeBody,
			Path:        "body",
			SourceValue: string(sourceBody),
			DestValue:   nil,
			Description: "Destination body is empty but source is not",
		}}, nil
	}

	// Parse JSON
	var sourceJSON, destJSON interface{}
	if err := json.Unmarshal(sourceBody, &sourceJSON); err != nil {
		return nil, fmt.Errorf("invalid JSON in source body: %w", err)
	}
	if err := json.Unmarshal(destBody, &destJSON); err != nil {
		return nil, fmt.Errorf("invalid JSON in destination body: %w", err)
	}

	// Calculate differences recursively
	bodyDiffs := uc.calculateJSONDifferences(sourceJSON, destJSON, "body")
	return bodyDiffs, nil
}

// calculateJSONDifferences recursively calculates differences between two JSON structures
func (uc *CompareResponsesUseCase) calculateJSONDifferences(source, destination interface{}, path string) []entities.Difference {
	var differences []entities.Difference

	sourceType := reflectType(source)
	destType := reflectType(destination)

	// Type mismatch
	if sourceType != destType {
		differences = append(differences, entities.Difference{
			Type:        entities.DifferenceTypeBody,
			Path:        path,
			SourceValue: source,
			DestValue:   destination,
			Description: fmt.Sprintf("Type mismatch at '%s': source is %s, destination is %s", path, sourceType, destType),
		})
		return differences
	}

	// Compare based on type
	switch sourceVal := source.(type) {
	case map[string]interface{}:
		destVal, ok := destination.(map[string]interface{})
		if !ok {
			return differences
		}

		// Find missing and different keys
		for key, sourceValue := range sourceVal {
			newPath := fmt.Sprintf("%s.%s", path, key)
			destValue, exists := destVal[key]
			if !exists {
				differences = append(differences, entities.Difference{
					Type:        entities.DifferenceTypeMissing,
					Path:        newPath,
					SourceValue: sourceValue,
					DestValue:   nil,
					Description: fmt.Sprintf("Field '%s' is missing in destination", newPath),
				})
			} else {
				nestedDiffs := uc.calculateJSONDifferences(sourceValue, destValue, newPath)
				differences = append(differences, nestedDiffs...)
			}
		}

		// Find extra keys
		for key, destValue := range destVal {
			if _, exists := sourceVal[key]; !exists {
				newPath := fmt.Sprintf("%s.%s", path, key)
				differences = append(differences, entities.Difference{
					Type:        entities.DifferenceTypeExtra,
					Path:        newPath,
					SourceValue: nil,
					DestValue:   destValue,
					Description: fmt.Sprintf("Field '%s' is extra in destination", newPath),
				})
			}
		}

	case []interface{}:
		destVal, ok := destination.([]interface{})
		if !ok {
			return differences
		}

		// Compare arrays element by element
		maxLen := len(sourceVal)
		if len(destVal) > maxLen {
			maxLen = len(destVal)
		}

		for i := 0; i < maxLen; i++ {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			if i >= len(sourceVal) {
				differences = append(differences, entities.Difference{
					Type:        entities.DifferenceTypeExtra,
					Path:        newPath,
					SourceValue: nil,
					DestValue:   destVal[i],
					Description: fmt.Sprintf("Array element at index %d is extra in destination", i),
				})
			} else if i >= len(destVal) {
				differences = append(differences, entities.Difference{
					Type:        entities.DifferenceTypeMissing,
					Path:        newPath,
					SourceValue: sourceVal[i],
					DestValue:   nil,
					Description: fmt.Sprintf("Array element at index %d is missing in destination", i),
				})
			} else {
				nestedDiffs := uc.calculateJSONDifferences(sourceVal[i], destVal[i], newPath)
				differences = append(differences, nestedDiffs...)
			}
		}

	default:
		// Primitive values
		if source != destination {
			sourceStr := uc.valueToString(source)
			destStr := uc.valueToString(destination)
			differences = append(differences, entities.Difference{
				Type:        entities.DifferenceTypeBody,
				Path:        path,
				SourceValue: source,
				DestValue:   destination,
				Description: fmt.Sprintf("Value differs at '%s': source=%s, destination=%s", path, sourceStr, destStr),
			})
		}
	}

	return differences
}

// valueToString converts a value to string representation
func (uc *CompareResponsesUseCase) valueToString(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case float64:
		// JSON numbers are unmarshaled as float64
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
