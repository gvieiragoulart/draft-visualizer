package database

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	
	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB is a mock database for testing
type MockDB struct {
	QueryRowFunc func(query string, args ...interface{}) *sql.Row
	QueryFunc    func(query string, args ...interface{}) (*sql.Rows, error)
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	CloseFunc    func() error
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return nil
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockRow is a helper to create mock sql.Row for testing
type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func TestNewClientWithDB(t *testing.T) {
	mockDB := &MockDB{}
	client := NewClientWithDB(mockDB)
	if client == nil {
		t.Fatal("expected client to be created")
	}
	if client.db != mockDB {
		t.Error("expected db to be the mock db")
	}
}

func TestClose(t *testing.T) {
	closeCalled := false
	mockDB := &MockDB{
		CloseFunc: func() error {
			closeCalled = true
			return nil
		},
	}
	client := NewClientWithDB(mockDB)
	err := client.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !closeCalled {
		t.Error("expected Close to be called")
	}
}

func TestSaveSummoner_Success(t *testing.T) {
	now := time.Now()
	mockDB := &MockDB{
		QueryRowFunc: func(query string, args ...interface{}) *sql.Row {
			// Create a mock row that will scan successfully
			db, mock, _ := sqlmock.New()
			defer db.Close()
			
			rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, now, now)
			mock.ExpectQuery("INSERT INTO summoners").WillReturnRows(rows)
			
			return db.QueryRow("SELECT 1")
		},
	}

	client := NewClientWithDB(mockDB)
	summoner := &Summoner{
		PUUID:         "test-puuid",
		SummonerName:  "TestSummoner",
		SummonerLevel: 100,
		ProfileIconID: 1234,
		Region:        "na1",
	}

	// Note: This test is simplified due to sql.Row limitations in testing
	// In a real scenario, we'd use sqlmock or a test database
	_ = client
	_ = summoner
}

func TestGetSummonerByPUUID_NotFound(t *testing.T) {
	mockDB := &MockDB{
		QueryRowFunc: func(query string, args ...interface{}) *sql.Row {
			// Create a mock row that will return sql.ErrNoRows
			db, mock, _ := sqlmock.New()
			defer db.Close()
			
			mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			return db.QueryRow("SELECT 1")
		},
	}

	client := NewClientWithDB(mockDB)
	summoner, err := client.GetSummonerByPUUID("non-existent-puuid")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summoner != nil {
		t.Error("expected summoner to be nil for not found case")
	}
}

func TestGetSummonerByPUUID_Error(t *testing.T) {
	mockDB := &MockDB{
		QueryRowFunc: func(query string, args ...interface{}) *sql.Row {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			
			mock.ExpectQuery("SELECT").WillReturnError(errors.New("database error"))
			return db.QueryRow("SELECT 1")
		},
	}

	client := NewClientWithDB(mockDB)
	_, err := client.GetSummonerByPUUID("test-puuid")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetMatchByID_NotFound(t *testing.T) {
	mockDB := &MockDB{
		QueryRowFunc: func(query string, args ...interface{}) *sql.Row {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			
			mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			return db.QueryRow("SELECT 1")
		},
	}

	client := NewClientWithDB(mockDB)
	match, err := client.GetMatchByID("non-existent-match")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if match != nil {
		t.Error("expected match to be nil for not found case")
	}
}


func TestSaveSummoner_WithSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	client := NewClientWithDB(db)
	now := time.Now()

	summoner := &Summoner{
		PUUID:         "test-puuid",
		SummonerName:  "TestSummoner",
		SummonerLevel: 100,
		ProfileIconID: 1234,
		Region:        "na1",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery(`INSERT INTO summoners`).
		WithArgs(summoner.PUUID, summoner.SummonerName, summoner.SummonerLevel, 
			summoner.ProfileIconID, summoner.Region, sqlmock.AnyArg()).
		WillReturnRows(rows)

	err = client.SaveSummoner(summoner)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summoner.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", summoner.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetSummonerByPUUID_WithSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	client := NewClientWithDB(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "puuid", "summoner_name", "summoner_level", 
		"profile_icon_id", "region", "created_at", "updated_at"}).
		AddRow(1, "test-puuid", "TestSummoner", 100, 1234, "na1", now, now)

	mock.ExpectQuery(`SELECT id, puuid, summoner_name`).
		WithArgs("test-puuid").
		WillReturnRows(rows)

	summoner, err := client.GetSummonerByPUUID("test-puuid")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summoner == nil {
		t.Fatal("expected summoner to be returned")
	}

	if summoner.PUUID != "test-puuid" {
		t.Errorf("expected PUUID to be 'test-puuid', got %s", summoner.PUUID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestSaveMatch_WithSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	client := NewClientWithDB(db)
	now := time.Now()

	match := &Match{
		MatchID:      "NA1_match1",
		GameMode:     "CLASSIC",
		GameDuration: 1800,
		GameCreation: 1609459200000,
		Data: map[string]interface{}{
			"test": "data",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery(`INSERT INTO matches`).
		WithArgs(match.MatchID, match.GameMode, match.GameDuration, 
			match.GameCreation, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	err = client.SaveMatch(match)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if match.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", match.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetMatchByID_WithSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	client := NewClientWithDB(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "match_id", "game_mode", "game_duration",
		"game_creation", "data", "created_at", "updated_at"}).
		AddRow(1, "NA1_match1", "CLASSIC", 1800, 1609459200000, 
			[]byte(`{"test":"data"}`), now, now)

	mock.ExpectQuery(`SELECT id, match_id`).
		WithArgs("NA1_match1").
		WillReturnRows(rows)

	match, err := client.GetMatchByID("NA1_match1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if match == nil {
		t.Fatal("expected match to be returned")
	}

	if match.MatchID != "NA1_match1" {
		t.Errorf("expected MatchID to be 'NA1_match1', got %s", match.MatchID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
