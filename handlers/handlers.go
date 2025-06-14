package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

// GetAPIInfo returns API information and available endpoints
func GetAPIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "AnoBoy API is running",
		"version": "2.0.0",
		"endpoints": []gin.H{
			{"path": "/api/detail?slug=<anime_slug>", "description": "Get anime episode details by slug"},
			{"path": "/api/list", "description": "Get complete anime list"},
			{"path": "/api/episodes?title=<anime_title>", "description": "Get episodes list for an anime by title"},
			{"path": "/api/schedule", "description": "Get anime broadcast schedule with images"},
		},
	})
}

// GetAnimeDetail gets anime episode details by slug or path
func GetAnimeDetail(c *gin.Context) {
	slug := c.Query("slug")
	path := c.Query("path")

	if slug == "" && path == "" {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Either slug or path parameter is required", http.StatusBadRequest))
		return
	}

	var animeURL string
	if slug != "" {
		animeURL = "https://ww1.anoboy.app/" + slug
	} else {
		animeURL = "https://ww1.anoboy.app/" + path
	}

	data, err := scraper.ScrapeAnimeDetail(animeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Failed to retrieve anime details: "+err.Error(), http.StatusInternalServerError))
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Anime details not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Anime details retrieved successfully"))
}

// GetAnimeList gets complete anime list
func GetAnimeList(c *gin.Context) {
	data, err := scraper.ScrapeAnimeList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Failed to retrieve anime list: "+err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Anime list retrieved successfully"))
}

// GetEpisodes gets episode list for an anime
func GetEpisodes(c *gin.Context) {
	title := c.Query("title")
	slug := c.Query("slug")
	path := c.Query("path")

	if title == "" && slug == "" && path == "" {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Either title, slug, or path parameter is required", http.StatusBadRequest))
		return
	}

	var animeURL string
	if title != "" {
		animeURL = "https://ww1.anoboy.app/anime/" + title
	} else if slug != "" {
		animeURL = "https://ww1.anoboy.app/" + slug
	} else {
		animeURL = "https://ww1.anoboy.app/" + path
	}

	data, err := scraper.ScrapeEpisodes(animeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Failed to retrieve episodes: "+err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Episodes retrieved successfully"))
}

// GetSchedule gets anime broadcast schedule with images
func GetSchedule(c *gin.Context) {
	data, err := scraper.ScrapeSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Failed to retrieve schedule: "+err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Schedule retrieved successfully"))
}
