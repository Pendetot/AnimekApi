package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type APIResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Status:    "ok",
		Message:   "AnoBoy API v2.0 - Golang Version is running",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

func apiInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	endpoints := []map[string]string{
		{"path": "/api/detail?slug=<anime_slug>", "description": "Get anime episode details by slug"},
		{"path": "/api/list", "description": "Get complete anime list"},
		{"path": "/api/episodes?title=<anime_title>", "description": "Get episodes list for an anime by title"},
		{"path": "/api/schedule", "description": "Get anime broadcast schedule with images"},
	}
	
	response := APIResponse{
		Status:  "success",
		Message: "AnoBoy API is running",
		Data: map[string]interface{}{
			"version":   "2.0.0",
			"endpoints": endpoints,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/", apiInfoHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			apiInfoHandler(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIResponse{
				Status:  "error",
				Message: "Route " + r.URL.Path + " not found",
			})
		}
	})

	port := "1408"
	fmt.Printf("Server starting on http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
