// Package config provides functionality for loading and managing application configuration
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application
type Config struct {
	DatabaseURL  string // URL for database connection
	APIKey       string // API key for external services
	APIBaseURL   string // Base URL for API endpoints
	Port         string // Server port number
	DatabaseName string // Name of the database
	Environment  string // Current environment (development/production/test)
	GinMode      string // Gin framework mode
	JwtSecretKey []byte // Secret key for JWT token generation
}

// Load reads configuration from .env file or environment variables
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables: %v", err)
	}

	config := &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		APIKey:       getEnv("API_KEY", ""),
		APIBaseURL:   getEnv("API_BASE_URL", ""),
		Port:         getEnv("PORT", "8080"),
		DatabaseName: getEnv("DATABASE_NAME", "stock_tracking"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		GinMode:      getEnv("GIN_MODE", "debug"),
		JwtSecretKey: []byte(getEnv("JWT_SECRET_KEY", "")),
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	return config
}

// LoadFromFile loads configuration from a specific file
func LoadFromFile(filename string) *Config {
	if err := godotenv.Load(filename); err != nil {
		log.Fatalf("Error loading configuration file %s: %v", filename, err)
	}
	return Load()
}

// LoadFromEnv loads configuration only from system environment variables (no .env file)
func LoadFromEnv() *Config {
	config := &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		APIKey:       os.Getenv("API_KEY"),
		APIBaseURL:   os.Getenv("API_BASE_URL"),
		Port:         os.Getenv("PORT"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		Environment:  os.Getenv("ENVIRONMENT"),
		GinMode:      os.Getenv("GIN_MODE"),
	}

	// Apply defaults if empty
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.DatabaseName == "" {
		config.DatabaseName = "stock_tracking"
	}
	if config.Environment == "" {
		config.Environment = "production"
	}
	if config.GinMode == "" {
		config.GinMode = "release"
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	return config
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.APIBaseURL == "" {
		return fmt.Errorf("API_BASE_URL is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY is required")
	}
	if c.APIBaseURL == "" {
		return fmt.Errorf("API_BASE_URL is required")
	}

	// Validate port
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be a valid number: %v", err)
	}

	return nil
}

// IsDevelopment checks if the environment is set to development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction checks if the environment is set to production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// IsTest checks if the environment is set to test mode
func (c *Config) IsTest() bool {
	return c.Environment == "test" || c.Environment == "testing"
}

// GetDatabaseConfig returns specific database configuration
func (c *Config) GetDatabaseConfig() map[string]string {
	return map[string]string{
		"url":  c.DatabaseURL,
		"name": c.DatabaseName,
	}
}

// getEnv retrieves environment variable with a default fallback value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// // Print displays the configuration (hiding sensitive data)
// func (c *Config) Print() {
// 	log.Println("=== Application Configuration ===")
// 	log.Printf("Environment: %s", c.Environment)
// 	log.Printf("Port: %s", c.Port)
// 	log.Printf("Database Name: %s", c.DatabaseName)
// 	log.Printf("API Base URL: %s", c.APIBaseURL)
// 	log.Printf("Gin Mode: %s", c.GinMode)
// 	log.Printf("Database URL: %s", maskSensitiveData(c.DatabaseURL))
// 	log.Printf("API Key: %s", maskSensitiveData(c.APIKey))
// 	log.Println("========================================")
// }

// // maskSensitiveData masks sensitive information for display purposes
// func maskSensitiveData(data string) string {
// 	if len(data) <= 8 {
// 		return "***"
// 	}
// 	return data[:4] + "..." + data[len(data)-4:]
// }
