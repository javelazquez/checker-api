package http

import (
	"checker-api/internal/adapters/entrypoint/http/handlers"
	"checker-api/internal/ports/input/service"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Router configures the HTTP routes of the application
type Router struct {
	comparisonHandler *handlers.ComparisonHandler
}

// NewRouter creates a new instance of the router
func NewRouter(comparisonService service.ComparisonService) *Router {
	return &Router{
		comparisonHandler: handlers.NewComparisonHandler(comparisonService),
	}
}

// SetupRoutes configures all HTTP routes
// TODO: Add middleware (logging, recovery, CORS, etc.)
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Route for comparing responses
	mux.HandleFunc("/api/v1/compare", r.comparisonHandler.Compare)

	// Route for getting a comparison by ID
	mux.HandleFunc("/api/v1/compare/", r.comparisonHandler.GetByID)

	// Swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	// TODO: Add health check route
	// mux.HandleFunc("/health", r.healthHandler.Health)

	return mux
}
