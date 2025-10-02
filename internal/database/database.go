package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB interface for database operations
type DB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}

// Client wraps the database connection
type Client struct {
	db DB
}

// NewClient creates a new database client
func NewClient(databaseURL string) (*Client, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{db: db}, nil
}

// NewClientWithDB creates a new client with an existing DB connection
func NewClientWithDB(db DB) *Client {
	return &Client{db: db}
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// Summoner represents a summoner in the database
type Summoner struct {
	ID            int
	PUUID         string
	SummonerName  string
	SummonerLevel int
	ProfileIconID int
	Region        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Match represents a match in the database
type Match struct {
	ID           int
	MatchID      string
	GameMode     string
	GameDuration int
	GameCreation int64
	Data         map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SaveSummoner saves a summoner to the database
func (c *Client) SaveSummoner(summoner *Summoner) error {
	query := `
		INSERT INTO summoners (puuid, summoner_name, summoner_level, profile_icon_id, region, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (puuid) 
		DO UPDATE SET 
			summoner_name = EXCLUDED.summoner_name,
			summoner_level = EXCLUDED.summoner_level,
			profile_icon_id = EXCLUDED.profile_icon_id,
			region = EXCLUDED.region,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`

	err := c.db.QueryRow(
		query,
		summoner.PUUID,
		summoner.SummonerName,
		summoner.SummonerLevel,
		summoner.ProfileIconID,
		summoner.Region,
		time.Now(),
	).Scan(&summoner.ID, &summoner.CreatedAt, &summoner.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to save summoner: %w", err)
	}

	return nil
}

// GetSummonerByPUUID retrieves a summoner by PUUID
func (c *Client) GetSummonerByPUUID(puuid string) (*Summoner, error) {
	query := `
		SELECT id, puuid, summoner_name, summoner_level, profile_icon_id, region, created_at, updated_at
		FROM summoners
		WHERE puuid = $1
	`

	summoner := &Summoner{}
	err := c.db.QueryRow(query, puuid).Scan(
		&summoner.ID,
		&summoner.PUUID,
		&summoner.SummonerName,
		&summoner.SummonerLevel,
		&summoner.ProfileIconID,
		&summoner.Region,
		&summoner.CreatedAt,
		&summoner.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get summoner: %w", err)
	}

	return summoner, nil
}

// SaveMatch saves a match to the database
func (c *Client) SaveMatch(match *Match) error {
	dataJSON, err := json.Marshal(match.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal match data: %w", err)
	}

	query := `
		INSERT INTO matches (match_id, game_mode, game_duration, game_creation, data, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (match_id)
		DO UPDATE SET
			game_mode = EXCLUDED.game_mode,
			game_duration = EXCLUDED.game_duration,
			game_creation = EXCLUDED.game_creation,
			data = EXCLUDED.data,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`

	err = c.db.QueryRow(
		query,
		match.MatchID,
		match.GameMode,
		match.GameDuration,
		match.GameCreation,
		dataJSON,
		time.Now(),
	).Scan(&match.ID, &match.CreatedAt, &match.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to save match: %w", err)
	}

	return nil
}

// GetMatchByID retrieves a match by match ID
func (c *Client) GetMatchByID(matchID string) (*Match, error) {
	query := `
		SELECT id, match_id, game_mode, game_duration, game_creation, data, created_at, updated_at
		FROM matches
		WHERE match_id = $1
	`

	match := &Match{}
	var dataJSON []byte
	err := c.db.QueryRow(query, matchID).Scan(
		&match.ID,
		&match.MatchID,
		&match.GameMode,
		&match.GameDuration,
		&match.GameCreation,
		&dataJSON,
		&match.CreatedAt,
		&match.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	if err := json.Unmarshal(dataJSON, &match.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match data: %w", err)
	}

	return match, nil
}
