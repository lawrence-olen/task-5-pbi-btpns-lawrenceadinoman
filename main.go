package main

import (
	"github.com/crocox/final-project/database"
	"github.com/crocox/final-project/router"

	// DB "github.com/crocox/final-project/database"
	"github.com/gin-gonic/gin"
)

func init() {
	database.ConnectionDB()
}

func main() {
	r := router.Routes()

	// DB.ConnectionDB()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
