package cache

import (
	"time"
)

// Cache interface defines the methods for interacting with the cache.
type Cache interface {
	// Set stores the given value in the cache with the specified expiration time.
	Set(key string, value interface{}, expiration time.Duration) error

	// Get retrieves the value from the cache with the specified key.
	Get(key string) (*CacheItem, error)

	// Delete removes the value from the cache with the specified key.
	Delete(key string) error

	// Flush removes all values from the cache.
	Flush() error

	// Has checks if the cache has the specified key.
	Has(key string) bool
}

// CacheItem represents a single cache item.
type CacheItem struct {
	Key        string        `json:"key"`
	Value      interface{}   `json:"value"`
	Expiration time.Duration `json:"expiration"`
	CreatedAt  time.Time     `json:"created_at"`
}

func (c *CacheItem) Expired() bool {
	return time.Now().After(c.CreatedAt.Add(c.Expiration))
}

func (c *CacheItem) ExpireIn(duration time.Duration) {
	c.Expiration = duration
}

func (c *CacheItem) ExpireAt(t time.Time) {
	c.Expiration = time.Since(t)
}

func (c *CacheItem) Set(value interface{}) {
	c.Value = value
}

func (c *CacheItem) Get() interface{} {
	return c.Value
}
