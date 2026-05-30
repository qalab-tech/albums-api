package handlers

import (
	"net/http"
	"strconv"

	"albums-api/internal/models"
	"albums-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	repo repository.AlbumRepository
}

func NewAlbumHandler(repo repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{repo: repo}
}

// GetAll godoc
// @Summary Получить все альбомы
// @Description Возвращает список всех альбомов из базы данных
// @Tags albums
// @Accept json
// @Produce json
// @Success 200 {object} handlers.SuccessResponse{data=[]models.Album}
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/v1/albums [get]
func (h *AlbumHandler) GetAll(c *gin.Context) {
	albums, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		NewError(c, http.StatusInternalServerError, err.Error())
		return
	}
	NewSuccess(c, albums)
}

// GetByID godoc
// @Summary Получить альбом по ID
// @Description Возвращает один альбом по его ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "ID альбома"
// @Success 200 {object} handlers.SuccessResponse{data=models.Album}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /api/v1/albums/{id} [get]
func (h *AlbumHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewError(c, http.StatusBadRequest, "неверный формат id")
		return
	}

	album, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		NewError(c, http.StatusNotFound, err.Error())
		return
	}
	NewSuccess(c, album)
}

// Create godoc
// @Summary Создать новый альбом
// @Description Добавляет новый альбом в базу
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.NewAlbum true "Данные альбома"
// @Success 201 {object} handlers.SuccessResponse{data=models.Album}
// @Failure 400 {object} handlers.ErrorResponse
// @Router /api/v1/albums [post]
func (h *AlbumHandler) Create(c *gin.Context) {
	var newAlbum models.NewAlbum
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	album, err := h.repo.Create(c.Request.Context(), newAlbum)
	if err != nil {
		NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccess(c, album, http.StatusCreated)
}

// Update godoc
// @Summary Обновить альбом
// @Description Обновляет информацию об альбоме (частично)
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "ID альбома"
// @Param album body models.UpdateAlbum true "Данные для обновления"
// @Success 200 {object} handlers.SuccessResponse{data=models.Album}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /api/v1/albums/{id} [put]
func (h *AlbumHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewError(c, http.StatusBadRequest, "неверный формат id")
		return
	}

	var update models.UpdateAlbum
	if err := c.ShouldBindJSON(&update); err != nil {
		NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	album, err := h.repo.Update(c.Request.Context(), id, update)
	if err != nil {
		NewError(c, http.StatusNotFound, err.Error())
		return
	}

	NewSuccess(c, album)
}

// Delete godoc
// @Summary Удалить альбом
// @Description Удаляет альбом по ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "ID альбома"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /api/v1/albums/{id} [delete]
func (h *AlbumHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewError(c, http.StatusBadRequest, "неверный формат id")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		NewError(c, http.StatusNotFound, err.Error())
		return
	}

	NewSuccess(c, gin.H{
		"message": "Альбом успешно удалён",
		"id":      id,
	})
}
