package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAPIInfo returns API information and available endpoints
func GetAPIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "AnoBoy API is running",
		"version": "2.0.0",
		"endpoints": []gin.H{
			{
				"path":        "/api/detail?slug=<anime_slug>",
				"description": "Get anime episode details by slug",
			},
			{
				"path":        "/api/list",
				"description": "Get complete anime list",
			},
			{
				"path":        "/api/episodes?title=<anime_title>",
				"description": "Get episodes list for an anime by title",
			},
			{
				"path":        "/api/schedule",
				"description": "Get anime broadcast schedule with images",
			},
		},
	})
}
