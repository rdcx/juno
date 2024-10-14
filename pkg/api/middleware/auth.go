package middleware

import (
	"context"
	"juno/pkg/api/auth"
	"juno/pkg/api/auth/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("SECRET") == "" {
			panic("SECRET environment variable required")
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort() // Prevent further handlers from running
			return
		}

		u, err := service.TokenToUser(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // Prevent further handlers from running
			return
		}

		c.Request = c.Request.WithContext(auth.WithUser(context.Background(), u))
		c.Next()
	}
}
