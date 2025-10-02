package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient interface for Redis operations
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}

// Client wraps the Redis client
type Client struct {
	redis RedisClient
}

// NewClient creates a new cache client
func NewClient(redisURL, password string) (*Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	if password != "" {
		opt.Password = password
	}

	rdb := redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{redis: rdb}, nil
}

// NewClientWithRedis creates a new client with an existing Redis client
func NewClientWithRedis(redis RedisClient) *Client {
	return &Client{redis: redis}
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.redis.Close()
}

// Get retrieves a value from cache
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}
	return []byte(val), nil
}

// Set stores a value in cache with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.redis.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set in cache: %w", err)
	}

	return nil
}

// Delete removes a value from cache
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// GetJSON retrieves a JSON value from cache and unmarshals it
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil // Cache miss
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return true, nil
}
