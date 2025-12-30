package config

import (
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Addr                 string
	DatabasePath         string
	Environment          string
	AnthropicAPIKey      string
	AutoApproveThreshold float64
	PriceImportToken     string // Secret token required to access price import feature
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Addr:                 getEnv("ADDR", ":8080"),
		DatabasePath:         getEnv("DATABASE_PATH", "quotes.db"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		AnthropicAPIKey:      getEnv("ANTHROPIC_API_KEY", ""),
		AutoApproveThreshold: getEnvFloat("AUTO_APPROVE_THRESHOLD", 0.9),
		PriceImportToken:     getEnv("PRICE_IMPORT_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
