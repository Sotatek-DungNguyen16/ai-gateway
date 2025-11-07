package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the AI Gateway
type Config struct {
	Port            string
	APIKeys         []string
	GoogleAPIKey    string
	OpenAIAPIKey    string
	AnthropicAPIKey string
	MaxDiffSize     int64 // Maximum diff size in bytes
	DefaultProvider string
	DefaultModel    string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		APIKeys:         parseAPIKeys(getEnv("API_KEYS", "")),
		GoogleAPIKey:    getEnv("GOOGLE_API_KEY", ""),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		MaxDiffSize:     10 * 1024 * 1024, // 10MB default
		DefaultProvider: getEnv("DEFAULT_AI_PROVIDER", "google"),
		DefaultModel:    getEnv("DEFAULT_AI_MODEL", "gemini-2.0-flash"),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.APIKeys) == 0 {
		return fmt.Errorf("at least one API key must be configured via API_KEYS environment variable")
	}

	if c.GoogleAPIKey == "" && c.OpenAIAPIKey == "" && c.AnthropicAPIKey == "" {
		return fmt.Errorf("at least one AI provider API key must be configured")
	}

	return nil
}

// parseAPIKeys splits comma-separated API keys
func parseAPIKeys(keys string) []string {
	if keys == "" {
		return []string{}
	}

	parts := strings.Split(keys, ",")
	result := make([]string, 0, len(parts))

	for _, key := range parts {
		trimmed := strings.TrimSpace(key)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
