package utils

import (
	"net/http"
	"time"

	"github.com/Pendetot/AnimekApi/config"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateSuccessResponse creates a success response
func CreateSuccessResponse(data interface{}, message string) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// CreateErrorResponse creates an error response
func CreateErrorResponse(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}

// CreateHTTPClient creates an HTTP client with timeout and user agent
func CreateHTTPClient(cfg *config.Config) *http.Client {
	return &http.Client{
		Timeout: cfg.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
		},
	}
}

// SetUserAgent sets the user agent for HTTP requests
func SetUserAgent(req *http.Request, userAgent string) {
	req.Header.Set("User-Agent", userAgent)
}
