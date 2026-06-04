package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"albums-api/internal/handlers" // для ErrorResponse

	"github.com/gin-gonic/gin"
)

type AuthServiceResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
	Error  string `json:"error,omitempty"`
}

func JWTAuth(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Error: "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Вызываем auth-service /auth/validate
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, "GET", authServiceURL+"/auth/validate", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, handlers.ErrorResponse{Error: "Auth service unavailable"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errResp AuthServiceResponse
			json.NewDecoder(resp.Body).Decode(&errResp)
			c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Error: errResp.Error})
			c.Abort()
			return
		}

		// Если токен валиден — продолжаем
		c.Next()
	}
}
