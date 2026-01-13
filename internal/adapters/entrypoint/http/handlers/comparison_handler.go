package handlers

import (
	"encoding/json"
	"net/http"

	"checker-api/internal/adapters/entrypoint/http/model"
	"checker-api/internal/domain/entities"
	"checker-api/internal/ports/input/service"
)

// ComparisonHandler handles HTTP requests related to comparisons
type ComparisonHandler struct {
	comparisonService service.ComparisonService
}

// NewComparisonHandler creates a new instance of the handler
func NewComparisonHandler(comparisonService service.ComparisonService) *ComparisonHandler {
	return &ComparisonHandler{
		comparisonService: comparisonService,
	}
}

// Compare handles the POST request to compare responses
// @Summary      Compare two API responses
// @Description  Compares two API responses and returns the differences between them
// @Tags         comparisons
// @Accept       json
// @Produce      json
// @Param        request  body      model.CompareRequest  true  "Comparison request with source and destination responses"
// @Success      200      {object}  model.CompareResponse  "Comparison result"
// @Failure      400      {object}  model.ErrorResponse    "Invalid request"
// @Failure      500      {object}  model.ErrorResponse    "Internal server error"
// @Router       /compare [post]
func (h *ComparisonHandler) Compare(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req model.CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Convert request to internal format (rawMessage)
	rawMessage, err := json.Marshal(req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to process request", err.Error())
		return
	}

	// Call application service
	comparison, err := h.comparisonService.ProcessMessage(r.Context(), rawMessage)
	if err != nil {
		// TODO: Improve error handling (different error types, appropriate HTTP codes)
		h.respondError(w, http.StatusInternalServerError, "Failed to process comparison", err.Error())
		return
	}

	// Convert entity to HTTP response
	response := h.toCompareResponse(comparison)

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// TODO: Log error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// toCompareResponse converts a Comparison entity to HTTP CompareResponse
func (h *ComparisonHandler) toCompareResponse(comparison *entities.Comparison) model.CompareResponse {
	differences := make([]model.Difference, len(comparison.Differences))
	for i, diff := range comparison.Differences {
		differences[i] = model.Difference{
			Type:        string(diff.Type),
			Path:        diff.Path,
			SourceValue: diff.SourceValue,
			DestValue:   diff.DestValue,
			Description: diff.Description,
		}
	}

	return model.CompareResponse{
		ID:          comparison.ID,
		IsEqual:     comparison.IsEqual,
		Differences: differences,
	}
}

// respondError responds with an HTTP error
func (h *ComparisonHandler) respondError(w http.ResponseWriter, statusCode int, errorMsg, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error:   errorMsg,
		Message: message,
	})
}
