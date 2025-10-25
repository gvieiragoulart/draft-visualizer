package service

import (
	"fmt"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients/wiki"
	"github.com/gvieiragoulart/draft-visualizer/internal/config"
)

type WikiService struct {
	wikiClient *wiki.WikiClient
}

func NewWikiService(cfg *config.Config) (*WikiService, error) {
	wikiClient, err := wiki.NewClient(cfg.WikiUsername, cfg.WikiPassword)
	if err != nil {
		return nil, fmt.Errorf("error creating wiki service: %w", err)
	}

	return &WikiService{wikiClient: wikiClient}, nil
}

func NewWikiServiceWithClient(wikiClient *wiki.WikiClient) *WikiService {
	return &WikiService{wikiClient: wikiClient}
}
