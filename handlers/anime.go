package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/config"
	"github.com/Pendetot/AnimekApi/scraper"
	"github.com/Pendetot/AnimekApi/utils"
)

var cfg = config.Load()

// GetAnimeDetail handles anime detail requests
func GetAnimeDetail(c *gin.Context) {
	slug := c.Query("slug")
	path := c.Query("path")
	
	var animeURL string
	if slug != "" {
		animeURL = fmt.Sprintf("%s/%s", cfg.BaseURL, slug)
	} else if path != "" {
		animeURL = fmt.Sprintf("%s/%s", cfg.BaseURL, path)
	} else {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Either slug or path parameter is required"))
		return
	}
	
	data, err := scraper.ScrapeAnimeDetail(animeURL, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse(fmt.Sprintf("Failed to retrieve anime details: %v", err)))
		return
	}
	
	if data == nil {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Anime details not found"))
		return
	}
	
	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Anime details retrieved successfully"))
}

// GetAnimeList handles anime list requests
func GetAnimeList(c *gin.Context) {
	data, err := scraper.ScrapeAnimeList(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse(fmt.Sprintf("Failed to retrieve anime list: %v", err)))
		return
	}
	
	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Anime list retrieved successfully"))
}

// GetEpisodes handles episodes list requests
func GetEpisodes(c *gin.Context) {
	title := c.Query("title")
	slug := c.Query("slug")
	path := c.Query("path")
	
	var animeURL string
	if title != "" {
		animeURL = fmt.Sprintf("%s/anime/%s", cfg.BaseURL, title)
	} else if slug != "" {
		animeURL = fmt.Sprintf("%s/%s", cfg.BaseURL, slug)
	} else if path != "" {
		animeURL = fmt.Sprintf("%s/%s", cfg.BaseURL, path)
	} else {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Either title, slug, or path parameter is required"))
		return
	}
	
	data, err := scraper.ScrapeEpisodes(animeURL, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse(fmt.Sprintf("Failed to retrieve episodes: %v", err)))
		return
	}
	
	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Episodes retrieved successfully"))
}

// GetSchedule handles schedule requests
func GetSchedule(c *gin.Context) {
	data, err := scraper.ScrapeSchedule(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse(fmt.Sprintf("Failed to retrieve schedule: %v", err)))
		return
	}
	
	c.JSON(http.StatusOK, utils.CreateSuccessResponse(data, "Schedule retrieved successfully"))
}
