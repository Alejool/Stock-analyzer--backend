// main.go
package main

import (
	"Backend/internal/api"
	"Backend/internal/config"
	"Backend/internal/database"
	// "Backend/internal/middleware"
	"Backend/internal/services"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Config env
	cfg := config.Load()

	// Conect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Error ejecutando migraciones:", err)
		log.Println("Error ejecutando migraciones:", err)
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
			// Wait for 40 minutes
			time.Sleep(40 * time.Minute)
		}
	}()

	//Config gin
	r := gin.Default()

	// Configurar middlewares - CORS	
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173", "https://43aa-167-0-100-7.ngrok-free.app"},
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Config routes
	api.SetupRoutes(r, stockService, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en el puerto %s", port)
	log.Fatal(r.Run(":" + port))
}
