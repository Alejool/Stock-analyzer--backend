
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	APIKey       string
	APIBaseURL   string
	Port         string
	DatabaseName string
	Environment  string
	GinMode      string
}

// 
func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Printf("No se encontró archivo .env, usando variables de entorno del sistema: %v", err)
	}

	config := &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		APIKey:       getEnv("API_KEY", ""),
		APIBaseURL:   getEnv("API_BASE_URL","" ),
		Port:         getEnv("PORT", "8080"),
		DatabaseName: getEnv("DATABASE_NAME", "stock_tracking"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		GinMode:      getEnv("GIN_MODE", "debug"),
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Error en configuración: %v", err)
	}

	return config
}

// LoadFromFile carga configuración desde un archivo específico
func LoadFromFile(filename string) *Config {
	if err := godotenv.Load(filename); err != nil {
		log.Fatalf("Error cargando archivo de configuración %s: %v", filename, err)
	}
	return Load()
}

// LoadFromEnv carga solo desde variables de entorno del sistema (sin .env)
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

	// Aplicar defaults si están vacías
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
		log.Fatalf("Error en configuración: %v", err)
	}

	return config
}

// Validate verifica que la configuración sea válida
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL es requerida")
	}
	if c.APIBaseURL == "" {
		return fmt.Errorf("API_BASE_URL es requerida")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY es requerida")
	}
	if c.APIBaseURL == "" {
		return fmt.Errorf("API_BASE_URL es requerida")
	}
	
	// Validar puerto
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT debe ser un número válido: %v", err)
	}

	return nil
}

// IsDevelopment verifica si está en modo desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction verifica si está en modo producción
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// IsTest verifica si está en modo test
func (c *Config) IsTest() bool {
	return c.Environment == "test" || c.Environment == "testing"
}

// GetDatabaseConfig retorna configuración específica de base de datos
func (c *Config) GetDatabaseConfig() map[string]string {
	return map[string]string{
		"url":  c.DatabaseURL,
		"name": c.DatabaseName,
	}
}

// Print imprime la configuración (sin mostrar datos sensibles)
func (c *Config) Print() {
	log.Println("=== Configuración de la aplicación ===")
	log.Printf("Environment: %s", c.Environment)
	log.Printf("Port: %s", c.Port)
	log.Printf("Database Name: %s", c.DatabaseName)
	log.Printf("API Base URL: %s", c.APIBaseURL)
	log.Printf("Gin Mode: %s", c.GinMode)
	log.Printf("Database URL: %s", maskSensitiveData(c.DatabaseURL))
	log.Printf("API Key: %s", maskSensitiveData(c.APIKey))
	log.Println("========================================")
}

// getEnv obtiene variable de entorno 
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskSensitiveData(data string) string {
	if len(data) <= 8 {
		return "***"
	}
	return data[:4] + "..." + data[len(data)-4:]
}