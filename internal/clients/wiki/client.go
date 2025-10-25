package wiki

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	mwclient "cgt.name/pkg/go-mwclient"
)

type WikiClient struct {
	client *mwclient.Client
}

func NewClient(username, password string) (*WikiClient, error) {
	w, err := mwclient.New("https://lol.fandom.com/api.php", "LolWiki")
	if err != nil {
		return nil, fmt.Errorf("error creating wiki client: %w", err)
	}

	err = w.Login(username, password)
	if err != nil {
		return nil, fmt.Errorf("error logging into wiki: %w", err)
	}

	log.Printf("Successfully logged into wiki as user: %s", username)

	return &WikiClient{client: w}, nil
}

func NewClientWithoutLogin() (*WikiClient, error) {
	w, err := mwclient.New("https://lol.fandom.com/api.php", "LolWiki")
	if err != nil {
		return nil, fmt.Errorf("error creating wiki client: %w", err)
	}

	return &WikiClient{client: w}, nil
}

func (c *WikiClient) GetWikiPage(ctx context.Context, page string) (string, error) {
	result, _, err := c.client.GetPageByName(page)
	if err != nil {
		return "", fmt.Errorf("error getting page: %w", err)
	}
	return result, nil
}

func (c *WikiClient) GetPageContent(ctx context.Context, page string) (string, error) {
	// Get page content using the API
	result, _, err := c.client.GetPageByName(page)
	if err != nil {
		return "", fmt.Errorf("error getting page content: %w", err)
	}
	return result, nil
}

func (c *WikiClient) GetPageInfo(ctx context.Context, page string) (interface{}, error) {
	// Get page info including infobox data
	params := map[string]string{
		"action": "query",
		"format": "json",
		"titles": page,
		"prop":   "revisions|pageprops",
		"rvprop": "content",
		"ppprop": "infoboxes",
	}

	result, err := c.client.Get(params)
	if err != nil {
		return nil, fmt.Errorf("error getting page info: %w", err)
	}

	return result, nil
}

func (c *WikiClient) SearchPages(ctx context.Context, query string) (interface{}, error) {
	// Search for pages
	params := map[string]string{
		"action":   "query",
		"format":   "json",
		"list":     "search",
		"srsearch": query,
		"srlimit":  "10",
	}

	result, err := c.client.Get(params)
	if err != nil {
		return nil, fmt.Errorf("error searching pages: %w", err)
	}

	return result, nil
}

// PlayerInfo represents information about a League of Legends player
type PlayerInfo struct {
	Name      string `json:"name"`
	Position  string `json:"position"`
	Country   string `json:"country"`
	JoinDate  string `json:"join_date"`
	LeaveDate string `json:"leave_date,omitempty"`
	Status    string `json:"status"` // "Active" or "Inactive"
}

// TeamRoster represents the roster of a team
type TeamRoster struct {
	TeamName string       `json:"team_name"`
	Players  []PlayerInfo `json:"players"`
}

func (c *WikiClient) GetTeamRoster(ctx context.Context, teamPage string) (*TeamRoster, error) {
	// Get the page content
	content, err := c.GetPageContent(ctx, teamPage)
	if err != nil {
		return nil, fmt.Errorf("error getting team page content: %w", err)
	}

	roster := &TeamRoster{
		TeamName: teamPage,
		Players:  []PlayerInfo{},
	}

	// Parse the content to extract player information
	players := c.parseTeamRoster(content)
	roster.Players = players

	return roster, nil
}

func (c *WikiClient) parseTeamRoster(content string) []PlayerInfo {
	var players []PlayerInfo

	// Regular expression to match player information in infobox format
	playerRegex := regexp.MustCompile(`(?i)\|\s*player(\d+)\s*=\s*([^|\n]+)`)

	// Find all players
	playerMatches := playerRegex.FindAllStringSubmatch(content, -1)

	for _, match := range playerMatches {
		if len(match) >= 3 {
			playerNum := match[1]
			playerName := strings.TrimSpace(match[2])

			// Skip empty or placeholder names
			if playerName == "" || playerName == "TBD" || playerName == "TBA" {
				continue
			}

			player := PlayerInfo{
				Name:   playerName,
				Status: "Active", // Default to active
			}

			// Find position for this player
			posPattern := fmt.Sprintf(`(?i)\|\s*position%s\s*=\s*([^|\n]+)`, playerNum)
			posRegex := regexp.MustCompile(posPattern)
			if posMatch := posRegex.FindStringSubmatch(content); len(posMatch) >= 2 {
				player.Position = strings.TrimSpace(posMatch[1])
			}

			// Find country for this player
			countryPattern := fmt.Sprintf(`(?i)\|\s*country%s\s*=\s*([^|\n]+)`, playerNum)
			countryRegex := regexp.MustCompile(countryPattern)
			if countryMatch := countryRegex.FindStringSubmatch(content); len(countryMatch) >= 2 {
				player.Country = strings.TrimSpace(countryMatch[1])
			}

			// Find join date
			joinPattern := fmt.Sprintf(`(?i)\|\s*join%s\s*=\s*([^|\n]+)`, playerNum)
			joinRegex := regexp.MustCompile(joinPattern)
			if joinMatch := joinRegex.FindStringSubmatch(content); len(joinMatch) >= 2 {
				player.JoinDate = strings.TrimSpace(joinMatch[1])
			}

			// Find leave date (if exists, player is inactive)
			leavePattern := fmt.Sprintf(`(?i)\|\s*leave%s\s*=\s*([^|\n]+)`, playerNum)
			leaveRegex := regexp.MustCompile(leavePattern)
			if leaveMatch := leaveRegex.FindStringSubmatch(content); len(leaveMatch) >= 2 {
				player.LeaveDate = strings.TrimSpace(leaveMatch[1])
				player.Status = "Inactive"
			}

			players = append(players, player)
		}
	}

	return players
}
