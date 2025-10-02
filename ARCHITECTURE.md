# System Architecture

## Overview Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Docker Compose                            │
│                                                                     │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐       │
│  │   PostgreSQL   │  │     Redis      │  │  Application   │       │
│  │   Port: 5432   │  │   Port: 6379   │  │   Port: 8080   │       │
│  │                │  │                │  │                │       │
│  │  ┌──────────┐  │  │  ┌──────────┐  │  │  ┌──────────┐  │       │
│  │  │summoners │  │  │  │ Cache    │  │  │  │  HTTP    │  │       │
│  │  │          │  │  │  │ Store    │  │  │  │  Server  │  │       │
│  │  │matches   │  │  │  │          │  │  │  │          │  │       │
│  │  └──────────┘  │  │  └──────────┘  │  │  └──────────┘  │       │
│  │                │  │                │  │       │         │       │
│  │  Health: ✓     │  │  Health: ✓     │  │       │         │       │
│  └────────────────┘  └────────────────┘  └───────┼─────────┘       │
│         ▲                    ▲                    │                 │
│         │                    │                    │                 │
│         └────────────────────┴────────────────────┘                 │
│                              │                                      │
└──────────────────────────────┼──────────────────────────────────────┘
                               │
                               ▼
                    ┌────────────────────┐
                    │   Client (HTTP)    │
                    │   curl, Browser    │
                    └────────────────────┘
```

## Request Flow Diagram

```
Client Request
     │
     ├─→ GET /health ─────────────────────→ [200 OK]
     │
     ├─→ GET /summoner?region=na1&name=X
     │         │
     │         ├─→ Service Layer
     │         │         │
     │         │         ├─→ Check Redis Cache ──→ [HIT] ──→ Return
     │         │         │                    ↓
     │         │         │                  [MISS]
     │         │         │                    │
     │         │         ├─→ Check PostgreSQL ─→ [HIT] ──→ Cache & Return
     │         │         │                    ↓
     │         │         │                  [MISS]
     │         │         │                    │
     │         │         └─→ Call Riot API ──→ Store in DB & Cache
     │         │                              │
     │         └─────────────────────────────→ Return to Client
     │
     ├─→ GET /matches?puuid=X&count=N
     │         │
     │         └─→ Similar flow with 5min cache
     │
     └─→ GET /match?id=X
               │
               └─→ Similar flow with 30min cache
```

## Component Interaction

```
┌──────────────────────────────────────────────────────────────┐
│                      Application Layer                       │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              cmd/server/main.go                        │  │
│  │  - HTTP Server (net/http)                              │  │
│  │  - Route Handlers                                      │  │
│  │  - Graceful Shutdown                                   │  │
│  └──────────────────────┬─────────────────────────────────┘  │
│                         │                                    │
│  ┌──────────────────────▼─────────────────────────────────┐  │
│  │         internal/service/service.go                    │  │
│  │  - GetSummoner()                                       │  │
│  │  - GetMatches()                                        │  │
│  │  - GetMatch()                                          │  │
│  │  - Cache & DB Coordination                             │  │
│  └─┬──────────────┬──────────────┬────────────────────────┘  │
│    │              │              │                           │
└────┼──────────────┼──────────────┼───────────────────────────┘
     │              │              │
     ▼              ▼              ▼
┌─────────┐   ┌──────────┐   ┌──────────┐
│ Riot    │   │ Database │   │  Cache   │
│ Client  │   │ Client   │   │  Client  │
└─────────┘   └──────────┘   └──────────┘
     │              │              │
     ▼              ▼              ▼
┌─────────┐   ┌──────────┐   ┌──────────┐
│ Riot    │   │PostgreSQL│   │  Redis   │
│  API    │   │   15     │   │    7     │
└─────────┘   └──────────┘   └──────────┘
  External       Docker         Docker
```

## Data Storage Schema

```
PostgreSQL Database: draftvisualizer
│
├── Table: summoners
│   ├── id (SERIAL PRIMARY KEY)
│   ├── puuid (VARCHAR(255) UNIQUE) ◄─── Indexed
│   ├── summoner_name (VARCHAR(255))
│   ├── summoner_level (INTEGER)
│   ├── profile_icon_id (INTEGER)
│   ├── region (VARCHAR(10))           ◄─── Indexed
│   ├── created_at (TIMESTAMP)
│   └── updated_at (TIMESTAMP)
│
└── Table: matches
    ├── id (SERIAL PRIMARY KEY)
    ├── match_id (VARCHAR(255) UNIQUE) ◄─── Indexed
    ├── game_mode (VARCHAR(50))
    ├── game_duration (INTEGER)
    ├── game_creation (BIGINT)         ◄─── Indexed
    ├── data (JSONB)                   ◄─── Full match data
    ├── created_at (TIMESTAMP)
    └── updated_at (TIMESTAMP)

