package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

// GetEpisodes handles requests for anime episodes list
func GetEpisodes(c *gin.Context) {
	// Get parameters from query
	title := c.Query("title")
	slug := c.Query("slug")
	path := c.Query("path")
	
	// Check if at least one parameter is provided
	if title == "" && slug == "" && path == "" {
		utils.SendBadRequest(c, "Parameter title, slug, atau path diperlukan")
		return
	}
	
	// Scrape episodes
	data, err := scraper.ScrapeEpisodes(title, slug, path)
	if err != nil {
		utils.LogError(err, "scraping episodes")
		utils.SendInternalError(c, fmt.Sprintf("Gagal mengambil daftar episode: %v", err))
		return
	}
	
	if data == nil || len(data.Episodes) == 0 {
		utils.SendNotFound(c, "Daftar episode tidak ditemukan")
		return
	}
	
	utils.SendSuccess(c, data, "Daftar episode berhasil diambil")
}
