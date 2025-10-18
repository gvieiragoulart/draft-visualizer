package service

import (
	"context"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients/esports"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/esports/dto"
)

type ScheduleService struct {
	scheduleClient *esports.EsportsClient
}

func NewScheduleService(scheduleClient *esports.EsportsClient) *ScheduleService {
	return &ScheduleService{
		scheduleClient: scheduleClient,
	}
}

func (s *ScheduleService) GetSchedule(ctx context.Context) (dto.ScheduleDTO, error) {
	return s.scheduleClient.GetSchedule()
}
