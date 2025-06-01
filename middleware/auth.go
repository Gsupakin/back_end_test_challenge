package middleware

import (
	"net/http"
	"strings"

	"github.com/Gsupakin/back_end_test_challeng/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenStr := strings.Split(authHeader, "Bearer ")[1]

		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// เก็บ user_id ไว้ใช้ใน handler
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
