package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
// It's called during application startup to load configuration
func LoadEnv() {
	// Load .env file if it exists
	_ = godotenv.Load()
}

// GetEnv retrieves an environment variable with a default value
// Parameters:
//   - key: the environment variable name
//   - defaultValue: the default value if the variable is not set
//
// Returns: the environment variable value or default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// AppConfig holds application configuration
type AppConfig struct {
	Port      string
	Env       string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

// LoadConfig loads application configuration from environment variables
// Returns: AppConfig struct with all configuration values
func LoadConfig() AppConfig {
	return AppConfig{
		Port:      GetEnv("PORT", "8080"),
		Env:       GetEnv("ENV", "development"),
		DBHost:    GetEnv("DB_HOST", "localhost"),
		DBPort:    GetEnv("DB_PORT", "5432"),
		DBUser:    GetEnv("DB_USER", "postgres"),
		DBPass:    GetEnv("DB_PASSWORD", "password"),
		DBName:    GetEnv("DB_NAME", "lms_db"),
		JWTSecret: GetEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}
