package utils

import "github.com/gin-gonic/gin"

// APIResponse represents a standard API response structure
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateSuccessResponse creates a success response
func CreateSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// CreateErrorResponse creates an error response
func CreateErrorResponse(message string) APIResponse {
	return APIResponse{
		Status:  "error",
		Message: message,
	}
}

// SendSuccessResponse sends a success response with data
func SendSuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, CreateSuccessResponse(data, message))
}

// SendErrorResponse sends an error response
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, CreateErrorResponse(message))
}
