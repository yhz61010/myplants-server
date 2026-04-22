package routes

import (
	"github.com/gin-gonic/gin"

	"myplants-server/internal/handlers"
	"myplants-server/internal/middleware"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Serve static assets for admin panel from templates/static
	router.Static("/admin/static", "./templates/static")

	// Admin page routes - serve HTML files
	// Admin page routes - serve HTML files. Use StaticFile for exact paths
	// to ensure serving works even if working dir differs.
	router.StaticFile("/admin", "./templates/login.html")
	router.StaticFile("/admin/login", "./templates/login.html")
	router.StaticFile("/admin/login.html", "./templates/login.html")
	router.StaticFile("/admin/", "./templates/index.html")
	router.StaticFile("/admin/users", "./templates/users.html")
	// Keep param routes as handlers to allow dynamic IDs but serve static HTML
	router.GET("/admin/users/:id", func(c *gin.Context) {
		c.File("./templates/user_detail.html")
	})
	router.StaticFile("/admin/diaries", "./templates/diaries.html")
	router.GET("/admin/diaries/:id", func(c *gin.Context) {
		c.File("./templates/diary_detail.html")
	})

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
		api.POST("/diaries", handlers.CreateDiary)
		api.GET("/timeline", handlers.ListContents)
		api.GET("/diaries", handlers.ListDiaries)
		api.GET("/diaries/:id", handlers.GetDiary)
		api.PUT("/diaries/:id", handlers.UpdateDiary)
		api.DELETE("/diaries/:id", handlers.DeleteDiary)

		// Upload endpoint for images (stores URL only)
		api.POST("/upload", handlers.UploadImage)

		// User management
		api.GET("/users", handlers.ListUsers)
		api.GET("/users/:id", handlers.GetUser)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
	}

	// Admin routes (require admin privileges)
	admin := api.Group("/admin")
	admin.Use(handlers.AdminMiddleware())
	{
		// Admin user management
		admin.GET("/users", handlers.AdminListUsers)
		admin.GET("/users/:id", handlers.AdminGetUser)
		admin.PUT("/users/:id", handlers.AdminUpdateUser)
		admin.DELETE("/users/:id", handlers.AdminDeleteUser)

		// Admin diary management
		admin.POST("/diaries", handlers.AdminCreateDiary)
		admin.PUT("/diaries/:id", handlers.AdminUpdateDiary)
		admin.DELETE("/diaries/:id", handlers.AdminDeleteDiary)
	}
}
