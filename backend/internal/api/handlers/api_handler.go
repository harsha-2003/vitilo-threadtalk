package handlers

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}

func SuccessResponse(c *gin.Context, status int, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["success"] = true
	c.JSON(status, data)
}
