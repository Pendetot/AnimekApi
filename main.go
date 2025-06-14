package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/Pendetot/AnimekApi/config"
	"github.com/Pendetot/AnimekApi/handlers"
	"github.com/Pendetot/AnimekApi/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware
	r.Use(middleware.RateLimit())

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Info route
		api.GET("/", handlers.GetAPIInfo)
		
		// Anime detail route
		api.GET("/detail", handlers.GetAnimeDetail)
		
		// Anime list route
		api.GET("/list", handlers.GetAnimeList)
		
		// Episodes route
		api.GET("/episodes", handlers.GetEpisodes)
		
		// Schedule route
		api.GET("/schedule", handlers.GetSchedule)
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Route " + c.Request.URL.Path + " not found",
		})
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "1408"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
