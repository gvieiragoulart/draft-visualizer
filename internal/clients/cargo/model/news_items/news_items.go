package news_items

import "time"

type NewsItems struct {
	DateDisplay  string     `json:"Date_Display"`
	DateSort     *time.Time `json:"Date_Sort"`
	IsApproxDate *bool      `json:"IsApproxDate"`

	EarliestPossibleDate *time.Time `json:"EarliestPossibleDate"`
	LatestPossibleDate   *time.Time `json:"LatestPossibleDate"`

	Sentence           string `json:"Sentence"`
	SentenceWithDate   string `json:"SentenceWithDate"`
	SentenceTeam       string `json:"Sentence_Team"`
	SentencePlayer     string `json:"Sentence_Player"`
	SentenceTournament string `json:"Sentence_Tournament"`

	Subject     string `json:"Subject"`
	SubjectType string `json:"SubjectType"`
	SubjectLink string `json:"SubjectLink"`
	Preload     string `json:"Preload"`
	Region      string `json:"Region"`

	Players     []string `json:"Players"`
	Teams       []string `json:"Teams"`
	Tournaments []string `json:"Tournaments"`
	Tags        []string `json:"Tags"`

	Source           string `json:"Source"`
	NLineInDate      *int   `json:"N_LineInDate"`
	NewsId           string `json:"NewsId"`
	ExcludeFrontpage *bool  `json:"ExcludeFrontpage"`
	ExcludePortal    *bool  `json:"ExcludePortal"`
	ExcludeArchive   *bool  `json:"ExcludeArchive"`
}

func GetFields() []string {
	return []string{
		"Date_Display",
		"Date_Sort",
		"IsApproxDate",
		"EarliestPossibleDate",
		"LatestPossibleDate",
		"Sentence",
		"SentenceWithDate",
		"Sentence_Team",
		"Sentence_Player",
		"Sentence_Tournament",
		"Subject",
		"SubjectType",
		"SubjectLink",
		"Preload",
		"Region",
		"Players",
		"Teams",
		"Tournaments",
		"Tags",
		"Source",
		"N_LineInDate",
		"NewsId",
		"ExcludeFrontpage",
		"ExcludePortal",
		"ExcludeArchive",
	}
}
