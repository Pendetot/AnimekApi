package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse(data, message))
}

// SendError sends an error response
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse(message))
}

// SendNotFound sends a 404 error response
func SendNotFound(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, message)
}

// SendBadRequest sends a 400 error response
func SendBadRequest(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, message)
}

// SendInternalError sends a 500 error response
func SendInternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, message)
}

// LogError logs an error with timestamp
func LogError(err error, context string) {
	fmt.Printf("[ERROR] %s - %s: %v
", time.Now().Format(time.RFC3339), context, err)
}

// LogInfo logs an info message with timestamp
func LogInfo(message string) {
	fmt.Printf("[INFO] %s - %s
", time.Now().Format(time.RFC3339), message)
}
