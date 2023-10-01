package router

import (
	"github.com/crocox/final-project/controllers"
	"github.com/crocox/final-project/middleware"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	router := gin.Default()

	// Group for API
	api := router.Group("/api")

	// Users
	api.POST("/auth/register", controllers.Register)
	api.POST("/auth/login", controllers.Login)
	api.Static("/photos", "./uploads")
	// api.GET("/validate", controllers.Validate)

	api.Use(middleware.Auth)
	// Photos
	api.GET("/photos", controllers.GetPhoto)
	api.POST("/photos", controllers.UploadPhoto)
	api.PUT("/photos/:id", controllers.UpdatePhoto)
	api.DELETE("/photos/:id", controllers.DeletePhoto)

	return router
}
