package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

// GetAnimeList handles anime list requests
func GetAnimeList(c *gin.Context) {
	data, err := scraper.ScrapeAnimeList()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar anime: "+err.Error())
		return
	}
	
	if len(data) == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "Daftar anime tidak ditemukan")
		return
	}
	
	utils.SendSuccessResponse(c, http.StatusOK, data, "Daftar anime berhasil diambil")
}