Redis Cache Structure:
│
├── summoner:{region}:{name} → Summoner JSON (TTL: 10min)
├── matches:{puuid}:{count}  → Match IDs JSON (TTL: 5min)
└── match:{matchId}          → Match JSON (TTL: 30min)
```

## Testing Architecture

```
┌────────────────────────────────────────────────────────────┐
│                     Test Pyramid                           │
│                                                            │
│                      ╱╲                                    │
│                     ╱  ╲   E2E Tests                       │
│                    ╱    ╲  (test/e2e)                      │
│                   ╱──────╲                                 │
│                  ╱        ╲                                │
│                 ╱ Integ.   ╲  Integration Tests           │
│                ╱   Tests    ╲ (test/integration)          │
│               ╱──────────────╲                             │
│              ╱                ╲                            │
│             ╱   Unit Tests     ╲  Unit Tests              │
│            ╱   (73.4% coverage) ╱ (*_test.go files)       │
│           ╱──────────────────────╲                        │
│          ╱                        ╲                       │
│         ╱──────────────────────────╲                      │
│                                                            │
└────────────────────────────────────────────────────────────┘

Unit Tests (Mock-Based):
├── config_test.go      (100% coverage)
├── client_test.go      (82.4% coverage)
├── database_test.go    (72.5% coverage)
├── cache_test.go       (65.7% coverage)
└── service_test.go     (62.5% coverage)

Integration Tests (Real Services):
├── Database Integration (PostgreSQL)
└── Cache Integration (Redis)

E2E Tests:
└── HTTP Server Validation
```

## Deployment Architecture

```
Development Environment:
┌────────────────────────────────────────────┐
│ Local Machine                              │
│                                            │
│  make docker-up                            │
│      │                                     │
│      ├─→ PostgreSQL (localhost:5432)      │
│      ├─→ Redis (localhost:6379)           │
│      └─→ App (localhost:8080)             │
│                                            │
└────────────────────────────────────────────┘

Production Environment (Example):
┌────────────────────────────────────────────┐
│ Cloud Provider (AWS/GCP/Azure)            │
│                                            │
│  ┌──────────────┐    ┌─────────────┐      │
│  │ Load Balancer│◄───┤ Application │      │
│  │ (ALB/nginx)  │    │  Container  │      │
│  └──────────────┘    └─────────────┘      │
│         │                    │             │
│         │            ┌───────┴────────┐    │
│         │            │                │    │
│    ┌────▼───┐   ┌───▼────┐    ┌─────▼──┐  │
│    │  App   │   │  RDS   │    │ Redis  │  │
│    │Instance│   │  (PG)  │    │ Elasti │  │
│    └────────┘   └────────┘    │ Cache  │  │
│                               └────────┘  │
└────────────────────────────────────────────┘
```

## Caching Strategy Flowchart

```
                    Receive Request
                          │
                          ▼
                    ┌──────────┐
                    │ Check    │──→ HIT ──→ Return Cached Data
                    │ Redis    │           (5-30min old)
                    └──────────┘
                          │
                        MISS
                          │
                          ▼
                    ┌──────────┐
                    │ Check    │──→ HIT ──→ Update Cache
                    │ Database │           Return DB Data
                    └──────────┘
                          │
                        MISS
                          │
                          ▼
                    ┌──────────┐
                    │  Call    │──→ SUCCESS ──→ Save to DB
                    │ Riot API │               Save to Cache
                    └──────────┘               Return Data
                          │
                        ERROR
                          │
                          ▼
                    Return Error
```

## Error Handling Flow

```
HTTP Request
     │
     ├─→ Invalid Parameters ──→ 400 Bad Request
     │
     ├─→ Cache Error
     │         │
     │         ├─→ Log Error
     │         └─→ Continue to Database
     │
     ├─→ Database Error
     │         │
     │         ├─→ Log Error
     │         └─→ Continue to API
     │
     ├─→ API Error
     │         │
     │         ├─→ Log Error
     │         └─→ 500 Internal Server Error
     │
     └─→ Success ──→ 200 OK with Data
```

## Security Layers

```
┌─────────────────────────────────────────┐
│         Security Considerations         │
│                                         │
│  1. Environment Variables               │
│     - API keys not in code              │
│     - .env file (gitignored)            │
│                                         │
│  2. Database Security                   │
│     - Connection pooling                │
│     - Prepared statements               │
│     - No SQL injection                  │
│                                         │
│  3. Docker Isolation                    │
│     - Network isolation                 │
│     - Health checks                     │
│     - Resource limits                   │
│                                         │
│  4. Application                         │
│     - Input validation                  │
│     - Error sanitization                │
│     - Graceful shutdown                 │
│                                         │
└─────────────────────────────────────────┘
```

This architecture provides:
- ✅ High performance (caching)
- ✅ Reliability (multi-layer fallback)
- ✅ Scalability (stateless design)
- ✅ Maintainability (clean separation)
- ✅ Testability (comprehensive tests)
