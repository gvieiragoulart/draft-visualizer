package service

import (
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo/model/news_items"
)

type CargoService struct {
	cargoClient *cargo.Client
}

func NewCargoService(cargoClient *cargo.Client) *CargoService {
	return &CargoService{cargoClient: cargoClient}
}

func (s *CargoService) GetNewsItems() ([]news_items.NewsItems, error) {
	return s.cargoClient.GetNewsLatest()
}
