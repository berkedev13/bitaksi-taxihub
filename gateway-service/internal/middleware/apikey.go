package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware() gin.HandlerFunc {
	requiredKey := os.Getenv("API_KEY")
	if requiredKey == "" {
		panic("API_KEY environment variable is not set")
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" || apiKey != requiredKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing API key",
			})
			return
		}

		c.Next()
	}
}
