// main.go
package main

import (
	"Backend/internal/api"
	"Backend/internal/config"
	"Backend/internal/database"
	"Backend/internal/services"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "Backend/docs" // This will be generated by swag init
)

// @title Stock Analyzer API
// @version 1.0
// @description A comprehensive stock analysis and recommendation API built with Go and Gin framework.
// @description This API provides endpoints for stock data retrieval, filtering, and investment recommendations.

// @contact.name API Support
// @contact.email support@stockanalyzer.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Config env
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Error executing migrations:", err)
		log.Println("Error executing migrations:", err)
	} else {
		log.Println("Database migrations executed successfully")
	}

	// Initialize services
	stockService := services.NewStockService(db)
	apiClient := services.NewAPIClient(cfg.APIKey, cfg.APIBaseURL)

	// Initialize stock data sync
	go func() {
		for {
			for {
				log.Printf("NewAPIClient", apiClient)
				if err := stockService.SyncAllData(apiClient); err != nil {
				}
				break
			}
			time.Sleep(40 * time.Minute)
		}
	}()

	// Config gin
	r := gin.Default()

	// Configure middlewares - CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8070" },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Config routes
	api.SetupRoutes(r, stockService, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on port %s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	log.Fatal(r.Run(":" + port))
}
