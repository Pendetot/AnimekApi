package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port           string
	Environment    string
	BaseURL        string
	UserAgent      string
	CacheDuration  time.Duration
	RequestTimeout time.Duration
	RateLimit      RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	WindowMs    time.Duration
	MaxRequests int
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "1408"),
		Environment: getEnv("ENVIRONMENT", "development"),
		BaseURL:     getEnv("BASE_URL", "https://ww1.anoboy.app"),
		UserAgent:   getEnv("USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		CacheDuration: time.Duration(getEnvAsInt("CACHE_DURATION", 30)) * time.Minute,
		RequestTimeout: time.Duration(getEnvAsInt("REQUEST_TIMEOUT", 10)) * time.Second,
		RateLimit: RateLimitConfig{
			WindowMs:    time.Duration(getEnvAsInt("RATE_LIMIT_WINDOW_MS", 15)) * time.Minute,
			MaxRequests: getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 100),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
