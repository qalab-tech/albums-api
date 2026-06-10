package tests

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

// getValidToken получает свежий токен
func getValidToken(t *testing.T) string {
	payload := map[string]string{"username": "test", "password": "test"}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://auth-service:5001/auth/login", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	require.NotEmpty(t, result.Token)
	return result.Token
}

// getLastAlbumID получает список и возвращает ID последнего альбома
func getLastAlbumID(t *testing.T) int {
	resp, err := http.Get("http://localhost:8080/api/v1/albums")
	require.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Data    []struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	require.True(t, result.Success)
	require.NotEmpty(t, result.Data)

	return result.Data[len(result.Data)-1].ID
}

func TestIntegration_ProtectedCRUD(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv := server.New(cfg)
	router := srv.GetRouter()
	token := getValidToken(t)

	t.Run("Create Album", func(t *testing.T) {
		payload := `{"title":"Dynamic Test Album","artist":"Test Artist","price":77.77}`
		req := httptest.NewRequest("POST", "/api/v1/albums", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update Last Album", func(t *testing.T) {
		id := getLastAlbumID(t)
		payload := `{"title":"Updated Dynamic Title","price":88.88}`

		req := httptest.NewRequest("PUT", "/api/v1/albums/"+string(rune(id)), bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Last Album", func(t *testing.T) {
		id := getLastAlbumID(t)

		req := httptest.NewRequest("DELETE", "/api/v1/albums/"+string(rune(id)), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
