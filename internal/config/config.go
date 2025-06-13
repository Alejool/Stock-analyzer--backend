

package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	APIBaseURL  string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://root@localhost:26257/truora_stocks?sslmode=disable"),
		APIKey:      getEnv("API_KEY", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MjMsImVtYWlsIjoib2xhcnRlYWxlamFuZHJvNDhAZ21haWwuY29tIiwiZXhwIjoxNzQ5NzkxMDM4LCJpZCI6IjAiLCJwYXNzd29yZCI6IicgT1IgJzEnPScxIn0.spHpSfpdsFxMYhhdBB6xJYQ6a3mXQZCaqJ2VfI6CW34"),
		APIBaseURL:  getEnv("API_BASE_URL", "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}