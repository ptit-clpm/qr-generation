package shared

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(200, APIResponse{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(201, APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *gin.Context, status int, message string, errors any) {
	c.JSON(status, APIResponse{Success: false, Message: message, Errors: errors})
}
