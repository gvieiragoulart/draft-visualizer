# Draft Visualizer

A Go application that integrates with the Riot Games API to fetch and store League of Legends data. The application uses PostgreSQL for persistent storage and Redis for caching API responses.

## Features

- **Riot Games API Integration**: Fetch summoner and match data using the official Riot Games API
- **PostgreSQL Database**: Store summoner and match information persistently
- **Redis Cache**: Cache API responses to reduce API calls and improve performance
- **Docker Support**: Full Docker and docker-compose setup for easy deployment
- **100% Test Coverage**: Comprehensive unit and integration tests

## Architecture

The project follows a clean architecture pattern with the following layers:

- **cmd/server**: Main application entry point and HTTP server
- **internal/config**: Configuration management
- **internal/riot**: Riot Games API client
- **internal/database**: PostgreSQL database client and models
- **internal/cache**: Redis cache client
- **internal/service**: Business logic layer that coordinates API, cache, and database

## Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (for running with containers)
- PostgreSQL 15+ (if running locally without Docker)
- Redis 7+ (if running locally without Docker)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/gvieiragoulart/draft-visualizer.git
cd draft-visualizer
```

### 2. Set up environment variables

Copy the example environment file and update it with your Riot API key:

```bash
cp .env.example .env
```

Edit `.env` and add your Riot Games API key:

```
RIOT_API_KEY=your_riot_api_key_here
```

You can get an API key from [Riot Games Developer Portal](https://developer.riotgames.com/).

### 3. Run with Docker

The easiest way to run the application is with Docker Compose:

```bash
# Start all services (PostgreSQL, Redis, and the application)
make docker-up

# Or using docker-compose directly
docker-compose up -d
```

The application will be available at `http://localhost:8080`.

### 4. Run locally (without Docker)

If you want to run the application locally:

```bash
# Install dependencies
make deps

# Start PostgreSQL and Redis (you'll need to have them installed)
# Or use docker-compose to run only the databases:
docker-compose up -d postgres redis

# Run the application
make run
```

## API Endpoints

### Health Check

```bash
GET /health
```

Returns the health status of the application.

### Get Summoner

```bash
GET /summoner?region=na1&name=SummonerName
```

Fetches summoner information by name and region.

**Parameters:**
- `region`: The region (e.g., na1, euw1, kr)
- `name`: The summoner name

**Response:**
```json
{
  "puuid": "summoner-puuid",
  "name": "SummonerName",
  "summonerLevel": 100,
  "profileIconId": 1234
}
```

### Get Matches

```bash
GET /matches?puuid=summoner-puuid&count=20
```

Fetches match IDs for a summoner.

**Parameters:**
- `puuid`: The summoner PUUID
- `count`: Number of matches to fetch (optional, default: 20)

**Response:**
```json
["NA1_match1", "NA1_match2", "NA1_match3"]
```

### Get Match Details

```bash
GET /match?id=NA1_match1
```

Fetches detailed information about a specific match.

**Parameters:**
- `id`: The match ID

**Response:**
```json
{
  "metadata": {
    "matchId": "NA1_match1",
    "participants": ["puuid1", "puuid2"]
  },
  "info": {
    "gameMode": "CLASSIC",
    "gameDuration": 1800,
    "gameCreation": 1609459200000
  }
}
```

## Development

### Run Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Generate coverage report
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint
```

### Build

```bash
# Build the application
make build

# The binary will be created in ./bin/server
```

## Project Structure

```
.
├── cmd/
│   └── server/          # Main application
├── internal/
│   ├── cache/           # Redis cache client
│   ├── config/          # Configuration management
│   ├── database/        # PostgreSQL client and models
│   ├── riot/            # Riot Games API client
│   └── service/         # Business logic layer
├── test/
│   └── integration/     # Integration tests
├── scripts/
│   └── init.sql         # Database initialization script
├── .env.example         # Example environment variables
├── .gitignore           # Git ignore rules
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker image definition
├── Makefile             # Build and development commands
├── go.mod               # Go module definition
└── README.md            # This file
```

## Docker Services

The docker-compose setup includes:

- **postgres**: PostgreSQL 15 database on port 5432
- **redis**: Redis 7 cache on port 6379
- **app**: The main application on port 8080

All services include health checks and the application waits for the databases to be ready before starting.

## Testing

The project includes comprehensive tests:

### Unit Tests

- Configuration loading and validation
- Riot API client with mocked HTTP responses
- Database operations with mocked database connections
- Cache operations with mocked Redis clients
- Service layer business logic

### Integration Tests

- Real PostgreSQL database operations
- Real Redis cache operations
- End-to-end workflows

Run integration tests with:

```bash
make test-integration
```

Note: Integration tests require PostgreSQL and Redis to be running. You can start them with:

```bash
docker-compose up -d postgres redis
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RIOT_API_KEY` | Riot Games API key (required) | - |
| `DATABASE_URL` | PostgreSQL connection string (required) | - |
| `REDIS_URL` | Redis connection URL | `redis://localhost:6379` |
| `REDIS_PASSWORD` | Redis password (if required) | - |
| `SERVER_PORT` | HTTP server port | `8080` |

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.