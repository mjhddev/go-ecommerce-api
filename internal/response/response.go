package response

import (
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Message: message,
	})
}
