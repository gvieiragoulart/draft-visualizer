.PHONY: help build test test-unit test-integration test-coverage run docker-up docker-down clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/server ./cmd/server

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./internal/...

test-integration: ## Run integration tests
	go test -v -race ./test/integration/...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

run: ## Run the application locally
	go run ./cmd/server

docker-up: ## Start all services with docker-compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-build: ## Build docker image
	docker-compose build

docker-logs: ## Show docker logs
	docker-compose logs -f

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.txt coverage.html
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run || true
