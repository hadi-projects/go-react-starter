package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// CacheService defines the interface for cache operations
type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	FlushAll() error
	Close() error
	Status() string
}

// redisCache implements CacheService using Redis
type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache service
func NewRedisCache(host, port, password string, db int) (CacheService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Warning: Failed to connect to redis: %v. Using NoOpCache.\n", err)
		return &NoOpCache{}, nil
	}

	return &redisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache and unmarshals it into dest
func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()
	requestID := ""
	if rid, ok := ctx.Value(logger.CtxKeyRequestID).(string); ok {
		requestID = rid
	}
	val, err := r.client.Get(ctx, key).Result()
	elapsed := time.Since(start)

	// Truncate response for logging
	truncatedVal := logger.Truncate(val, 65536)

	status := 200
	if err != nil && err != redis.Nil {
		status = 500
	}

	logger.SystemLogger.Info().
		Str("method", "REDIS:GET").
		Str("path", key).
		Int("status_code", status).
		Int64("latency", elapsed.Milliseconds()).
		Str("response_body", truncatedVal).
		Msg("redis operation")

	if logger.SystemLogRepo != nil {
		_ = logger.SystemLogRepo.Create(ctx, &logger.SystemLog{
			RequestID:    requestID,
			Method:       "REDIS:GET",
			Path:         key,
			StatusCode:   status,
			Latency:      elapsed.Milliseconds(),
			ResponseBody: truncatedVal,
		})
	}

	if err == redis.Nil {
		return fmt.Errorf("cache miss: key not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Set stores a value in cache with the specified TTL
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	requestID := ""
	if rid, ok := ctx.Value(logger.CtxKeyRequestID).(string); ok {
		requestID = rid
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(r.ctx, key, jsonValue, ttl).Err()
	elapsed := time.Since(start)

	// Truncate request body for logging
	truncatedReq := logger.Truncate(string(jsonValue), 65536)

	status := 200
	if err != nil {
		status = 500
	}

	logger.SystemLogger.Info().
		Str("method", "REDIS:SET").
		Str("path", key).
		Int("status_code", status).
		Int64("latency", elapsed.Milliseconds()).
		Str("request_body", truncatedReq).
		Msg("redis operation")

	if logger.SystemLogRepo != nil {
		_ = logger.SystemLogRepo.Create(ctx, &logger.SystemLog{
			RequestID:   requestID,
			Method:      "REDIS:SET",
			Path:        key,
			StatusCode:  status,
			Latency:     elapsed.Milliseconds(),
			RequestBody: truncatedReq,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a specific key from cache
func (r *redisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// DeletePattern removes all keys matching the pattern
func (r *redisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}
	return nil
}

// FlushAll removes all keys from the current database
func (r *redisCache) FlushAll() error {
	if err := r.client.FlushDB(r.ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}

// Status returns the connection status
func (r *redisCache) Status() string {
	if err := r.client.Ping(r.ctx).Err(); err != nil {
		return "disconnected"
	}
	return "connected"
}

// Close closes the Redis connection
func (r *redisCache) Close() error {
	return r.client.Close()
}

// NoOpCache implements CacheService but does nothing
type NoOpCache struct{}

func (n *NoOpCache) Get(ctx context.Context, key string, dest interface{}) error {
	return fmt.Errorf("cache miss: no-op cache")
}

func (n *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (n *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *NoOpCache) DeletePattern(ctx context.Context, pattern string) error {
	return nil
}

func (n *NoOpCache) FlushAll() error {
	return nil
}

func (n *NoOpCache) Close() error {
	return nil
}

func (n *NoOpCache) Status() string {
	return "disconnected"
}
