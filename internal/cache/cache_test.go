package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockRedisClient is a mock Redis client for testing
type MockRedisClient struct {
	GetFunc  func(ctx context.Context, key string) *redis.StringCmd
	SetFunc  func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	DelFunc  func(ctx context.Context, keys ...string) *redis.IntCmd
	PingFunc func(ctx context.Context) *redis.StatusCmd
	CloseFunc func() error
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return redis.NewStringCmd(ctx)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, keys...)
	}
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestNewClientWithRedis(t *testing.T) {
	mockRedis := &MockRedisClient{}
	client := NewClientWithRedis(mockRedis)
	if client == nil {
		t.Fatal("expected client to be created")
	}
	if client.redis != mockRedis {
		t.Error("expected redis to be the mock redis")
	}
}

func TestClose(t *testing.T) {
	closeCalled := false
	mockRedis := &MockRedisClient{
		CloseFunc: func() error {
			closeCalled = true
			return nil
		},
	}
	client := NewClientWithRedis(mockRedis)
	err := client.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !closeCalled {
		t.Error("expected Close to be called")
	}
}

func TestGet_Success(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal(`{"test":"value"}`)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	data, err := client.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != `{"test":"value"}` {
		t.Errorf("expected data to be '{\"test\":\"value\"}', got %s", string(data))
	}
}

func TestGet_CacheMiss(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetErr(redis.Nil)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	data, err := client.Get(ctx, "non-existent-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data != nil {
		t.Error("expected data to be nil for cache miss")
	}
}

func TestGet_Error(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetErr(errors.New("redis error"))
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	_, err := client.Get(ctx, "test-key")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSet_Success(t *testing.T) {
	setCalled := false
	var capturedKey string
	var capturedExpiration time.Duration

	mockRedis := &MockRedisClient{
		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			setCalled = true
			capturedKey = key
			capturedExpiration = expiration
			cmd := redis.NewStatusCmd(ctx)
			cmd.SetVal("OK")
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	testData := map[string]string{"test": "value"}
	err := client.Set(ctx, "test-key", testData, 5*time.Minute)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !setCalled {
		t.Error("expected Set to be called")
	}

	if capturedKey != "test-key" {
		t.Errorf("expected key to be 'test-key', got %s", capturedKey)
	}

	if capturedExpiration != 5*time.Minute {
		t.Errorf("expected expiration to be 5 minutes, got %v", capturedExpiration)
	}
}

func TestSet_Error(t *testing.T) {
	mockRedis := &MockRedisClient{
		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			cmd := redis.NewStatusCmd(ctx)
			cmd.SetErr(errors.New("redis error"))
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	testData := map[string]string{"test": "value"}
	err := client.Set(ctx, "test-key", testData, 5*time.Minute)
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_Success(t *testing.T) {
	delCalled := false
	mockRedis := &MockRedisClient{
		DelFunc: func(ctx context.Context, keys ...string) *redis.IntCmd {
			delCalled = true
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(1)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	err := client.Delete(ctx, "test-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !delCalled {
		t.Error("expected Del to be called")
	}
}

func TestDelete_Error(t *testing.T) {
	mockRedis := &MockRedisClient{
		DelFunc: func(ctx context.Context, keys ...string) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetErr(errors.New("redis error"))
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	err := client.Delete(ctx, "test-key")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetJSON_Success(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal(`{"name":"test","value":123}`)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	
	var result TestStruct
	found, err := client.GetJSON(ctx, "test-key", &result)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !found {
		t.Fatal("expected cache hit, got miss")
	}

	if result.Name != "test" {
		t.Errorf("expected Name to be 'test', got %s", result.Name)
	}

	if result.Value != 123 {
		t.Errorf("expected Value to be 123, got %d", result.Value)
	}
}

func TestGetJSON_CacheMiss(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetErr(redis.Nil)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	var result map[string]interface{}
	found, err := client.GetJSON(ctx, "non-existent-key", &result)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found {
		t.Error("expected cache miss, got hit")
	}
}

func TestGetJSON_InvalidJSON(t *testing.T) {
	mockRedis := &MockRedisClient{
		GetFunc: func(ctx context.Context, key string) *redis.StringCmd {
			cmd := redis.NewStringCmd(ctx)
			cmd.SetVal(`invalid json`)
			return cmd
		},
	}

	client := NewClientWithRedis(mockRedis)
	ctx := context.Background()
	
	var result map[string]interface{}
	_, err := client.GetJSON(ctx, "test-key", &result)
	
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
