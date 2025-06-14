package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

// GetSchedule handles requests for anime broadcast schedule
func GetSchedule(c *gin.Context) {
	// Scrape schedule
	data, err := scraper.ScrapeSchedule()
	if err != nil {
		utils.LogError(err, "scraping schedule")
		utils.SendInternalError(c, fmt.Sprintf("Gagal mengambil jadwal tayang: %v", err))
		return
	}
	
	if data == nil || len(data.Items) == 0 {
		utils.SendNotFound(c, "Jadwal tayang tidak ditemukan")
		return
	}
	
	utils.SendSuccess(c, data, "Jadwal tayang berhasil diambil")
}
