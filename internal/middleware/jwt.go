package middleware

import (
	"time"
	"strings"
	"net/http"
	"fmt"

	"Backend/internal/config"
	"github.com/gin-gonic/gin"
	"Backend/internal/entity"
	"github.com/golang-jwt/jwt/v5"


)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {


		// Get token from header
		// token, err := c.Cookie("token")
				token := c.GetHeader("Authorization");

		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Unauthorized",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// }

		// Check if header has the correct format
		parts := strings.Split(token, ".")
		if len(parts) != 3 || 
				parts[0] != "Bearer" ||
				len(parts[1]) == 0 || 
				len(parts[2]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token does not exist",
			})
			c.Abort()
			return
		}

		// Parse and validate  token
		
		tokenJwt, err := jwt.ParseWithClaims(token, &entity.Claimes{}, func(token *jwt.Token) (interface{}, error) {

			// validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKeyType
			}
			return cfg.JwtSecretKey, nil

		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
			
			// check if the token is valid

			if claims, ok := tokenJwt.Claims.(*entity.Claimes); ok && tokenJwt.Valid {
				// if token is valid, create a new userJwt
				userJwt := &entity.UserJwt{
					UserId: claims.UserId,
					Username: claims.Username,
				}
				c.Set("user", userJwt)
				c.Next()

			}else{
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token",
				})
			}
			// return entity.JwtSecretKey, nil
	

	
	}
}

func GenerateToken(user *entity.UserJwt,cfg *config.Config) (string, error) {	
	expirationtime := time.Now().Add(time.Hour * 24 * 7)

	// Create claims with user data
	claims := &entity.Claimes{
		UserId: user.UserId,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationtime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// For debugging purposes, print the token details
	fmt.Printf("Generated Token: %+v\n", token)
	fmt.Printf("JwtSecretKey: %+v\n", cfg.JwtSecretKey)

		return token.SignedString(cfg.JwtSecretKey) 
	
	
	// return token.SignedString("1234141") 
}

