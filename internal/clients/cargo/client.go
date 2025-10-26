package cargo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo/model/cargo_query"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo/model/news_items"
)

type Client struct {
	clients.Client
}

type CargoResponse struct {
	CargoQuery struct {
		Count  int         `json:"count"`
		Format string      `json:"format"`
		Items  []CargoItem `json:"items"`
	} `json:"cargoquery"`
}

type CargoItem struct {
	Title map[string]interface{} `json:"title"`
}

func NewClient() *Client {
	return &Client{
		Client: clients.Client{
			HttpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
			BaseURL: "https://lol.fandom.com/api.php",
		},
	}
}

func NewClientWithHTTPClient(httpClient clients.HTTPClient) *Client {
	return &Client{
		Client: clients.Client{
			HttpClient: httpClient,
			BaseURL:    "https://lol.fandom.com/api.php",
		},
	}
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) Query(query *cargo_query.CargoQuery) (*CargoResponse, error) {
	queryString := query.ToQuery()
	fullURL := fmt.Sprintf("%s?format=json&%s", c.BaseURL, queryString)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "DraftVisualizer/1.0")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var cargoResponse CargoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cargoResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &cargoResponse, nil
}

func (c *Client) GetNewsLatest() ([]news_items.NewsItems, error) {
	query := cargo_query.NewCargoQuery(
		[]string{"NewsItems"},
		news_items.GetFields(),
		"IsApproxDate = 1",
		"",
		"",
		"",
		"",
		0,
		500,
	)
	response, err := c.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error querying news items: %w", err)
	}

	var newsItems []news_items.NewsItems
	for _, item := range response.CargoQuery.Items {
		newsItem := c.parseNewsItem(item.Title)
		newsItems = append(newsItems, newsItem)
	}
	return newsItems, nil
}

func (c *Client) parseNewsItem(title map[string]interface{}) news_items.NewsItems {
	return news_items.NewsItems{
		DateDisplay:          parseString(title["Date_Display"]),
		DateSort:             parseTime(title["Date_Sort"]),
		IsApproxDate:         parseBool(title["IsApproxDate"]),
		EarliestPossibleDate: parseTime(title["EarliestPossibleDate"]),
		LatestPossibleDate:   parseTime(title["LatestPossibleDate"]),
		Sentence:             parseString(title["Sentence"]),
		SentenceWithDate:     parseString(title["SentenceWithDate"]),
		SentenceTeam:         parseString(title["Sentence_Team"]),
		SentencePlayer:       parseString(title["Sentence_Player"]),
		SentenceTournament:   parseString(title["Sentence_Tournament"]),
		Subject:              parseString(title["Subject"]),
		SubjectType:          parseString(title["SubjectType"]),
		SubjectLink:          parseString(title["SubjectLink"]),
		Preload:              parseString(title["Preload"]),
		Region:               parseString(title["Region"]),
		Players:              parseStringSlice(title["Players"]),
		Teams:                parseStringSlice(title["Teams"]),
		Tournaments:          parseStringSlice(title["Tournaments"]),
		Tags:                 parseStringSlice(title["Tags"]),
		Source:               parseString(title["Source"]),
		NLineInDate:          parseInt(title["N_LineInDate"]),
		NewsId:               parseString(title["NewsId"]),
		ExcludeFrontpage:     parseBool(title["ExcludeFrontpage"]),
		ExcludePortal:        parseBool(title["ExcludePortal"]),
		ExcludeArchive:       parseBool(title["ExcludeArchive"]),
	}
}
