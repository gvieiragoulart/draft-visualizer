package controller

import (
	"fmt"

	"github.com/gvieiragoulart/draft-visualizer/internal/config"
	"github.com/gvieiragoulart/draft-visualizer/internal/service"
)

type WikiHandler struct {
	wikiService *service.WikiService
}

func NewWikiHandler(cfg *config.Config, wikiService *service.WikiService) (*WikiHandler, error) {
	if wikiService == nil {
		var err error
		wikiService, err = service.NewWikiService(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating wiki service: %w", err)
		}
	}
	return &WikiHandler{wikiService: wikiService}, nil
}
