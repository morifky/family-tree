package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

// Config holds all the environment-based configuration for the application.
type Config struct {
	Port         string `validate:"required,numeric"`
	DatabasePath string `validate:"required"`
	PhotosDir    string `validate:"required"`
	LogLevel     string `validate:"required,oneof=debug info warn error"`
	// Additional config fields can be added here
}

// MustLoad reads the environment variables, validates them, and returns a Config struct.
// It panics if any required environment variable is missing or invalid.
func MustLoad() *Config {
	cfg := &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		DatabasePath: getEnvOrDefault("DATABASE_PATH", "/data/brayat.db"),
		PhotosDir:    getEnvOrDefault("PHOTOS_DIR", "/data/photos"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(fmt.Errorf("configuration validation failed: %w", err))
	}

	return cfg
}

// getEnvOrDefault returns the value of an environment variable, or a fallback value if it is not set or empty.
func getEnvOrDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
