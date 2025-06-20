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

	// Debug: Print stockService details for debugging
	// log.Printf("StockService initialized with database connection: %+v\n", stockService)

	r.GET("/health", healthCheck)
	r.POST("/get-token", getToken(cfg))

	api := r.Group("/api/v1")
	// api.Use(middleware.AuthMiddleware(cfg))
	{
		api.GET("/stocks", getStocks(stockService))
		api.GET("/recommendations", getRecommendations(stockService))

	}
	// r.Use(middleware.AuthMiddleware())
}

func getToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginRequest entity.LoginRequest

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Faltan parámetros de entrada",
			})
			return
		}

		if(loginRequest.Username != "dashboard" && loginRequest.Username != "admin"){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Usuario no autorizado",
			})
			return
		}

		if loginRequest.Username == "dashboard" {
			user := &entity.UserJwt{
				UserId:   1,
				Username: loginRequest.Username,
			}

			
			fmt.Printf("user: ", user)
			fmt.Printf("cfg: ", cfg)
			token, err := middleware.GenerateToken(user, cfg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error al generar token - "+err.Error(),
				})
				return
			}



			// user, err := services.Login(username, password)
			// if err != nil {
			// 	c.JSON(http.StatusUnauthorized, gin.H{
			// 		"error": "Error al iniciar sesión",
			// 	})
			// 	return
			// }


			

			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}
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
