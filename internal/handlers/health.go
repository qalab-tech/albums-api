package handlers

import (
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	NewSuccess(c, gin.H{
		"status":  "ok",
		"service": "albums-api",
		"version": "1.0.0",
	})
}
