package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gvieiragoulart/draft-visualizer/internal/cache"
	"github.com/gvieiragoulart/draft-visualizer/internal/database"
	"github.com/gvieiragoulart/draft-visualizer/internal/riot"
)

// Service provides business logic for the application
type Service struct {
	riotClient *riot.Client
	db         *database.Client
	cache      *cache.Client
}

// NewService creates a new service
func NewService(riotClient *riot.Client, db *database.Client, cache *cache.Client) *Service {
	return &Service{
		riotClient: riotClient,
		db:         db,
		cache:      cache,
	}
}

// GetSummoner retrieves a summoner by name, using cache and database
func (s *Service) GetSummoner(ctx context.Context, region, summonerName string) (*riot.Summoner, error) {
	cacheKey := fmt.Sprintf("summoner:%s:%s", region, summonerName)

	// Try cache first
	var cachedSummoner riot.Summoner
	found, err := s.cache.GetJSON(ctx, cacheKey, &cachedSummoner)
	if err != nil {
		// Log error but continue to API/DB
		fmt.Printf("cache error: %v\n", err)
	}
	if found {
		return &cachedSummoner, nil
	}

	// Try API
	summoner, err := s.riotClient.GetSummonerByName(region, summonerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner from API: %w", err)
	}

	// Store in database
	dbSummoner := &database.Summoner{
		PUUID:         summoner.PUUID,
		SummonerName:  summoner.SummonerName,
		SummonerLevel: summoner.SummonerLevel,
		ProfileIconID: summoner.ProfileIconID,
		Region:        region,
	}
	if err := s.db.SaveSummoner(dbSummoner); err != nil {
		// Log error but continue
		fmt.Printf("database error: %v\n", err)
	}

	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, summoner, 10*time.Minute); err != nil {
		// Log error but continue
		fmt.Printf("cache error: %v\n", err)
	}

	return summoner, nil
}

// GetMatches retrieves matches for a summoner by PUUID
func (s *Service) GetMatches(ctx context.Context, puuid string, count int) ([]string, error) {
	cacheKey := fmt.Sprintf("matches:%s:%d", puuid, count)

	// Try cache first
	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		fmt.Printf("cache error: %v\n", err)
	}
	if data != nil {
		var matchIDs []string
		if err := json.Unmarshal(data, &matchIDs); err == nil {
			return matchIDs, nil
		}
	}

	// Get from API
	matchIDs, err := s.riotClient.GetMatchesByPUUID(puuid, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches from API: %w", err)
	}

	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, matchIDs, 5*time.Minute); err != nil {
		fmt.Printf("cache error: %v\n", err)
	}

	return matchIDs, nil
}

// GetMatch retrieves a match by ID, using cache and database
func (s *Service) GetMatch(ctx context.Context, matchID string) (*riot.Match, error) {
	cacheKey := fmt.Sprintf("match:%s", matchID)

	// Try cache first
	var cachedMatch riot.Match
	found, err := s.cache.GetJSON(ctx, cacheKey, &cachedMatch)
	if err != nil {
		fmt.Printf("cache error: %v\n", err)
	}
	if found {
		return &cachedMatch, nil
	}

	// Try database
	dbMatch, err := s.db.GetMatchByID(matchID)
	if err != nil {
		fmt.Printf("database error: %v\n", err)
	}
	if dbMatch != nil {
		// Reconstruct match from database
		match := &riot.Match{
			Metadata: riot.MatchMetadata{
				MatchID: dbMatch.MatchID,
			},
			Info: riot.MatchInfo{
				GameMode:     dbMatch.GameMode,
				GameDuration: dbMatch.GameDuration,
				GameCreation: dbMatch.GameCreation,
			},
		}

		// Store in cache
		if err := s.cache.Set(ctx, cacheKey, match, 30*time.Minute); err != nil {
			fmt.Printf("cache error: %v\n", err)
		}

		return match, nil
	}

	// Get from API
	match, err := s.riotClient.GetMatchByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match from API: %w", err)
	}

	// Store in database
	matchData := map[string]interface{}{
		"metadata": match.Metadata,
		"info":     match.Info,
	}
	dbMatch = &database.Match{
		MatchID:      match.Metadata.MatchID,
		GameMode:     match.Info.GameMode,
		GameDuration: match.Info.GameDuration,
		GameCreation: match.Info.GameCreation,
		Data:         matchData,
	}
	if err := s.db.SaveMatch(dbMatch); err != nil {
		fmt.Printf("database error: %v\n", err)
	}

	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, match, 30*time.Minute); err != nil {
		fmt.Printf("cache error: %v\n", err)
	}

	return match, nil
}
