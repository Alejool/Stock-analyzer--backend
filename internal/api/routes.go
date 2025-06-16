package api

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, stockService *services.StockService) {
	// Debug: Print stockService details for debugging
	// log.Printf("StockService initialized with database connection: %+v\n", stockService)
	api := r.Group("/api/v1")
	{
		// Debug: Log each route registration
		log.Println("Registering route: GET /api/v1/stocks")
		api.GET("/stocks", getStocks(stockService))

		log.Println("Registering route: GET /api/v1/recommendations")
		api.GET("/recommendations", getRecommendations(stockService))

		log.Println("Registering route: GET /api/v1/health")
		api.GET("/health", healthCheck)
	}
}

func getStocks(stockService *services.StockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters models.StockFilters

		// Obtener parámetros de query
		// log.Println("filters: ", filters)

		// Parsear parámetros de query
		filters.Ticker = c.Query("ticker")
		filters.Company = c.Query("company")
		filters.Brokerage = c.Query("brokerage")
		filters.Action = c.Query("action")
		filters.Rating = c.Query("rating")
		filters.SortBy = c.Query("sort_by")
		filters.Order = c.Query("order")

		if page := c.Query("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil {
				filters.Page = p
			}
		}

		if limit := c.Query("limit"); limit != "" && limit != "-1" {
			if l, err := strconv.Atoi(limit); err == nil {
				filters.Limit = l
			}
		}

		stocks, err := stockService.GetStocks(filters)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stocks)
	}
}

func getRecommendations(stockService *services.StockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		recommendations, err := stockService.GetRecommendations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API funcionando correctamente",
	})
}
