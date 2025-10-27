package cargo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo/model/cargo_query"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/cargo/model/news_items"
)

type Client struct {
	clients.Client
}

type CargoResponse struct {
	CargoQuery []CargoItem `json:"cargoquery"`
}

type CargoItem struct {
	Title CargoItemTitle `json:"title"`
}

type CargoItemTitle struct {
	Tournaments          string `json:"Tournaments"`
	Teams                string `json:"Teams"`
	Players              string `json:"Players"`
	Region               string `json:"Region"`
	SubjectType          string `json:"SubjectType"`
	Subject              string `json:"Subject"`
	Sentence             string `json:"Sentence"`
	DateDisplay          string `json:"Date Display"`
	DateSort             string `json:"Date Sort"`
	IsApproxDate         string `json:"IsApproxDate"`
	DateSortPrecision    string `json:"Date Sort__precision"`
	DateDisplayAlias     string `json:"Date_Display,omitempty"`
	DateSortAlias        string `json:"Date_Sort,omitempty"`
	SentenceWithDate     string `json:"SentenceWithDate,omitempty"`
	SentenceTeam         string `json:"Sentence_Team,omitempty"`
	SentencePlayer       string `json:"Sentence_Player,omitempty"`
	SentenceTournament   string `json:"Sentence_Tournament,omitempty"`
	SubjectLink          string `json:"SubjectLink,omitempty"`
	Preload              string `json:"Preload,omitempty"`
	EarliestPossibleDate string `json:"EarliestPossibleDate,omitempty"`
	LatestPossibleDate   string `json:"LatestPossibleDate,omitempty"`
	Tags                 string `json:"Tags,omitempty"`
	Source               string `json:"Source,omitempty"`
	NLineInDate          string `json:"N_LineInDate,omitempty"`
	NewsId               string `json:"NewsId,omitempty"`
	ExcludeFrontpage     string `json:"ExcludeFrontpage,omitempty"`
	ExcludePortal        string `json:"ExcludePortal,omitempty"`
	ExcludeArchive       string `json:"ExcludeArchive,omitempty"`
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

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
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
		"IsApproxDate%20=%201",
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
	for _, item := range response.CargoQuery {
		newsItem := c.parseNewsItem(item.Title)
		newsItems = append(newsItems, newsItem)
	}
	return newsItems, nil
}

func (c *Client) parseNewsItem(title CargoItemTitle) news_items.NewsItems {
	// Get the field value, prioritizing the alias if available
	getField := func(primary, alias string) string {
		if alias != "" {
			return alias
		}
		return primary
	}

	// Get date display
	dateDisplay := getField(title.DateDisplay, title.DateDisplayAlias)

	// Get date sort
	dateSort := getField(title.DateSort, title.DateSortAlias)

	return news_items.NewsItems{
		DateDisplay:          dateDisplay,
		DateSort:             c.parseTimeString(dateSort),
		IsApproxDate:         c.parseBoolString(title.IsApproxDate),
		EarliestPossibleDate: c.parseTimeString(title.EarliestPossibleDate),
		LatestPossibleDate:   c.parseTimeString(title.LatestPossibleDate),
		Sentence:             title.Sentence,
		SentenceWithDate:     title.SentenceWithDate,
		SentenceTeam:         title.SentenceTeam,
		SentencePlayer:       title.SentencePlayer,
		SentenceTournament:   title.SentenceTournament,
		Subject:              title.Subject,
		SubjectType:          title.SubjectType,
		SubjectLink:          title.SubjectLink,
		Preload:              title.Preload,
		Region:               title.Region,
		Players:              parseStringSliceString(title.Players),
		Teams:                parseStringSliceString(title.Teams),
		Tournaments:          parseStringSliceString(title.Tournaments),
		Tags:                 parseStringSliceString(title.Tags),
		Source:               title.Source,
		NLineInDate:          parseStringInt(title.NLineInDate),
		NewsId:               title.NewsId,
		ExcludeFrontpage:     c.parseBoolString(title.ExcludeFrontpage),
		ExcludePortal:        c.parseBoolString(title.ExcludePortal),
		ExcludeArchive:       c.parseBoolString(title.ExcludeArchive),
	}
}

// Helper functions for parsing string values
func (c *Client) parseTimeString(v string) *time.Time {
	if v == "" {
		return nil
	}
	// Try parsing with the format used in the response: "2006-01-02 15:04:05"
	t, err := time.Parse("2006-01-02 15:04:05", v)
	if err != nil {
		return nil
	}
	return &t
}

func (c *Client) parseBoolString(v string) *bool {
	if v == "1" || v == "true" {
		b := true
		return &b
	}
	b := false
	return &b
}

func parseStringSliceString(v string) []string {
	if v == "" {
		return []string{}
	}
	// Simple split by comma for now
	items := strings.Split(v, ",")
	var result []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseStringInt(v string) *int {
	if v == "" {
		return nil
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &num
}
