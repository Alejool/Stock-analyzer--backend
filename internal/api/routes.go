package api

import (
	"Backend/internal/config"
	"Backend/internal/entity"
	"Backend/internal/middleware"
	"Backend/internal/models"
	"Backend/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, stockService *services.StockService, cfg *config.Config) {
	r.GET("/health", healthCheck)
	r.POST("/get-token", getToken(cfg))

	api := r.Group("/api/v1")
	{
		api.GET("/stocks", getStocks(stockService))
		api.GET("/recommendations", getRecommendations(stockService))
	}
}

// @Summary Get authentication token
// @Description Authenticate user and receive JWT token for API access
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body entity.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /get-token [post]
func getToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest entity.LoginRequest

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing input parameters",
			})
			return
		}

		if loginRequest.Username != "dashboard" && loginRequest.Username != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unauthorized user",
			})
			return
		}

		if loginRequest.Username == "dashboard" || loginRequest.Username == "admin" {
			user := &entity.UserJwt{
				UserId:   1,
				Username: loginRequest.Username,
			}

			fmt.Printf("user: %+v\n", user)
			fmt.Printf("cfg: %+v\n", cfg)

			token, err := middleware.GenerateToken(user, cfg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error generating token - " + err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}
	}
}

// @Summary      Get stocks with filtering
// @Description  Retrieve a list of stocks with optional filtering and pagination
// @Tags         Stocks
// @Accept       json
// @Produce      json
// @Param        ticker     query  string  false  "Stock ticker symbol"
// @Param        company    query  string  false  "Company name"
// @Param        brokerage  query  string  false  "Brokerage firm"
// @Param        action     query  string  false  "Recommended action (buy, sell, hold)"
// @Param        rating     query  string  false  "Stock rating"
// @Param        sort_by    query  string  false  "Sort field"
// @Param        order      query  string  false  "Sort order (asc, desc)"
// @Param        page       query  int     false  "Page number for pagination"
// @Param        limit      query  int     false  "Number of items per page"
// @Param        today      query  string  false  "Filter for today's data"
// @Success      200        {object}  models.StockResponse  "List of stocks with metadata"

// @Security     BearerAuth
// @Router       /api/v1/stocks [get]
func getStocks(stockService *services.StockService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters models.StockFilters

		// Parse query parameters
		filters.Ticker = c.Query("ticker")
		filters.Company = c.Query("company")
		filters.Action = c.Query("action")
		filters.Rating = c.Query("rating")
		filters.SortBy = c.Query("sort_by")
		filters.Order = c.Query("order")
		filters.Today = c.Query("today")

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

// @Summary Get stock recommendations
// @Description Retrieve a list of stock recommendations
// @Tags Recommendations
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]models.Recommendation
// @Failure 500 {object} map[string]string "error"
// @Security BearerAuth
// @Router /api/v1/recommendations [get]
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

// @Summary Health check
// @Description Check if the API is running and healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API working correctly",
	})
}
