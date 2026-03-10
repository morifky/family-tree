package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse creates a standard JSON envelope for success contexts
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// CreatedResponse creates a standard JSON envelope for created resources
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

// ErrorResponse creates a standard JSON envelope for contextual API failures
func ErrorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": msg,
		},
	})
}
