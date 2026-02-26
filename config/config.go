package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	AppName string
	AppEnv  string
	Port    string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", ""),
		AppEnv:  getEnv("APP_ENV", "development"),
		Port:    getEnv("PORT", "8080"),
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets integer environment variable with default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// IsDevelopment checks if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction checks if app is in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
