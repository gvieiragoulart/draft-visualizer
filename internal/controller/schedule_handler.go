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

	enriched := r.URL.Query().Get("enriched") == "true"

	var response interface{}
	var err error

	if enriched {
		response, err = sh.service.GetScheduleEnriched(ctx)
	} else {
		response, err = sh.service.GetSchedule(ctx)
	}

	if err != nil {
		log.Printf("Error getting schedule: %v", err)
		http.Error(w, fmt.Sprintf("Error getting schedule: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (sh *ScheduleHandler) TeamsInScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	teamCodes, err := sh.service.GetTeamsInSchedule(ctx)
	if err != nil {
		log.Printf("Error getting teams in schedule: %v", err)
		http.Error(w, fmt.Sprintf("Error getting teams in schedule: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"teams": teamCodes,
		"count": len(teamCodes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (sh *ScheduleHandler) TeamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	teams, err := sh.service.GetTeamsData(ctx)
	if err != nil {
		log.Printf("Error getting teams: %v", err)
		http.Error(w, fmt.Sprintf("Error getting teams: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
