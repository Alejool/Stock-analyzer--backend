// main.go
package main

import (
	"log"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/alejool/RetosTecnicos/Backend/internal/api"
	"Backend/internal/config"
	"Backend/internal/database"
	"Backend/internal/services"

)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Conectar a la base de datos
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close()

	// Ejecutar migraciones
	if err := database.Migrate(db); err != nil {
		log.Fatal("Error ejecutando migraciones:", err)
	}

	// Inicializar servicios
	stockService := services.NewStockService(db)
	apiClient := services.NewAPIClient(cfg.APIKey)

	// Sincronizar datos iniciales
	go func() {
		if err := stockService.SyncAllData(apiClient); err != nil {
			log.Printf("Error sincronizando datos: %v", err)
		}
	}()

	// Configurar router
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Configurar rutas
	api.SetupRoutes(r, stockService)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Servidor iniciado en el puerto %s", port)
	log.Fatal(r.Run(":" + port))
}

