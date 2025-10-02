package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var serverURL = flag.String("url", "http://localhost:8080", "Server URL to test")

func main() {
	flag.Parse()

	fmt.Println("=== Draft Visualizer E2E Test ===")
	fmt.Printf("Testing server at: %s\n\n", *serverURL)

	// Test 1: Health check
	if err := testHealthCheck(); err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Health check passed")

	fmt.Println("\n=== All E2E tests passed! ===")
}

func testHealthCheck() error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(*serverURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result["status"] != "ok" {
		return fmt.Errorf("expected status 'ok', got '%s'", result["status"])
	}

	return nil
}
