package riot

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	if client == nil {
		t.Fatal("expected client to be created")
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey to be 'test-api-key', got %s", client.apiKey)
	}
	if client.baseURL != "https://americas.api.riotgames.com" {
		t.Errorf("expected baseURL to be 'https://americas.api.riotgames.com', got %s", client.baseURL)
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	mockClient := &MockHTTPClient{}
	client := NewClientWithHTTPClient("test-api-key", mockClient)
	if client == nil {
		t.Fatal("expected client to be created")
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey to be 'test-api-key', got %s", client.apiKey)
	}
	if client.httpClient != mockClient {
		t.Error("expected httpClient to be the mock client")
	}
}

func TestSetBaseURL(t *testing.T) {
	client := NewClient("test-api-key")
	client.SetBaseURL("https://test.example.com")
	if client.baseURL != "https://test.example.com" {
		t.Errorf("expected baseURL to be 'https://test.example.com', got %s", client.baseURL)
	}
}

func TestGetSummonerByName_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.Header.Get("X-Riot-Token") != "test-api-key" {
				t.Errorf("expected X-Riot-Token header to be 'test-api-key', got %s", req.Header.Get("X-Riot-Token"))
			}
			
			// Return mock response
			responseBody := `{
				"puuid": "test-puuid",
				"name": "TestSummoner",
				"summonerLevel": 100,
				"profileIconId": 1234
			}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	summoner, err := client.GetSummonerByName("na1", "TestSummoner")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if summoner.PUUID != "test-puuid" {
		t.Errorf("expected PUUID to be 'test-puuid', got %s", summoner.PUUID)
	}
	if summoner.SummonerName != "TestSummoner" {
		t.Errorf("expected SummonerName to be 'TestSummoner', got %s", summoner.SummonerName)
	}
	if summoner.SummonerLevel != 100 {
		t.Errorf("expected SummonerLevel to be 100, got %d", summoner.SummonerLevel)
	}
	if summoner.ProfileIconID != 1234 {
		t.Errorf("expected ProfileIconID to be 1234, got %d", summoner.ProfileIconID)
	}
}

func TestGetSummonerByName_APIError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":{"message":"Summoner not found"}}`)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	_, err := client.GetSummonerByName("na1", "NonExistent")
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetMatchesByPUUID_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.Header.Get("X-Riot-Token") != "test-api-key" {
				t.Errorf("expected X-Riot-Token header to be 'test-api-key', got %s", req.Header.Get("X-Riot-Token"))
			}
			
			// Return mock response
			responseBody := `["NA1_match1", "NA1_match2", "NA1_match3"]`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	matches, err := client.GetMatchesByPUUID("test-puuid", 3)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if len(matches) != 3 {
		t.Errorf("expected 3 matches, got %d", len(matches))
	}
	
	if matches[0] != "NA1_match1" {
		t.Errorf("expected first match to be 'NA1_match1', got %s", matches[0])
	}
}

func TestGetMatchesByPUUID_APIError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":{"message":"Bad request"}}`)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	_, err := client.GetMatchesByPUUID("invalid-puuid", 3)
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetMatchByID_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.Header.Get("X-Riot-Token") != "test-api-key" {
				t.Errorf("expected X-Riot-Token header to be 'test-api-key', got %s", req.Header.Get("X-Riot-Token"))
			}
			
			// Return mock response
			responseBody := `{
				"metadata": {
					"matchId": "NA1_match1",
					"participants": ["puuid1", "puuid2"]
				},
				"info": {
					"gameMode": "CLASSIC",
					"gameDuration": 1800,
					"gameCreation": 1609459200000
				}
			}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	match, err := client.GetMatchByID("NA1_match1")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if match.Metadata.MatchID != "NA1_match1" {
		t.Errorf("expected MatchID to be 'NA1_match1', got %s", match.Metadata.MatchID)
	}
	if match.Info.GameMode != "CLASSIC" {
		t.Errorf("expected GameMode to be 'CLASSIC', got %s", match.Info.GameMode)
	}
	if match.Info.GameDuration != 1800 {
		t.Errorf("expected GameDuration to be 1800, got %d", match.Info.GameDuration)
	}
}

func TestGetMatchByID_APIError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":{"message":"Match not found"}}`)),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-api-key", mockClient)
	_, err := client.GetMatchByID("invalid-match-id")
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
