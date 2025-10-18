package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	RiotAPIKey    string
	EsportsAPIKey string
	DatabaseURL   string
	RedisURL      string
	RedisPassword string
	ServerPort    int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	riotAPIKey := os.Getenv("RIOT_API_KEY")
	if riotAPIKey == "" {
		return nil, fmt.Errorf("RIOT_API_KEY environment variable is required")
	}

	esportsAPIKey := os.Getenv("ESPORTS_API_KEY")
	if esportsAPIKey == "" {
		return nil, fmt.Errorf("ESPORTS_API_KEY environment variable is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	serverPortStr := os.Getenv("SERVER_PORT")
	if serverPortStr == "" {
		serverPortStr = "8080"
	}

	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	return &Config{
		RiotAPIKey:    riotAPIKey,
		EsportsAPIKey: esportsAPIKey,
		DatabaseURL:   databaseURL,
		RedisURL:      redisURL,
		RedisPassword: redisPassword,
		ServerPort:    serverPort,
	}, nil
}
