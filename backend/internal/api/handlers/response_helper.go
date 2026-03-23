package handlers

import "github.com/gin-gonic/gin"

func PaginatedResponse(c *gin.Context, key string, data interface{}, total int64, page, limit int) {
	c.JSON(200, gin.H{
		key:          data,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}
