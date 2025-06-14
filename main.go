package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := APIResponse{
			Status:  "ok",
			Message: "Server is running",
			Data: map[string]string{
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// API info endpoint
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := APIResponse{
			Status:  "success",
			Message: "AnoBoy API is running",
			Data: map[string]interface{}{
				"version": "2.0.0",
				"endpoints": []map[string]string{
					{
						"path":        "/api/detail?slug=<anime_slug>",
						"description": "Dapatkan detail episode anime berdasarkan slug",
					},
					{
						"path":        "/api/list",
						"description": "Dapatkan daftar lengkap anime",
					},
					{
						"path":        "/api/episodes?title=<anime_title>",
						"description": "Dapatkan daftar episode untuk anime berdasarkan judul",
					},
					{
						"path":        "/api/schedule",
						"description": "Dapatkan jadwal tayang anime dengan gambar",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// API detail endpoint (placeholder)
	http.HandleFunc("/api/detail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		slug := r.URL.Query().Get("slug")
		if slug == "" {
			response := APIResponse{
				Status:  "error",
				Message: "Parameter slug diperlukan",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		response := APIResponse{
			Status:  "success",
			Message: "Detail anime berhasil diambil",
			Data: map[string]string{
				"title": "Sample Anime",
				"slug":  slug,
				"note":  "This is a placeholder response - scraping functionality to be implemented",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// API list endpoint (placeholder)
	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := APIResponse{
			Status:  "success",
			Message: "Daftar anime berhasil diambil",
			Data: []map[string]string{
				{
					"title": "Sample Anime 1",
					"url":   "/anime/sample-1",
				},
				{
					"title": "Sample Anime 2", 
					"url":   "/anime/sample-2",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// API episodes endpoint (placeholder)
	http.HandleFunc("/api/episodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		title := r.URL.Query().Get("title")
		if title == "" {
			response := APIResponse{
				Status:  "error",
				Message: "Parameter title diperlukan",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		response := APIResponse{
			Status:  "success",
			Message: "Daftar episode berhasil diambil",
			Data: []map[string]string{
				{
					"title":   title + " Episode 1",
					"episode": "1",
					"url":     "/episode/" + title + "-episode-1",
				},
				{
					"title":   title + " Episode 2",
					"episode": "2", 
					"url":     "/episode/" + title + "-episode-2",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// API schedule endpoint (placeholder)
	http.HandleFunc("/api/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := APIResponse{
			Status:  "success",
			Message: "Jadwal anime berhasil diambil",
			Data: []map[string]string{
				{
					"title": "Sample Anime Schedule 1",
					"day":   "Monday",
					"time":  "20:00",
				},
				{
					"title": "Sample Anime Schedule 2",
					"day":   "Tuesday", 
					"time":  "21:00",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	port := "1408"
	fmt.Printf("Server is running on http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
