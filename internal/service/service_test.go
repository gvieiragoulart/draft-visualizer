package service

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gvieiragoulart/draft-visualizer/internal/riot"
	"github.com/redis/go-redis/v9"
)

func TestNewService(t *testing.T) {
	mockRiot := riot.NewClient("test-key")

	service := NewService(mockRiot)
	if service == nil {
		t.Fatal("expected service to be created")
	}
}

// mockRedisClient implements the cache.RedisClient interface
type mockRedisClient struct{}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	cmd.SetErr(redis.Nil)
	return cmd
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("OK")
	return cmd
}

func (m *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)
	return cmd
}

func (m *mockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("PONG")
	return cmd
}

func (m *mockRedisClient) Close() error {
	return nil
}

// mockHTTPClient implements riot.HTTPClient interface
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
	}, nil
}

func TestGetSummoner_FromAPI(t *testing.T) {
	expectedSummoner := &riot.Summoner{
		PUUID:         "test-puuid",
		SummonerName:  "TestSummoner",
		SummonerLevel: 100,
		ProfileIconID: 1234,
	}

	mockHTTP := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
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

	riotClient := riot.NewClientWithHTTPClient("test-key", mockHTTP)

	service := NewService(riotClient)
	ctx := context.Background()

	summoner, err := service.GetSummoner(ctx, "na1", "TestSummoner")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summoner.PUUID != expectedSummoner.PUUID {
		t.Errorf("expected PUUID to be %s, got %s", expectedSummoner.PUUID, summoner.PUUID)
	}
}

func TestGetMatches_FromAPI(t *testing.T) {
	expectedMatches := []string{"NA1_match1", "NA1_match2", "NA1_match3"}

	mockHTTP := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `["NA1_match1", "NA1_match2", "NA1_match3"]`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		},
	}

	riotClient := riot.NewClientWithHTTPClient("test-key", mockHTTP)

	service := NewService(riotClient)
	ctx := context.Background()

	matches, err := service.GetMatches(ctx, "test-puuid", 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(matches) != len(expectedMatches) {
		t.Errorf("expected %d matches, got %d", len(expectedMatches), len(matches))
	}
}

func TestGetMatch_FromAPI(t *testing.T) {
	mockHTTP := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
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

	riotClient := riot.NewClientWithHTTPClient("test-key", mockHTTP)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	// First query - GetMatchByID returns no rows (not found)
	mock.ExpectQuery(`SELECT id, match_id`).WillReturnError(sql.ErrNoRows)

	// Second query - SaveMatch inserts the match
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())
	mock.ExpectQuery(`INSERT INTO matches`).WillReturnRows(rows)

	service := NewService(riotClient)
	ctx := context.Background()

	match, err := service.GetMatch(ctx, "NA1_match1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if match.Metadata.MatchID != "NA1_match1" {
		t.Errorf("expected match ID to be NA1_match1, got %s", match.Metadata.MatchID)
	}
}
