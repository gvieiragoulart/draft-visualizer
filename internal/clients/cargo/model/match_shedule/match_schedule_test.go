package match_schedule

import (
	"testing"
	"time"
)

func TestMatchSchedule_Fields(t *testing.T) {
	tests := []struct {
		name  string
		ms    MatchSchedule
		check func(*testing.T, MatchSchedule)
	}{
		{
			name: "complete match schedule",
			ms: MatchSchedule{
				Team1:        "Team A",
				Team2:        "Team B",
				Team1Final:   "Team A",
				Team2Final:   "Team B",
				Winner:       "Team A",
				Team1Poster:  "poster1.jpg",
				Team2Poster:  "poster2.jpg",
				Team1Points:  intPtr(2),
				Team2Points:  intPtr(0),
				Team1Score:   intPtr(15),
				Team2Score:   intPtr(5),
				BestOf:       intPtr(3),
				OverviewPage: "2024_Summer",
				Round:        "Quarterfinals",
				Phase:        "Playoffs",
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.Team1 != "Team A" {
					t.Errorf("expected Team1 to be 'Team A', got %s", ms.Team1)
				}
				if ms.Winner != "Team A" {
					t.Errorf("expected Winner to be 'Team A', got %s", ms.Winner)
				}
				if ms.Team1Points == nil || *ms.Team1Points != 2 {
					t.Errorf("expected Team1Points to be 2, got %v", ms.Team1Points)
				}
			},
		},
		{
			name: "match with null score",
			ms: MatchSchedule{
				Team1:       "Team X",
				Team2:       "Team Y",
				IsNullified: boolPtr(true),
				FF:          intPtr(1),
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.IsNullified == nil || !*ms.IsNullified {
					t.Errorf("expected IsNullified to be true, got %v", ms.IsNullified)
				}
				if ms.FF == nil || *ms.FF != 1 {
					t.Errorf("expected FF to be 1, got %v", ms.FF)
				}
			},
		},
		{
			name: "match with date and time",
			ms: MatchSchedule{
				Team1:           "Team C",
				Team2:           "Team D",
				DateTimeUTC:     timePtr(time.Date(2024, 7, 15, 14, 30, 0, 0, time.UTC)),
				HasTime:         boolPtr(true),
				DST:             "UTC",
				MatchDay:        intPtr(15),
				IsReschedulable: boolPtr(true),
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.DateTimeUTC == nil {
					t.Errorf("expected DateTimeUTC to be set, got nil")
				}
				if ms.HasTime == nil || !*ms.HasTime {
					t.Errorf("expected HasTime to be true, got %v", ms.HasTime)
				}
				if ms.MatchDay == nil || *ms.MatchDay != 15 {
					t.Errorf("expected MatchDay to be 15, got %v", ms.MatchDay)
				}
			},
		},
		{
			name: "match with disabled champions",
			ms: MatchSchedule{
				Team1:             "Team E",
				Team2:             "Team F",
				Patch:             "14.12",
				LegacyPatch:       "14.12.1",
				DisabledChampions: []string{"Champion1", "Champion2", "Champion3"},
				BestOf:            intPtr(5),
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if len(ms.DisabledChampions) != 3 {
					t.Errorf("expected 3 disabled champions, got %d", len(ms.DisabledChampions))
				}
				if ms.Patch != "14.12" {
					t.Errorf("expected Patch to be '14.12', got %s", ms.Patch)
				}
				if ms.BestOf == nil || *ms.BestOf != 5 {
					t.Errorf("expected BestOf to be 5, got %v", ms.BestOf)
				}
			},
		},
		{
			name: "match with casters and stream",
			ms: MatchSchedule{
				Team1:         "Team G",
				Team2:         "Team H",
				Stream:        "Twitch Stream",
				StreamDisplay: "twitch.tv/lol",
				Venue:         "Arena",
				Casters:       []string{"Caster1", "Caster2"},
				CastersPBP:    "Play-by-Play",
				CastersColor:  "Blue",
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.StreamDisplay != "twitch.tv/lol" {
					t.Errorf("expected StreamDisplay to be 'twitch.tv/lol', got %s", ms.StreamDisplay)
				}
				if len(ms.Casters) != 2 {
					t.Errorf("expected 2 casters, got %d", len(ms.Casters))
				}
				if ms.Venue != "Arena" {
					t.Errorf("expected Venue to be 'Arena', got %s", ms.Venue)
				}
			},
		},
		{
			name: "match with MVP and VOD",
			ms: MatchSchedule{
				Team1:         "Team I",
				Team2:         "Team J",
				MVP:           "Player1",
				MVPPoints:     intPtr(120),
				VodInterview:  "Interview Link",
				VodHighlights: "Highlights Link",
				Recap:         "Match recap text",
				Reddit:        "Reddit thread",
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.MVP != "Player1" {
					t.Errorf("expected MVP to be 'Player1', got %s", ms.MVP)
				}
				if ms.MVPPoints == nil || *ms.MVPPoints != 120 {
					t.Errorf("expected MVPPoints to be 120, got %v", ms.MVPPoints)
				}
				if ms.Recap == "" {
					t.Errorf("expected Recap to be set, got empty string")
				}
			},
		},
		{
			name: "match with page navigation",
			ms: MatchSchedule{
				Team1:               "Team K",
				Team2:               "Team L",
				N_MatchInPage:       intPtr(3),
				Tab:                 "Main",
				N_MatchInTab:        intPtr(2),
				N_TabInPage:         intPtr(1),
				N_Page:              intPtr(5),
				InitialN_MatchInTab: intPtr(2),
				InitialPageAndTab:   "5-1",
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.N_Page == nil || *ms.N_Page != 5 {
					t.Errorf("expected N_Page to be 5, got %v", ms.N_Page)
				}
				if ms.Tab != "Main" {
					t.Errorf("expected Tab to be 'Main', got %s", ms.Tab)
				}
				if ms.N_MatchInPage == nil || *ms.N_MatchInPage != 3 {
					t.Errorf("expected N_MatchInPage to be 3, got %v", ms.N_MatchInPage)
				}
			},
		},
		{
			name: "match with tags and unique identifier",
			ms: MatchSchedule{
				Team1:       "Team M",
				Team2:       "Team N",
				UniqueMatch: "2024_Summer_Playoffs_QF1",
				MatchId:     "match-001",
				Tags:        []string{"live", "featured", "playoffs"},
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.UniqueMatch != "2024_Summer_Playoffs_QF1" {
					t.Errorf("expected UniqueMatch to be '2024_Summer_Playoffs_QF1', got %s", ms.UniqueMatch)
				}
				if ms.MatchId != "match-001" {
					t.Errorf("expected MatchId to be 'match-001', got %s", ms.MatchId)
				}
				if len(ms.Tags) != 3 {
					t.Errorf("expected 3 tags, got %d", len(ms.Tags))
				}
			},
		},
		{
			name: "match with prediction settings",
			ms: MatchSchedule{
				Team1:                       "Team O",
				Team2:                       "Team P",
				OverrideAllowPredictions:    boolPtr(true),
				OverrideDisallowPredictions: boolPtr(false),
				IsTiebreaker:                boolPtr(false),
				IsFlexibleStart:             boolPtr(true),
			},
			check: func(t *testing.T, ms MatchSchedule) {
				if ms.OverrideAllowPredictions == nil || !*ms.OverrideAllowPredictions {
					t.Errorf("expected OverrideAllowPredictions to be true, got %v", ms.OverrideAllowPredictions)
				}
				if ms.IsTiebreaker == nil || *ms.IsTiebreaker {
					t.Errorf("expected IsTiebreaker to be false, got %v", ms.IsTiebreaker)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.ms)
		})
	}
}

func TestGetFields(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
		checkFields   func(*testing.T, []string)
	}{
		{
			name:          "returns all field names",
			expectedCount: 81,
			checkFields: func(t *testing.T, fields []string) {
				if len(fields) != 81 {
					t.Errorf("expected 81 fields, got %d", len(fields))
				}
				// Check for some expected field names
				expectedFields := []string{
					"Team1", "Team2", "Winner", "DateTime_UTC",
					"BestOf", "MatchId", "UniqueMatch", "Casters",
				}
				fieldMap := make(map[string]bool)
				for _, field := range fields {
					fieldMap[field] = true
				}
				for _, expected := range expectedFields {
					if !fieldMap[expected] {
						t.Errorf("expected field '%s' not found", expected)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := GetFields()
			if len(fields) != tt.expectedCount {
				t.Errorf("expected %d fields, got %d", tt.expectedCount, len(fields))
			}
			if tt.checkFields != nil {
				tt.checkFields(t, fields)
			}
		})
	}
}

// Helper functions for creating pointers
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func timePtr(t time.Time) *time.Time {
	return &t
}
