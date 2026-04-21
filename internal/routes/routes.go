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
		// Diary routes
		api.POST("/diaries", handlers.CreateDiary)
		// Add more protected routes here
	}
}
