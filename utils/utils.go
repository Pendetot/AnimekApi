package utils

// APIResponse represents a standard API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
}

// CreateSuccessResponse creates a success response
func CreateSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Code:    200,
	}
}

// CreateErrorResponse creates an error response
func CreateErrorResponse(message string, code int) APIResponse {
	return APIResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}

// GetUserAgent returns a random user agent string
func GetUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
	}
	
	// For simplicity, return the first one. In a real implementation, you might want to randomize
	return userAgents[0]
}
