/*{
  "data": {
    "schedule": {
      "updated": "2025-10-05T16:27:59Z",
      "pages": {
        "older": "string",
        "newer": "string"
      },
      "events": [
        {
          "startTime": "2025-10-05T16:27:59Z",
          "blockName": "string",
          "match": {
            "teams": [
              {
                "code": "string",
                "image": "string",
                "name": "string",
                "result": {
                  "gameWins": 0,
                  "outcome": "loss"
                },
                "record": {
                  "losses": 0,
                  "wins": 0
                }
              },
              {
                "code": "string",
                "image": "string",
                "name": "string"
              }
            ],
            "id": "string",
            "strategy": {
              "count": 1,
              "type": "bestOf"
            }
          },
          "state": "completed",
          "type": "match",
          "league": {
            "name": "string",
            "slug": "string"
          }
        }
      ]
    }
  }
}*/

package dto

type ScheduleDTO struct {
	Data struct {
		Schedule struct {
			Updated string `json:"updated"`
			Pages   struct {
				Older string `json:"older"`
				Newer string `json:"newer"`
			} `json:"pages"`
			Events []struct {
				StartTime string `json:"startTime"`
				BlockName string `json:"blockName"`
				Match     struct {
					Teams []struct {
						Code   string `json:"code"`
						Image  string `json:"image"`
						Name   string `json:"name"`
						Result struct {
							GameWins int    `json:"gameWins"`
							Outcome  string `json:"outcome"`
						} `json:"result,omitempty"`
						Record struct {
							Losses int `json:"losses"`
							Wins   int `json:"wins"`
						} `json:"record,omitempty"`
					} `json:"teams"`
					ID       string `json:"id"`
					Strategy struct {
						Count int    `json:"count"`
						Type  string `json:"type"`
					} `json:"strategy"`
				} `json:"match"`
				State  string `json:"state"`
				Type   string `json:"type"`
				League struct {
					Name string `json:"name"`
					Slug string `json:"slug"`
				} `json:"league"`
			} `json:"events"`
		} `json:"schedule"`
	} `json:"data"`
}

func (s *ScheduleDTO) ToSchedule() *Schedule {
	events := make([]Event, len(s.Data.Schedule.Events))
	for i, e := range s.Data.Schedule.Events {
		teams := make([]Team, len(e.Match.Teams))
		for j, t := range e.Match.Teams {
			teams[j] = Team{
				Code:  t.Code,
				Image: t.Image,
				Name:  t.Name,
			}
			if t.Result != (struct {
				GameWins int    `json:"gameWins"`
				Outcome  string `json:"outcome"`
			}{}) {
				teams[j].Result = &Result{
					GameWins: t.Result.GameWins,
					Outcome:  t.Result.Outcome,
				}
			}
			if t.Record != (struct {
				Losses int `json:"losses"`
				Wins   int `json:"wins"`
			}{}) {
				teams[j].Record = &Record{
					Losses: t.Record.Losses,
					Wins:   t.Record.Wins,
				}
			}
		}
		events[i] = Event{
			StartTime: e.StartTime,
			BlockName: e.BlockName,
			Match: Match{
				Teams: teams,
				ID:    e.Match.ID,
				Strategy: Strategy{
					Count: e.Match.Strategy.Count,
					Type:  e.Match.Strategy.Type,
				},
			},
			State: e.State,
			Type:  e.Type,
			League: League{
				Name: e.League.Name,
				Slug: e.League.Slug,
			},
		}
	}
	return &Schedule{
		Updated: s.Data.Schedule.Updated,
		Pages: Pages{
			Older: s.Data.Schedule.Pages.Older,
			Newer: s.Data.Schedule.Pages.Newer,
		},
		Events: events,
	}
}

type Schedule struct {
	Updated string  `json:"updated"`
	Pages   Pages   `json:"pages"`
	Events  []Event `json:"events"`
}

type Pages struct {
	Older string `json:"older"`
	Newer string `json:"newer"`
}

type Event struct {
	StartTime string `json:"startTime"`
	BlockName string `json:"blockName"`
	Match     Match  `json:"match"`
	State     string `json:"state"`
	Type      string `json:"type"`
	League    League `json:"league"`
}

type Match struct {
	Teams    []Team   `json:"teams"`
	ID       string   `json:"id"`
	Strategy Strategy `json:"strategy"`
}

type Team struct {
	Code   string  `json:"code"`
	Image  string  `json:"image"`
	Name   string  `json:"name"`
	Result *Result `json:"result,omitempty"`
	Record *Record `json:"record,omitempty"`
}

type Result struct {
	GameWins int    `json:"gameWins"`
	Outcome  string `json:"outcome"`
}

type Record struct {
	Losses int `json:"losses"`
	Wins   int `json:"wins"`
}

type Strategy struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
}

type League struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
