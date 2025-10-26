package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gvieiragoulart/draft-visualizer/internal/service"
)

type CargoHandler interface {
	GetNewsLatest(w http.ResponseWriter, r *http.Request)
}

type CargoHandlerImpl struct {
	service *service.CargoService
}

func NewCargoHandler(service *service.CargoService) CargoHandler {
	return &CargoHandlerImpl{service: service}
}

func (h *CargoHandlerImpl) GetNewsLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newsItems, err := h.service.GetNewsItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsItems)
}
