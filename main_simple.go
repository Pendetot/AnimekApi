package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "AnoBoy API v2.0 - Golang Version",
			"status":  "success",
		})
	})
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	fmt.Println("Server starting on port 1408...")
	log.Fatal(http.ListenAndServe(":1408", r))
}
