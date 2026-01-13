package server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config contains the HTTP server configuration
// Environment variables are automatically loaded using the "SERVER" prefix
type Config struct {
	// Port is the port where the HTTP server will listen
	// Environment variable: SERVER_PORT (default: "8080")
	Port string `envconfig:"PORT" default:"8080"`

	// ReadTimeout is the maximum time to read the full request body
	// Environment variable: SERVER_READ_TIMEOUT (default: "15s")
	ReadTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"15s"`

	// WriteTimeout is the maximum time before canceling the response write
	// Environment variable: SERVER_WRITE_TIMEOUT (default: "15s")
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"15s"`

	// IdleTimeout is the maximum time to wait for the next request when keep-alive is enabled
	// Environment variable: SERVER_IDLE_TIMEOUT (default: "60s")
	IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`

	// ShutdownTimeout is the maximum time to perform graceful shutdown
	// Environment variable: SERVER_SHUTDOWN_TIMEOUT (default: "30s")
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// NewConfig creates a new configuration reading values from environment variables
// Uses the "SERVER" prefix for environment variables
// Expected environment variables:
//   - SERVER_PORT: HTTP server port (default: "8080")
//   - SERVER_READ_TIMEOUT: Read timeout (default: "15s")
//   - SERVER_WRITE_TIMEOUT: Write timeout (default: "15s")
//   - SERVER_IDLE_TIMEOUT: Idle timeout (default: "60s")
//   - SERVER_SHUTDOWN_TIMEOUT: Shutdown timeout (default: "30s")
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("server", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load server config from environment: %w", err)
	}

	return &cfg, nil
}
