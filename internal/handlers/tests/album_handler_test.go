package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"albums-api/internal/handlers"
	"albums-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAll(ctx context.Context) ([]models.Album, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Album), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id int) (models.Album, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, album models.NewAlbum) (models.Album, error) {
	args := m.Called(ctx, album)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id int, update models.UpdateAlbum) (models.Album, error) {
	args := m.Called(ctx, id, update)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *handlers.AlbumHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		albums := api.Group("/albums")
		{
			albums.GET("", handler.GetAll)
			albums.GET("/:id", handler.GetByID)
			albums.POST("", handler.Create)
			albums.PUT("/:id", handler.Update)
			albums.DELETE("/:id", handler.Delete)
		}
	}
	return r
}

// ==================== ТЕСТЫ ====================

func TestGetAllAlbums(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	albums := []models.Album{
		{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: 2, Title: "Kind of Blue", Artist: "Miles Davis", Price: 39.99},
	}
	mockRepo.On("GetAll", mock.Anything).Return(albums, nil)

	router := setupTestRouter(h)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/albums", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Новый формат ответа
	var resp handlers.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)
}

func TestGetAlbumByID(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	album := models.Album{ID: 5, Title: "Test Album", Artist: "Test Artist", Price: 29.99}
	mockRepo.On("GetByID", mock.Anything, 5).Return(album, nil)

	router := setupTestRouter(h)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/albums/5", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp handlers.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
}

func TestGetAlbumByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	mockRepo.On("GetByID", mock.Anything, 999).Return(models.Album{}, assert.AnError)

	router := setupTestRouter(h)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/albums/999", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp handlers.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
}

func TestCreateAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	newAlbum := models.NewAlbum{Title: "New Album", Artist: "New Artist", Price: 49.99}
	created := models.Album{ID: 10, Title: "New Album", Artist: "New Artist", Price: 49.99}

	mockRepo.On("Create", mock.Anything, newAlbum).Return(created, nil)

	router := setupTestRouter(h)
	body := `{"title":"New Album","artist":"New Artist","price":49.99}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/albums", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	mockRepo.On("Delete", mock.Anything, 3).Return(nil)

	router := setupTestRouter(h)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/albums/3", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	h := handlers.NewAlbumHandler(mockRepo)

	update := models.UpdateAlbum{Title: strPtr("Updated Title")}
	updated := models.Album{ID: 1, Title: "Updated Title", Artist: "John Coltrane", Price: 56.99}

	mockRepo.On("Update", mock.Anything, 1, update).Return(updated, nil)

	router := setupTestRouter(h)
	body := `{"title":"Updated Title"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/albums/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func strPtr(s string) *string {
	return &s
}
