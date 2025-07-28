// File 2: internal/config/config.go
package config

import (
	"os"
	"time"
)

type Config struct {
	Environment        string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiry          time.Duration
	RedisURL           string
	PlaidClientID      string
	PlaidSecret        string
	OpenAIKey          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	return &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost/finance_db?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiry:     getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
		PlaidClientID: getEnv("PLAID_CLIENT_ID", ""),
		PlaidSecret:   getEnv("PLAID_SECRET", ""),
		OpenAIKey:     getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
