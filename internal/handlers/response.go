package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse — стандартный успешный ответ
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse — стандартный ответ с ошибкой
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// NewSuccess — обёртка для успешных ответов
func NewSuccess(c *gin.Context, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	c.JSON(code, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// NewError — обёртка для ошибок
func NewError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
	})
}
