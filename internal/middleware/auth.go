package middleware

import (
	"net/http"
	"strings"

	"myplants-server/internal/auth"

	"github.com/gin-gonic/gin"

	"myplants-server/internal/database"
	"myplants-server/internal/models"
)

// jwtSecret now managed by internal/auth

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Parse and validate token via internal/auth
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID := uint(claims["userId"].(float64))
		c.Set("userId", userID)
		c.Set("username", claims["username"].(string))

		// Fetch user from DB to get isAdmin
		var user models.User
		if err := database.GetDB().First(&user, userID).Error; err == nil {
			c.Set("isAdmin", user.IsAdmin)
		} else {
			c.Set("isAdmin", false)
		}

		c.Next()
	}
}
