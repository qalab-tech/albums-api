package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"albums-api/internal/config"
	"albums-api/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getValidToken получает свежий токен от auth-service
func getValidToken(t *testing.T) string {
	payload := map[string]string{
		"username": "test",
		"password": "test",
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://auth-service:5001/auth/login", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to connect to auth-service")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Auth service should return 200")

	var result struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to parse token response")
	require.NotEmpty(t, result.Token, "Token should not be empty")

	return result.Token
}

func TestIntegration_PublicEndpoints(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv := server.New(cfg)
	router := srv.GetRouter()

	t.Run("GET /api/v1/albums - public access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/albums", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /api/v1/health - public access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestIntegration_ProtectedEndpoints(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv := server.New(cfg)
	router := srv.GetRouter()
	token := getValidToken(t)

	t.Run("POST /api/v1/albums - with valid token", func(t *testing.T) {
		payload := `{"title":"Integration Test Album","artist":"Test Artist","price":99.99}`
		req := httptest.NewRequest("POST", "/api/v1/albums", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("POST /api/v1/albums - without token should fail", func(t *testing.T) {
		payload := `{"title":"No Token Test","artist":"Test","price":10}`
		req := httptest.NewRequest("POST", "/api/v1/albums", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}