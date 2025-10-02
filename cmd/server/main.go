package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gvieiragoulart/draft-visualizer/internal/cache"
	"github.com/gvieiragoulart/draft-visualizer/internal/config"
	"github.com/gvieiragoulart/draft-visualizer/internal/database"
	"github.com/gvieiragoulart/draft-visualizer/internal/riot"
	"github.com/gvieiragoulart/draft-visualizer/internal/service"
)

type Server struct {
	service *service.Service
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Riot API client
	riotClient := riot.NewClient(cfg.RiotAPIKey)

	// Initialize database client
	dbClient, err := database.NewClient(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	// Initialize cache client
	cacheClient, err := cache.NewClient(cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %v", err)
	}
	defer cacheClient.Close()

	// Initialize service
	svc := service.NewService(riotClient, dbClient, cacheClient)

	// Create server
	server := &Server{service: svc}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.healthHandler)
	mux.HandleFunc("/summoner", server.summonerHandler)
	mux.HandleFunc("/matches", server.matchesHandler)
	mux.HandleFunc("/match", server.matchHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) summonerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	region := r.URL.Query().Get("region")
	name := r.URL.Query().Get("name")

	if region == "" || name == "" {
		http.Error(w, "region and name parameters are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	summoner, err := s.service.GetSummoner(ctx, region, name)
	if err != nil {
		log.Printf("Error getting summoner: %v", err)
		http.Error(w, fmt.Sprintf("Error getting summoner: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summoner)
}

func (s *Server) matchesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	puuid := r.URL.Query().Get("puuid")
	if puuid == "" {
		http.Error(w, "puuid parameter is required", http.StatusBadRequest)
		return
	}

	count := 20 // default
	if countParam := r.URL.Query().Get("count"); countParam != "" {
		fmt.Sscanf(countParam, "%d", &count)
	}

	ctx := r.Context()
	matches, err := s.service.GetMatches(ctx, puuid, count)
	if err != nil {
		log.Printf("Error getting matches: %v", err)
		http.Error(w, fmt.Sprintf("Error getting matches: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func (s *Server) matchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	matchID := r.URL.Query().Get("id")
	if matchID == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	match, err := s.service.GetMatch(ctx, matchID)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		http.Error(w, fmt.Sprintf("Error getting match: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}
