package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/farms", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Tania",
		})
	})

	router.Run()
}
