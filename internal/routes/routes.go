package routes

import (
	"github.com/gin-gonic/gin"

	"myplants-server/internal/handlers"
	"myplants-server/internal/middleware"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Content CRUD routes (protected)
		api.POST("/contents", handlers.CreateContent)
		api.GET("/contents", handlers.ListContents)
		api.GET("/contents/:id", handlers.GetContent)
		api.PUT("/contents/:id", handlers.UpdateContent)
		api.DELETE("/contents/:id", handlers.DeleteContent)

		// Alias routes for backwards compatibility
		api.POST("/diaries", handlers.CreateContent)
		api.GET("/timeline", handlers.ListContents)

		// Upload endpoint for images (stores URL only)
		api.POST("/upload", handlers.UploadImage)
	}
}
