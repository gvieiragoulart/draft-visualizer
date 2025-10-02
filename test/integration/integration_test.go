package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/gvieiragoulart/draft-visualizer/internal/cache"
	"github.com/gvieiragoulart/draft-visualizer/internal/database"
	_ "github.com/lib/pq"
)

// TestDatabaseIntegration tests database operations with a real PostgreSQL instance
func TestDatabaseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use environment variable or default test database URL
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/draftvisualizer?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skipf("Failed to ping database: %v", err)
	}

	client := database.NewClientWithDB(db)

	// Test saving and retrieving a summoner
	t.Run("SaveAndGetSummoner", func(t *testing.T) {
		summoner := &database.Summoner{
			PUUID:         "test-puuid-integration",
			SummonerName:  "TestIntegrationSummoner",
			SummonerLevel: 150,
			ProfileIconID: 5678,
			Region:        "na1",
		}

		// Save summoner
		err := client.SaveSummoner(summoner)
		if err != nil {
			t.Fatalf("Failed to save summoner: %v", err)
		}

		// Retrieve summoner
		retrieved, err := client.GetSummonerByPUUID("test-puuid-integration")
		if err != nil {
			t.Fatalf("Failed to get summoner: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected to retrieve summoner, got nil")
		}

		if retrieved.PUUID != summoner.PUUID {
			t.Errorf("Expected PUUID %s, got %s", summoner.PUUID, retrieved.PUUID)
		}

		if retrieved.SummonerName != summoner.SummonerName {
			t.Errorf("Expected SummonerName %s, got %s", summoner.SummonerName, retrieved.SummonerName)
		}
	})

	// Test saving and retrieving a match
	t.Run("SaveAndGetMatch", func(t *testing.T) {
		match := &database.Match{
			MatchID:      "NA1_test_integration_match",
			GameMode:     "CLASSIC",
			GameDuration: 2400,
			GameCreation: time.Now().Unix(),
			Data: map[string]interface{}{
				"test": "integration data",
			},
		}

		// Save match
		err := client.SaveMatch(match)
		if err != nil {
			t.Fatalf("Failed to save match: %v", err)
		}

		// Retrieve match
		retrieved, err := client.GetMatchByID("NA1_test_integration_match")
		if err != nil {
			t.Fatalf("Failed to get match: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected to retrieve match, got nil")
		}

		if retrieved.MatchID != match.MatchID {
			t.Errorf("Expected MatchID %s, got %s", match.MatchID, retrieved.MatchID)
		}

		if retrieved.GameMode != match.GameMode {
			t.Errorf("Expected GameMode %s, got %s", match.GameMode, retrieved.GameMode)
		}
	})
}

// TestCacheIntegration tests cache operations with a real Redis instance
func TestCacheIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use environment variable or default test Redis URL
	redisURL := os.Getenv("TEST_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	// Connect to Redis
	client, err := cache.NewClient(redisURL, "")
	if err != nil {
		t.Skipf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Test Set and Get
	t.Run("SetAndGet", func(t *testing.T) {
		key := "test:integration:key"
		value := map[string]interface{}{
			"name":  "test",
			"value": 123,
		}

		// Set value
		err := client.Set(ctx, key, value, 1*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Get value
		var retrieved map[string]interface{}
		found, err := client.GetJSON(ctx, key, &retrieved)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}

		if !found {
			t.Fatal("Expected to find value in cache")
		}

		if retrieved["name"] != "test" {
			t.Errorf("Expected name to be 'test', got %v", retrieved["name"])
		}

		// Clean up
		client.Delete(ctx, key)
	})

	// Test expiration
	t.Run("Expiration", func(t *testing.T) {
		key := "test:integration:expire"
		value := "expire-test"

		// Set value with short expiration
		err := client.Set(ctx, key, value, 1*time.Second)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Value should exist immediately
		data, err := client.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if data == nil {
			t.Fatal("Expected value to exist")
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Value should be gone
		data, err = client.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if data != nil {
			t.Error("Expected value to be expired")
		}
	})

	// Test delete
	t.Run("Delete", func(t *testing.T) {
		key := "test:integration:delete"
		value := "delete-test"

		// Set value
		err := client.Set(ctx, key, value, 1*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Delete value
		err = client.Delete(ctx, key)
		if err != nil {
			t.Fatalf("Failed to delete value: %v", err)
		}

		// Value should be gone
		data, err := client.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if data != nil {
			t.Error("Expected value to be deleted")
		}
	})
}
