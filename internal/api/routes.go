
package api

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"../models"
	"Backend/internal/services"
)

func SetupRoutes(r *gin.Engine, stockService *services.StockService) {
	api := r.Group("/api/v1")
	{
		api.GET("/stocks", getStocks(stockService))
		api.GET("/recommendations", getRecommendations(stockService))
		api.GET("/health", healthCheck)
	}
}

func getStocks(stockService *services.StockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters models.StockFilters
		
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
		
		if limit := c.Query("limit"); limit != "" {
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
		"status": "ok",
		"message": "API funcionando correctamente",
	})
}