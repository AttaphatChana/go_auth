package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_auth/application/ports"
)

func JWTMiddleware(provider ports.JwtProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userID, err := provider.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
