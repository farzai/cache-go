package cache

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// RedisDriver implements the Cache interface using Redis as the storage driver.
type RedisDriver struct {
	client *redis.Client
}

// NewRedisDriver creates a new RedisDriver instance.
func NewRedisDriver(addr, password string, db int) *RedisDriver {
	return NewRedisDriverWithOptions(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// NewRedisDriverWithOptions creates a new RedisDriver instance with the specified options.
func NewRedisDriverWithOptions(opt *redis.Options) *RedisDriver {
	client := redis.NewClient(opt)
	return &RedisDriver{
		client: client,
	}
}

// Set stores the given value in the cache with the specified expiration time.
func (c *RedisDriver) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(key, data, expiration).Err()
}

// Get retrieves the value from the cache for the given key.
// If the key is found, the value is unmarshaled into the provided variable.
// Returns true if the key is found in the cache, false otherwise.
func (c *RedisDriver) Get(key string, value interface{}) (bool, error) {
	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

// OptimizedRedisDriver implements the Cache interface using Redis as the storage driver with performance optimizations.
type OptimizedRedisDriver struct {
	client *redis.Client
	cache  sync.Map
}

// NewOptimizedRedisDriver creates a new OptimizedRedisDriver instance.
func NewOptimizedRedisDriver(addr, password string, db int) *OptimizedRedisDriver {
	return NewOptimizedRedisDriverWithOptions(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// NewOptimizedRedisDriverWithOptions creates a new OptimizedRedisDriver instance with the specified options.
func NewOptimizedRedisDriverWithOptions(opt *redis.Options) *OptimizedRedisDriver {
	client := redis.NewClient(opt)

	return &OptimizedRedisDriver{
		client: client,
		cache:  sync.Map{},
	}
}

// Set stores the given value in the cache with the specified expiration time.
func (c *OptimizedRedisDriver) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.client.Set(key, data, expiration).Err()
	if err != nil {
		return err
	}

	c.cache.Store(key, value)

	return nil
}

// Get retrieves the value from the cache for the given key.
// If the key is found, the value is returned from the cache.
// Returns true if the key is found in the cache, false otherwise.
func (c *OptimizedRedisDriver) Get(key string, value interface{}) (bool, error) {
	cachedValue, ok := c.cache.Load(key)
	if ok {
		return true, json.Unmarshal(cachedValue.([]byte), value)
	}

	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return false, err
	}

	c.cache.Store(key, data)

	return true, nil
}
