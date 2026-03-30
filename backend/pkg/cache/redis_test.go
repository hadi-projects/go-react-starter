package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoOpCache(t *testing.T) {
	cache := &NoOpCache{}

	t.Run("Get returns error", func(t *testing.T) {
		var dest string
		err := cache.Get(context.Background(), "key", &dest)
		assert.Error(t, err)
		assert.Equal(t, "cache miss: no-op cache", err.Error())
	})

	t.Run("Set returns nil", func(t *testing.T) {
		err := cache.Set(context.Background(), "key", "value", time.Minute)
		assert.NoError(t, err)
	})

	t.Run("Delete returns nil", func(t *testing.T) {
		err := cache.Delete(context.Background(), "key")
		assert.NoError(t, err)
	})

	t.Run("DeletePattern returns nil", func(t *testing.T) {
		err := cache.DeletePattern(context.Background(), "pattern")
		assert.NoError(t, err)
	})

	t.Run("FlushAll returns nil", func(t *testing.T) {
		err := cache.FlushAll()
		assert.NoError(t, err)
	})

	t.Run("Close returns nil", func(t *testing.T) {
		err := cache.Close()
		assert.NoError(t, err)
	})

	t.Run("Status returns disconnected", func(t *testing.T) {
		status := cache.Status()
		assert.Equal(t, "disconnected", status)
	})
}

func TestNewRedisCache_Fallback(t *testing.T) {
	// Try to connect to an invalid port
	cacheService, err := NewRedisCache("localhost", "12345", "", 0)

	assert.NoError(t, err)
	assert.NotNil(t, cacheService)
	assert.IsType(t, &NoOpCache{}, cacheService)
	assert.Equal(t, "disconnected", cacheService.Status())
}
