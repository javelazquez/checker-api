// @title           Checker API
// @version         1.0
// @description     API for comparing API responses and detecting differences
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "checker-api/docs" // Swagger documentation - import for side effects
	"checker-api/internal/adapters/persistence"
	"checker-api/internal/adapters/utility"
	"checker-api/internal/application/services"
	"checker-api/internal/application/usecases"
	"checker-api/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if environment is LOCAL
	loadEnvIfLocal()

	ctx := context.Background()

	// Initialize output adapters (repositories)
	messageParser := utility.NewJSONMessageParser()
	comparisonRepo, err := persistence.NewDynamoDBComparisonRepository(nil)
	if err != nil {
		log.Fatalf("Failed to create comparison repository: %v", err)
	}
	idGenerator := utility.NewUUIDIDGenerator()

	// Initialize use cases (application layer)
	responseComparator, differenceCalculator := usecases.NewCompareResponsesUseCase()

	// Initialize application service
	comparisonService := services.NewComparisonService(
		messageParser,
		responseComparator,
		differenceCalculator,
		comparisonRepo,
		idGenerator,
	)

	// Create server with configuration from environment variables
	srv, err := server.NewServer(nil, comparisonService)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-shutdown

	// Perform graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
}

// loadEnvIfLocal loads environment variables from .env file
// only if the environment is LOCAL
func loadEnvIfLocal() {
	env := strings.ToLower(os.Getenv("ENV"))
	if env == "" {
		env = strings.ToLower(os.Getenv("ENVIRONMENT"))
	}

	// Load .env only if environment is LOCAL or not defined (local development)
	if env == "local" || env == "" {
		if err := godotenv.Load(); err != nil {
			// Not critical if .env file doesn't exist
			log.Printf("Warning: Could not load .env file: %v (continuing without it)", err)
		} else {
			log.Println("Environment variables loaded from .env (LOCAL environment)")
		}
	}
}
