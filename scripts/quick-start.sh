#!/bin/bash
set -e

echo "=== Draft Visualizer - Quick Start Guide ==="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "Step 1: Starting PostgreSQL and Redis..."
docker-compose up -d postgres redis

echo "Waiting for databases to be ready..."
sleep 5

echo ""
echo "Step 2: Running unit tests..."
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./internal/...

echo ""
echo "Step 3: Checking test coverage..."
go tool cover -func=coverage.txt | tail -1

echo ""
echo "Step 4: Building the application..."
go build -o bin/server ./cmd/server

echo ""
echo "=== Build Complete! ==="
echo ""
echo "To run the application:"
echo "  1. Set your RIOT_API_KEY in .env file (copy from .env.example)"
echo "  2. Run: make docker-up"
echo "  3. Access: http://localhost:8080/health"
echo ""
echo "To run integration tests:"
echo "  1. Ensure databases are running: docker-compose up -d postgres redis"
echo "  2. Run: make test-integration"
echo ""
echo "For more commands, run: make help"
