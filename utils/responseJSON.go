package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status   bool        `json:"status"`
	Message  string      `json:"message,omitempty"`
	Duration string      `json:"duration"`
	Data     interface{} `json:"data,omitempty"`
}

func ResponseJSON(c *gin.Context, statusCode int, status bool, message string, start time.Time, data interface{}) {
	duration := time.Since(start).Milliseconds()

	c.JSON(statusCode, APIResponse{
		Duration: fmt.Sprintf("%d ms", duration),
		Status:   status,
		Message:  message,
		Data:     data,
	})
}
