/*
{
  "data": {
    "teams": [
      {
        "code": "string",
        "image": "string",
        "name": "string",
        "id": "string",
        "slug": "string",
        "alternativeImage": "string",
        "homeLeague": {
          "name": "string",
          "region": "string"
        },
        "players": [
          {
            "id": "string",
            "summonerName": "string",
            "firstName": "string",
            "lastName": "string",
            "image": "string",
            "role": "string"
          }
        ]
      }
    ]
  }
}
*/

package dto

type TeamsDTO struct {
	Data struct {
		Teams []struct {
			Code             string `json:"code"`
			Image            string `json:"image"`
			Name             string `json:"name"`
			ID               string `json:"id"`
			Slug             string `json:"slug"`
			AlternativeImage string `json:"alternativeImage"`
			HomeLeague       struct {
				Name   string `json:"name"`
				Region string `json:"region"`
			} `json:"homeLeague"`
			Players []struct {
				ID           string `json:"id"`
				SummonerName string `json:"summonerName"`
				FirstName    string `json:"firstName"`
				LastName     string `json:"lastName"`
				Image        string `json:"image"`
				Role         string `json:"role"`
			} `json:"players"`
		} `json:"teams"`
	} `json:"data"`
}

type Teams struct {
	Code             string `json:"code"`
	Image            string `json:"image"`
	Name             string `json:"name"`
	ID               string `json:"id"`
	Slug             string `json:"slug"`
	AlternativeImage string `json:"alternativeImage"`
	HomeLeague       struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	} `json:"homeLeague"`
	Players []struct {
		ID           string `json:"id"`
		SummonerName string `json:"summonerName"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Image        string `json:"image"`
		Role         string `json:"role"`
	} `json:"players"`
}

func (t *TeamsDTO) ToTeams() []Teams {
	teams := make([]Teams, len(t.Data.Teams))
	for i, team := range t.Data.Teams {
		teams[i] = Teams{
			Code:             team.Code,
			Image:            team.Image,
			Name:             team.Name,
			ID:               team.ID,
			Slug:             team.Slug,
			AlternativeImage: team.AlternativeImage,
			HomeLeague:       team.HomeLeague,
			Players:          team.Players,
		}
	}
	return teams
}
