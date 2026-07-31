package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
