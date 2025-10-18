package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gvieiragoulart/draft-visualizer/internal/service"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

func (sh *ScheduleHandler) ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	schedule, err := sh.service.GetSchedule(ctx)
	if err != nil {
		log.Printf("Error getting schedule: %v", err)
		http.Error(w, fmt.Sprintf("Error getting schedule: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
