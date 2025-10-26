package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Set required environment variables
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.RiotAPIKey != "test-api-key" {
		t.Errorf("expected RiotAPIKey to be 'test-api-key', got %s", cfg.RiotAPIKey)
	}

	if cfg.DatabaseURL != "postgres://localhost:5432/test" {
		t.Errorf("expected DatabaseURL to be 'postgres://localhost:5432/test', got %s", cfg.DatabaseURL)
	}

	if cfg.RedisURL != "redis://localhost:6379" {
		t.Errorf("expected RedisURL to be 'redis://localhost:6379', got %s", cfg.RedisURL)
	}

	if cfg.ServerPort != 8080 {
		t.Errorf("expected ServerPort to be 8080, got %d", cfg.ServerPort)
	}
}

func TestLoad_MissingRiotAPIKey(t *testing.T) {
	os.Unsetenv("RIOT_API_KEY")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing RIOT_API_KEY, got nil")
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	os.Unsetenv("DATABASE_URL")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL, got nil")
	}
}

func TestLoad_MissingEsportsAPIKey(t *testing.T) {
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Unsetenv("ESPORTS_API_KEY")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing ESPORTS_API_KEY, got nil")
	}
}

func TestLoad_MissingWikiCredentials(t *testing.T) {
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Unsetenv("WIKI_USERNAME")
	os.Unsetenv("WIKI_PASSWORD")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("DATABASE_URL")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing WIKI_USERNAME and WIKI_PASSWORD, got nil")
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	os.Unsetenv("REDIS_URL")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("REDIS_PASSWORD")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.RedisURL != "redis://localhost:6379" {
		t.Errorf("expected default RedisURL to be 'redis://localhost:6379', got %s", cfg.RedisURL)
	}

	if cfg.ServerPort != 8080 {
		t.Errorf("expected default ServerPort to be 8080, got %d", cfg.ServerPort)
	}

	if cfg.RedisPassword != "" {
		t.Errorf("expected RedisPassword to be empty, got %s", cfg.RedisPassword)
	}
}

func TestLoad_InvalidServerPort(t *testing.T) {
	os.Setenv("RIOT_API_KEY", "test-api-key")
	os.Setenv("ESPORTS_API_KEY", "test-esports-api-key")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/test")
	os.Setenv("WIKI_USERNAME", "test-user")
	os.Setenv("WIKI_PASSWORD", "test-password")
	os.Setenv("SERVER_PORT", "invalid")
	defer func() {
		os.Unsetenv("RIOT_API_KEY")
		os.Unsetenv("ESPORTS_API_KEY")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("WIKI_USERNAME")
		os.Unsetenv("WIKI_PASSWORD")
		os.Unsetenv("SERVER_PORT")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid SERVER_PORT, got nil")
	}
}
