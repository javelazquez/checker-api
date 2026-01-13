package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	httpadapter "checker-api/internal/adapters/entrypoint/http"
	sqsadapter "checker-api/internal/adapters/entrypoint/sqs"
	"checker-api/internal/ports/input/entrypoint"
	"checker-api/internal/ports/input/service"
)

// Server represents the application server
type Server struct {
	httpServer        *http.Server
	sqsConsumer       entrypoint.MessageConsumer
	comparisonService service.ComparisonService
	config            *Config
}

// NewServer creates a new instance of the server with dependencies already initialized
func NewServer(cfg *Config, comparisonService service.ComparisonService) (*Server, error) {
	if cfg == nil {
		var err error
		cfg, err = NewConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create server config: %w", err)
		}
	}

	// Setup HTTP server
	httpServer := setupHTTPServer(cfg, comparisonService)

	// Setup SQS consumer (optional)
	sqsConsumer := setupSQSConsumer()

	return &Server{
		httpServer:        httpServer,
		sqsConsumer:       sqsConsumer,
		comparisonService: comparisonService,
		config:            cfg,
	}, nil
}

// Start starts the HTTP server and SQS consumer
func (s *Server) Start(ctx context.Context) error {
	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP server starting on port %s", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	// Start SQS consumer in a goroutine (if configured)
	if s.sqsConsumer != nil {
		go func() {
			log.Println("SQS consumer starting...")
			handler := func(ctx context.Context, rawMessage []byte) error {
				_, err := s.comparisonService.ProcessMessage(ctx, rawMessage)
				if err != nil {
					log.Printf("Error processing SQS message: %v", err)
					return err
				}
				log.Println("SQS message processed successfully")
				return nil
			}

			if err := s.sqsConsumer.Consume(ctx, handler); err != nil {
				log.Printf("SQS consumer error: %v", err)
			}
		}()
	}

	log.Println("Application started successfully")
	log.Println("Press Ctrl+C to shutdown")

	return nil
}

// Shutdown performs graceful shutdown of the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutdown signal received, gracefully shutting down...")

	// Create context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	// Close HTTP server gracefully
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error shutting down HTTP server: %w", err)
	}

	log.Println("HTTP server shut down successfully")
	log.Println("Application shutdown complete")

	return nil
}

// setupHTTPServer configures and returns the HTTP server
func setupHTTPServer(cfg *Config, comparisonService service.ComparisonService) *http.Server {
	router := httpadapter.NewRouter(comparisonService)
	mux := router.SetupRoutes()

	// Add health check route
	mux.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return server
}

// setupSQSConsumer configures the SQS consumer if available
// Returns nil if not configured (SQS_QUEUE_URL is not defined)
func setupSQSConsumer() entrypoint.MessageConsumer {
	// Create SQS consumer (pass nil to use environment configuration)
	consumer, err := sqsadapter.NewSQSMessageConsumer(nil)
	if err != nil {
		log.Printf("Failed to create SQS consumer: %v", err)
		return nil
	}

	if consumer == nil {
		log.Println("SQS consumer not configured (SQS_QUEUE_URL not set)")
		return nil
	}

	log.Println("SQS consumer configured (will start if SQS_QUEUE_URL is set)")
	return consumer
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
