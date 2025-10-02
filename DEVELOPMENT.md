# Development Workflow

## Quick Start

1. **Clone and setup**:
```bash
git clone https://github.com/gvieiragoulart/draft-visualizer.git
cd draft-visualizer
cp .env.example .env
# Edit .env and add your RIOT_API_KEY
```

2. **Run with Docker** (Recommended):
```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f app

# Stop services
docker-compose down
```

3. **Access the API**:
```bash
# Health check
curl http://localhost:8080/health

# Get summoner (replace YOUR_API_KEY and test with real summoner name)
curl "http://localhost:8080/summoner?region=na1&name=Faker"
```

## Development Commands

### Testing

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests (requires databases running)
make test-integration

# Generate coverage report
make test-coverage
# Open coverage.html in browser
```

### Building

```bash
# Build the application
make build

# Run the application locally
make run

# Format code
make fmt

# Run go vet
make vet
```

### Docker

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs

# Rebuild docker image
make docker-build
```

## Project Architecture

```
draft-visualizer/
├── cmd/
│   └── server/          # HTTP server entry point
├── internal/
│   ├── cache/           # Redis cache layer
│   ├── config/          # Configuration management
│   ├── database/        # PostgreSQL database layer
│   ├── riot/            # Riot Games API client
│   └── service/         # Business logic layer
├── test/
│   ├── integration/     # Integration tests
│   └── e2e/            # End-to-end tests
├── scripts/
│   ├── init.sql        # Database initialization
│   └── quick-start.sh  # Quick start script
└── docker-compose.yml   # Docker services
```

## Testing Strategy

### Unit Tests (73.4% coverage)
- Config: 100%
- Riot API Client: 82.4%
- Database: 72.5%
- Cache: 65.7%
- Service: 62.5%

All unit tests use mocks and don't require external dependencies.

### Integration Tests
Test real interactions with PostgreSQL and Redis using Docker containers.

Run with:
```bash
docker-compose up -d postgres redis
make test-integration
```

### End-to-End Tests
Test the complete application workflow including HTTP endpoints.

## API Flow

1. **Request comes in** → HTTP handler in `cmd/server/main.go`
2. **Service layer** → `internal/service/` coordinates the flow
3. **Check cache** → `internal/cache/` checks Redis for cached data
4. **Check database** → `internal/database/` checks PostgreSQL
5. **Fetch from API** → `internal/riot/` calls Riot Games API
6. **Store results** → Save to cache and database
7. **Return response** → Send JSON response to client

## Caching Strategy

- **Summoner data**: Cached for 10 minutes
- **Match IDs**: Cached for 5 minutes  
- **Match details**: Cached for 30 minutes

All API responses are also stored in PostgreSQL for persistence.

## Database Schema

### Summoners Table
```sql
CREATE TABLE summoners (
    id SERIAL PRIMARY KEY,
    puuid VARCHAR(255) UNIQUE NOT NULL,
    summoner_name VARCHAR(255) NOT NULL,
    summoner_level INTEGER NOT NULL,
    profile_icon_id INTEGER NOT NULL,
    region VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Matches Table
```sql
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    match_id VARCHAR(255) UNIQUE NOT NULL,
    game_mode VARCHAR(50),
    game_duration INTEGER,
    game_creation BIGINT,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Troubleshooting

### Port Already in Use
If port 8080 is already in use:
```bash
# Change SERVER_PORT in .env
SERVER_PORT=8081

# Or change in docker-compose.yml
ports:
  - "8081:8080"
```

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Redis Connection Failed
```bash
# Check if Redis is running
docker-compose ps redis

# View Redis logs
docker-compose logs redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

### Cannot Connect to Riot API
1. Check your API key in `.env`
2. Verify the summoner name and region are correct
3. Check Riot API status: https://developer.riotgames.com/
4. Review rate limits in your API key dashboard

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Submit a pull request

## Additional Resources

- [Riot Games API Documentation](https://developer.riotgames.com/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
