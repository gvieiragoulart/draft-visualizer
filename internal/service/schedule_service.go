package service

import (
	"context"
	"fmt"

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

// GetSchedule returns the basic schedule data
func (s *ScheduleService) GetSchedule(ctx context.Context) (dto.ScheduleDTO, error) {
	schedule, err := s.scheduleClient.GetSchedule()
	if err != nil {
		return dto.ScheduleDTO{}, fmt.Errorf("error getting schedule: %w", err)
	}

	return schedule, nil
}

func (s *ScheduleService) GetScheduleEnriched(ctx context.Context) (*dto.ScheduleEnriched, error) {
	schedule, err := s.scheduleClient.GetSchedule()
	if err != nil {
		return nil, fmt.Errorf("error getting schedule: %w", err)
	}

	teams, err := s.scheduleClient.GetTeams()
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	// Merge schedule with team data
	enrichedSchedule := schedule.MergeWithTeams(teams)

	return enrichedSchedule, nil
}

// GetTeamsInSchedule returns all unique teams that appear in the current schedule
func (s *ScheduleService) GetTeamsInSchedule(ctx context.Context) ([]string, error) {
	schedule, err := s.scheduleClient.GetSchedule()
	if err != nil {
		return nil, fmt.Errorf("error getting schedule: %w", err)
	}

	teamCodes := schedule.GetUniqueTeamCodes()
	return teamCodes, nil
}

// GetTeamsData returns the full teams data
func (s *ScheduleService) GetTeamsData(ctx context.Context) (dto.TeamsDTO, error) {
	teams, err := s.scheduleClient.GetTeams()
	if err != nil {
		return dto.TeamsDTO{}, fmt.Errorf("error getting teams: %w", err)
	}

	return teams, nil
}
