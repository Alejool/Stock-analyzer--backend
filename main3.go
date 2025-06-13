package main

import "github.com/gin-gonic/gin"


func main11() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello mundo!",
		})
	})

	r.Run(":8080")
}
