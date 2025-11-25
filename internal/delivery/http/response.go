package http

import "github.com/gin-gonic/gin"

// База для ответа на любой запрос

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func RespondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func RespondError(c *gin.Context, code int, error string, message string) {
	c.JSON(code, ErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
	})
}

func RespondErrorWithDetails(c *gin.Context, code int, error string, message string, details interface{}) {
	c.JSON(code, ErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
		Details: details,
	})
}
