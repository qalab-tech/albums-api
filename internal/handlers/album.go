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

func (h *AlbumHandler) GetAll(c *gin.Context) {
	albums, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		NewError(c, http.StatusInternalServerError, err.Error())
		return
	}
	NewSuccess(c, albums)
}

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
