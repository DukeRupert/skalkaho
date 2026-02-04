package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration.
type Config struct {
	Addr                 string
	DatabaseURL          string // PostgreSQL connection string
	DatabasePath         string // SQLite path (deprecated, for backward compatibility)
	Environment          string
	AnthropicAPIKey      string
	AutoApproveThreshold float64
	PriceImportToken     string        // Secret token required to access price import feature
	SessionSecret        string        // Secret for session token generation
	SessionDuration      time.Duration // Session duration (default: 720h = 30 days)
	SessionCookieName    string        // Session cookie name (default: skalkaho_session)
	SecureCookies        bool          // true if ENVIRONMENT=production
}

// Load reads configuration from environment variables.
func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")
	sessionDuration := getEnvDuration("SESSION_DURATION", 720*time.Hour) // 30 days default

	return &Config{
		Addr:                 getEnv("ADDR", ":8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		DatabasePath:         getEnv("DATABASE_PATH", "quotes.db"),
		Environment:          env,
		AnthropicAPIKey:      getEnv("ANTHROPIC_API_KEY", ""),
		AutoApproveThreshold: getEnvFloat("AUTO_APPROVE_THRESHOLD", 0.9),
		PriceImportToken:     getEnv("PRICE_IMPORT_TOKEN", ""),
		SessionSecret:        getEnv("SESSION_SECRET", ""),
		SessionDuration:      sessionDuration,
		SessionCookieName:    getEnv("SESSION_COOKIE_NAME", "skalkaho_session"),
		SecureCookies:        env == "production",
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

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
