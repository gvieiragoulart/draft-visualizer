package esports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/esports/dto"
)

type EsportsClient struct {
	clients.Client
}

func NewClient(apiKey string) *EsportsClient {
	return &EsportsClient{
		Client: clients.Client{
			ApiKey:     apiKey,
			HttpClient: &http.Client{},
			BaseURL:    "https://esports-api.lolesports.com/persisted/gw",
		},
	}
}

func (e *EsportsClient) GetSchedule() (dto.ScheduleDTO, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/getSchedule?hl=pt-BR", e.BaseURL), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("x-api-key", e.ApiKey)

	resp, err := e.HttpClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	var scheduleResponse dto.ScheduleDTO
	if err := json.NewDecoder(resp.Body).Decode(&scheduleResponse); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	return scheduleResponse, nil
}
