package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"myplants-server/internal/database"
	"myplants-server/internal/routes"
)

func main() {
	// Set release mode for production
	// gin.SetMode(gin.ReleaseMode)

	// Initialize database
	database.InitDB()

	// Create Gin router
	router := gin.New() // Use New() instead of Default() to avoid default middleware in release mode

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
