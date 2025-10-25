package service

import (
	"context"
	"fmt"

	"github.com/gvieiragoulart/draft-visualizer/internal/riot"
)

// Service provides business logic for the application
type Service struct {
	riotClient *riot.Client
}

// NewService creates a new service
func NewService(riotClient *riot.Client) *Service {
	return &Service{
		riotClient: riotClient,
	}
}

// GetSummoner retrieves a summoner by name, using cache and database
func (s *Service) GetSummoner(ctx context.Context, region, summonerName string) (*riot.Summoner, error) {
	// Try API
	summoner, err := s.riotClient.GetSummonerByName(region, summonerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner from API: %w", err)
	}

	return summoner, nil
}

// GetMatches retrieves matches for a summoner by PUUID
func (s *Service) GetMatches(ctx context.Context, puuid string, count int) ([]string, error) {
	matchIDs, err := s.riotClient.GetMatchesByPUUID(puuid, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches from API: %w", err)
	}

	return matchIDs, nil
}

// GetMatch retrieves a match by ID, using cache and database
func (s *Service) GetMatch(ctx context.Context, matchID string) (*riot.Match, error) {

	// Get from API
	match, err := s.riotClient.GetMatchByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match from API: %w", err)
	}

	return match, nil
}
