package utils

import "time"

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Success:   false,
		Code:      400,
		Message:   message,
		Timestamp: time.Now(),
	}
}
