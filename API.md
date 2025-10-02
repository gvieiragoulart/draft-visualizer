# API Documentation

## Base URL

When running locally with Docker:
```
http://localhost:8080
```

## Authentication

The application requires a Riot Games API key to be configured via the `RIOT_API_KEY` environment variable. This is used internally to fetch data from the Riot Games API. No authentication is required for the client to use this service.

## Endpoints

### Health Check

Check if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok"
}
```

**Status Codes:**
- `200 OK`: Service is healthy

**Example:**
```bash
curl http://localhost:8080/health
```

---

### Get Summoner

Retrieve summoner information by name and region.

**Endpoint:** `GET /summoner`

**Query Parameters:**
| Parameter | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| region    | string | Yes      | Region code (e.g., na1, euw1)  |
| name      | string | Yes      | Summoner name                  |

**Response:**
```json
{
  "puuid": "string",
  "name": "string",
  "summonerLevel": 0,
  "profileIconId": 0
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Missing required parameters
- `500 Internal Server Error`: Server error

**Example:**
```bash
curl "http://localhost:8080/summoner?region=na1&name=Faker"
```

**Response Example:**
```json
{
  "puuid": "xyz123abc456def789...",
  "name": "Faker",
  "summonerLevel": 500,
  "profileIconId": 4568
}
```

**Caching:**
- Cached for 10 minutes
- Stored in database permanently
- Updates database on each request

**Supported Regions:**
- `br1`: Brazil
- `eun1`: Europe Nordic & East
- `euw1`: Europe West
- `jp1`: Japan
- `kr`: Korea
- `la1`: Latin America North
- `la2`: Latin America South
- `na1`: North America
- `oc1`: Oceania
- `ru`: Russia
- `tr1`: Turkey
- `ph2`: Philippines
- `sg2`: Singapore
- `th2`: Thailand
- `tw2`: Taiwan
- `vn2`: Vietnam

---

### Get Matches

Retrieve match IDs for a summoner.

**Endpoint:** `GET /matches`

**Query Parameters:**
| Parameter | Type   | Required | Description                           |
|-----------|--------|----------|---------------------------------------|
| puuid     | string | Yes      | Player Universally Unique Identifier  |
| count     | number | No       | Number of matches (default: 20)       |

**Response:**
```json
[
  "NA1_1234567890",
  "NA1_0987654321",
  "NA1_5555555555"
]
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Missing required parameters
- `500 Internal Server Error`: Server error

**Example:**
```bash
curl "http://localhost:8080/matches?puuid=xyz123abc456&count=5"
```

**Response Example:**
```json
[
  "NA1_4567890123",
  "NA1_4567890124",
  "NA1_4567890125",
  "NA1_4567890126",
  "NA1_4567890127"
]
```

**Caching:**
- Cached for 5 minutes
- Not stored in database (only IDs, data fetched separately)

**Notes:**
- Maximum count is determined by Riot API limits (typically 100)
- Results are ordered by most recent first
- PUUID can be obtained from the `/summoner` endpoint

---

### Get Match Details

Retrieve detailed information about a specific match.

**Endpoint:** `GET /match`

**Query Parameters:**
| Parameter | Type   | Required | Description                     |
|-----------|--------|----------|----------------------------------|
| id        | string | Yes      | Match ID (e.g., NA1_1234567890) |

**Response:**
```json
{
  "metadata": {
    "matchId": "string",
    "participants": ["string"]
  },
  "info": {
    "gameMode": "string",
    "gameDuration": 0,
    "gameCreation": 0
  }
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Missing required parameters
- `500 Internal Server Error`: Server error

**Example:**
```bash
curl "http://localhost:8080/match?id=NA1_4567890123"
```

**Response Example:**
```json
{
  "metadata": {
    "matchId": "NA1_4567890123",
    "participants": [
      "puuid1",
      "puuid2",
      "puuid3",
      "puuid4",
      "puuid5",
      "puuid6",
      "puuid7",
      "puuid8",
      "puuid9",
      "puuid10"
    ]
  },
  "info": {
    "gameMode": "CLASSIC",
    "gameDuration": 1856,
    "gameCreation": 1696000000000
  }
}
```

**Caching:**
- Cached for 30 minutes
- Stored in database permanently
- Full match data stored in database as JSONB

**Game Modes:**
- `CLASSIC`: Summoner's Rift
- `ARAM`: All Random All Mid
- `URF`: Ultra Rapid Fire
- `TUTORIAL`: Tutorial games
- And others...

---

## Error Responses

All endpoints may return the following error format:

```json
{
  "error": "Error message description"
}
```

### Common Error Status Codes

- `400 Bad Request`: Invalid or missing parameters
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded (from Riot API)
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service temporarily unavailable

---

## Rate Limits

The service respects Riot Games API rate limits:
- Development API Key: 20 requests per second, 100 requests per 2 minutes
- Production API Key: Higher limits based on your approval

The service implements caching to minimize API calls and stay within rate limits.

---

## Data Flow

1. **Request** → Client makes HTTP request to the service
2. **Cache Check** → Service checks Redis cache for existing data
3. **Database Check** → If not in cache, check PostgreSQL database
4. **API Call** → If not in database, fetch from Riot Games API
5. **Store** → Save data to cache and database
6. **Response** → Return data to client

This multi-layer approach ensures:
- Fast response times (cache hits)
- Reduced API calls (rate limit compliance)
- Data persistence (database storage)
- High availability (cached and stored data)

---

## Examples

### Complete Workflow Example

1. Get summoner information:
```bash
curl "http://localhost:8080/summoner?region=na1&name=Faker"
```

2. Use PUUID from response to get matches:
```bash
curl "http://localhost:8080/matches?puuid=xyz123abc456&count=10"
```

3. Get details for a specific match:
```bash
curl "http://localhost:8080/match?id=NA1_4567890123"
```

### Using with jq (JSON processor)

```bash
# Pretty print summoner info
curl -s "http://localhost:8080/summoner?region=na1&name=Faker" | jq '.'

# Extract just the PUUID
curl -s "http://localhost:8080/summoner?region=na1&name=Faker" | jq -r '.puuid'

# Get match IDs and format nicely
curl -s "http://localhost:8080/matches?puuid=xyz123abc456&count=5" | jq '.[]'
```

### Using with curl in a script

```bash
#!/bin/bash

REGION="na1"
SUMMONER_NAME="Faker"
SERVER="http://localhost:8080"

# Get summoner
echo "Fetching summoner: $SUMMONER_NAME"
SUMMONER=$(curl -s "$SERVER/summoner?region=$REGION&name=$SUMMONER_NAME")
echo $SUMMONER | jq '.'

# Extract PUUID
PUUID=$(echo $SUMMONER | jq -r '.puuid')
echo "PUUID: $PUUID"

# Get recent matches
echo "Fetching recent matches..."
MATCHES=$(curl -s "$SERVER/matches?puuid=$PUUID&count=5")
echo $MATCHES | jq '.'
```

---

## Testing

### Using Postman

1. Import the following collection or create requests manually
2. Set base URL as variable: `{{baseUrl}}` = `http://localhost:8080`
3. Create requests for each endpoint

### Using HTTPie

```bash
# Health check
http GET localhost:8080/health

# Get summoner
http GET localhost:8080/summoner region==na1 name==Faker

# Get matches
http GET localhost:8080/matches puuid==xyz123abc456 count==10

# Get match details
http GET localhost:8080/match id==NA1_4567890123
```

---

## Support

For issues or questions:
1. Check the logs: `docker-compose logs -f app`
2. Verify databases are running: `docker-compose ps`
3. Check Riot API status: https://developer.riotgames.com/
4. Review your API key permissions

For more information, see the [README.md](README.md) and [DEVELOPMENT.md](DEVELOPMENT.md).
