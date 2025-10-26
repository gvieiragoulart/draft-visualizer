package match_schedule

import (
	"time"
)

type MatchSchedule struct {
	Team1       string `json:"Team1"`
	Team2       string `json:"Team2"`
	Team1Final  string `json:"Team1Final"`
	Team2Final  string `json:"Team2Final"`
	Winner      string `json:"Winner"`
	Team1Poster string `json:"Team1Poster"`
	Team2Poster string `json:"Team2Poster"`

	Team1Points    *int `json:"Team1Points"`
	Team2Points    *int `json:"Team2Points"`
	Team1PointsTB  *int `json:"Team1PointsTB"`
	Team2PointsTB  *int `json:"Team2PointsTB"`
	Team1Score     *int `json:"Team1Score"`
	Team2Score     *int `json:"Team2Score"`
	Team1Advantage *int `json:"Team1Advantage"`
	Team2Advantage *int `json:"Team2Advantage"`

	FF          *int  `json:"FF"`
	IsNullified *bool `json:"IsNullified"`

	Player1 string `json:"Player1"`
	Player2 string `json:"Player2"`

	MatchDay        *int       `json:"MatchDay"`
	DateTimeUTC     *time.Time `json:"DateTime_UTC"`
	HasTime         *bool      `json:"HasTime"`
	DST             string     `json:"DST"`
	IsFlexibleStart *bool      `json:"IsFlexibleStart"`
	IsReschedulable *bool      `json:"IsReschedulable"`

	OverrideAllowPredictions    *bool `json:"OverrideAllowPredictions"`
	OverrideDisallowPredictions *bool `json:"OverrideDisallowPredictions"`
	IsTiebreaker                *bool `json:"IsTiebreaker"`
	BestOf                      *int  `json:"BestOf"`

	OverviewPage string `json:"OverviewPage"`
	ShownName    string `json:"ShownName"`
	ShownRound   string `json:"ShownRound"`
	Round        string `json:"Round"`
	Phase        string `json:"Phase"`
	GroupName    string `json:"GroupName"`

	N_MatchInPage       *int   `json:"N_MatchInPage"`
	Tab                 string `json:"Tab"`
	N_MatchInTab        *int   `json:"N_MatchInTab"`
	N_TabInPage         *int   `json:"N_TabInPage"`
	N_Page              *int   `json:"N_Page"`
	InitialN_MatchInTab *int   `json:"InitialN_MatchInTab"`
	InitialPageAndTab   string `json:"InitialPageAndTab"`

	Patch             string   `json:"Patch"`
	LegacyPatch       string   `json:"LegacyPatch"`
	PatchPage         string   `json:"PatchPage"`
	Hotfix            string   `json:"Hotfix"`
	DisabledChampions []string `json:"DisabledChampions"`
	PatchFootnote     string   `json:"PatchFootnote"`

	Stream        string   `json:"Stream"`
	StreamDisplay string   `json:"StreamDisplay"`
	Venue         string   `json:"Venue"`
	CastersPBP    string   `json:"CastersPBP"`
	CastersColor  string   `json:"CastersColor"`
	Casters       []string `json:"Casters"`

	MVP       string `json:"MVP"`
	MVPPoints *int   `json:"MVPPoints"`

	VodInterview  string   `json:"VodInterview"`
	VodHighlights string   `json:"VodHighlights"`
	InterviewWith []string `json:"InterviewWith"`
	Recap         string   `json:"Recap"`
	Reddit        string   `json:"Reddit"`

	QQ            *int   `json:"QQ"`
	Wanplus       string `json:"Wanplus"`
	WanplusId     *int   `json:"WanplusId"`
	PageAndTeam1  string `json:"PageAndTeam1"`
	PageAndTeam2  string `json:"PageAndTeam2"`
	Team1Footnote string `json:"Team1Footnote"`
	Team2Footnote string `json:"Team2Footnote"`
	Footnote      string `json:"Footnote"`

	UniqueMatch string   `json:"UniqueMatch"`
	MatchId     string   `json:"MatchId"`
	Tags        []string `json:"Tags"`
}

// GetFields returns all field names for MatchSchedule
func GetFields() []string {
	return []string{
		"Team1", "Team2", "Team1Final", "Team2Final", "Winner",
		"Team1Points", "Team2Points", "Team1PointsTB", "Team2PointsTB",
		"Team1Score", "Team2Score", "Team1Poster", "Team2Poster",
		"Team1Advantage", "Team2Advantage", "FF", "IsNullified",
		"Player1", "Player2", "MatchDay", "DateTime_UTC", "HasTime",
		"DST", "IsFlexibleStart", "IsReschedulable",
		"OverrideAllowPredictions", "OverrideDisallowPredictions",
		"IsTiebreaker", "OverviewPage", "ShownName", "ShownRound",
		"BestOf", "Round", "Phase", "N_MatchInPage", "Tab",
		"N_MatchInTab", "N_TabInPage", "N_Page",
		"InitialN_MatchInTab", "InitialPageAndTab",
		"Patch", "LegacyPatch", "PatchPage", "Hotfix",
		"DisabledChampions", "PatchFootnote", "GroupName",
		"Stream", "StreamDisplay", "Venue", "CastersPBP",
		"CastersColor", "Casters", "MVP", "MVPPoints",
		"VodInterview", "VodHighlights", "InterviewWith",
		"Recap", "Reddit", "QQ", "Wanplus", "WanplusId",
		"PageAndTeam1", "PageAndTeam2", "Team1Footnote",
		"Team2Footnote", "Footnote", "UniqueMatch", "MatchId", "Tags",
	}
}
