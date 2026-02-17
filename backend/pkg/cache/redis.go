package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService defines the interface for cache operations
type CacheService interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	DeletePattern(pattern string) error
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
func (r *redisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
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
func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := r.client.Set(r.ctx, key, jsonValue, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a specific key from cache
func (r *redisCache) Delete(key string) error {
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// DeletePattern removes all keys matching the pattern
func (r *redisCache) DeletePattern(pattern string) error {
	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()
	for iter.Next(r.ctx) {
		if err := r.client.Del(r.ctx, iter.Val()).Err(); err != nil {
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
	return "connected"
}

// Close closes the Redis connection
func (r *redisCache) Close() error {
	return r.client.Close()
}

// NoOpCache implements CacheService but does nothing
type NoOpCache struct{}

func (n *NoOpCache) Get(key string, dest interface{}) error {
	return fmt.Errorf("cache miss: no-op cache")
}

func (n *NoOpCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (n *NoOpCache) Delete(key string) error {
	return nil
}

func (n *NoOpCache) DeletePattern(pattern string) error {
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
