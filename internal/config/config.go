package config

import (
	"os"
	"time"
)

// Config holds application configuration.
type Config struct {
	Addr              string
	DatabaseURL       string        // PostgreSQL connection string (required)
	Environment       string        // development or production
	SessionSecret     string        // Secret for session token generation
	SessionDuration   time.Duration // Session duration (default: 720h = 30 days)
	SessionCookieName string        // Session cookie name (default: skalkaho_session)
	SecureCookies     bool          // true if ENVIRONMENT=production
	PostmarkAPIKey    string        // Postmark transactional email API key
}

// Load reads configuration from environment variables.
func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")
	sessionDuration := getEnvDuration("SESSION_DURATION", 720*time.Hour)

	return &Config{
		Addr:              getEnv("ADDR", ":8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		Environment:       env,
		SessionSecret:     getEnv("SESSION_SECRET", ""),
		SessionDuration:   sessionDuration,
		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "skalkaho_session"),
		SecureCookies:     env == "production",
		PostmarkAPIKey:    getEnv("POSTMARK_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
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
