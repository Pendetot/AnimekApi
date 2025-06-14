package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/config"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

// GetAnimeDetail handles anime detail requests
func GetAnimeDetail(c *gin.Context) {
	slug := c.Query("slug")
	path := c.Query("path")
	
	if slug == "" && path == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Parameter slug atau path diperlukan")
		return
	}
	
	cfg := config.Load()
	var animeURL string
	
	if slug != "" {
		animeURL = cfg.BaseURL + "/" + slug
	} else {
		animeURL = cfg.BaseURL + "/" + path
	}
	
	data, err := scraper.ScrapeAnimeDetail(animeURL)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail anime: "+err.Error())
		return
	}
	
	if data == nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Detail anime tidak ditemukan")
		return
	}
	
	utils.SendSuccessResponse(c, http.StatusOK, data, "Detail anime berhasil diambil")
}
