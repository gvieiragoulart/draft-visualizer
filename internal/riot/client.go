package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the Riot Games API client
type Client struct {
	apiKey     string
	httpClient HTTPClient
	baseURL    string
}

// NewClient creates a new Riot Games API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://americas.api.riotgames.com",
	}
}

// NewClientWithHTTPClient creates a new client with a custom HTTP client
func NewClientWithHTTPClient(apiKey string, httpClient HTTPClient) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    "https://americas.api.riotgames.com",
	}
}

// SetBaseURL sets the base URL for the client (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// Summoner represents a summoner from the Riot API
type Summoner struct {
	PUUID         string `json:"puuid"`
	SummonerName  string `json:"name"`
	SummonerLevel int    `json:"summonerLevel"`
	ProfileIconID int    `json:"profileIconId"`
}

// Match represents a match from the Riot API
type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

// MatchMetadata contains match metadata
type MatchMetadata struct {
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

// MatchInfo contains match information
type MatchInfo struct {
	GameMode     string `json:"gameMode"`
	GameDuration int    `json:"gameDuration"`
	GameCreation int64  `json:"gameCreation"`
}

// GetSummonerByName retrieves a summoner by name
func (c *Client) GetSummonerByName(region, summonerName string) (*Summoner, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", region, summonerName)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("X-Riot-Token", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var summoner Summoner
	if err := json.NewDecoder(resp.Body).Decode(&summoner); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &summoner, nil
}

// GetMatchesByPUUID retrieves match IDs for a player by PUUID
func (c *Client) GetMatchesByPUUID(puuid string, count int) ([]string, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?count=%d", c.baseURL, puuid, count)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("X-Riot-Token", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var matchIDs []string
	if err := json.NewDecoder(resp.Body).Decode(&matchIDs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return matchIDs, nil
}

// GetMatchByID retrieves a match by ID
func (c *Client) GetMatchByID(matchID string) (*Match, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/%s", c.baseURL, matchID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("X-Riot-Token", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var match Match
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &match, nil
}
