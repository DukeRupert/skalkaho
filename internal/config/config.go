package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	Addr         string
	DatabasePath string
	Environment  string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Addr:         getEnv("ADDR", ":8080"),
		DatabasePath: getEnv("DATABASE_PATH", "quotes.db"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
