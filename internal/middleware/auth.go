package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"myplants-server/internal/database"
	"myplants-server/internal/models"
)

var jwtSecret []byte

func init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-secret"
	}
	jwtSecret = []byte(s)
}

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

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
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
		}

		c.Next()
	}
}