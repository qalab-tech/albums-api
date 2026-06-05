package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"albums-api/internal/client"
	"albums-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func JWTAuth(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Error: "Authorization header required. Use Bearer token",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверяем токен через Auth Service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := authClient.ValidateToken(ctx, token); err != nil {
			c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Error: err.Error(),
			})
			c.Abort()
			return
		}

		// Токен валиден — продолжаем
		c.Next()
	}
}
