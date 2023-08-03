package local

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/farzai/cache-go"
)

var _ cache.Cache = (*LocalFileDriver)(nil)

// LocalFileDriver implements the Cache interface using the local file system as the storage driver.
type LocalFileDriver struct {
	cache.Cache
	mutex sync.RWMutex
	path  string
}

// NewLocalFileDriver creates a new LocalFileDriver instance with the specified file path.
func NewLocalFileDriver(path string) *LocalFileDriver {
	// Ensure the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}

	dirSeparator := string(os.PathSeparator)

	return &LocalFileDriver{
		path: strings.TrimRight(path, dirSeparator) + dirSeparator,
	}
}

// Get retrieves the value from the cache for the given key.
// If the key is found, the value is unmarshaled into the provided variable.
// Returns true if the key is found in the cache, false otherwise.
func (c *LocalFileDriver) Get(key string) (*cache.CacheItem, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	path := c.filePathForKey(key)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	for {
		var cacheData map[string]interface{}
		if err := dec.Decode(&cacheData); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		cacheKey, ok := cacheData["key"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid cache data format")
		}

		if cacheKey == key {
			expiration, ok := cacheData["expiration"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid cache data format")
			}

			if int64(expiration) > time.Now().UnixNano() {
				cacheValue, ok := cacheData["value"].([]byte)
				if !ok {
					return nil, fmt.Errorf("invalid cache data format")
				}

				item := new(cache.CacheItem)

				err = json.Unmarshal(cacheValue, item)
				if err != nil {
					return nil, err
				}

				return item, nil
			}
		}
	}

	return nil, nil
}

// Set stores the given value in the cache with the specified expiration time.
func (c *LocalFileDriver) Set(key string, value interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if value == nil {
		return nil
	}

	cacheData := map[string]interface{}{
		"key":        key,
		"value":      value,
		"expiration": time.Now().Add(expiration).UnixNano(),
		"created_at": time.Now().UnixNano(),
	}

	path := c.filePathForKey(key)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(cacheData)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes the value from the cache with the specified key.
func (c *LocalFileDriver) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	path := c.filePathForKey(key)

	// Check if the file exists before deleting it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Delete the file
	return os.Remove(path)
}

func (c *LocalFileDriver) filePathForKey(key string) string {
	// Encode key to base64 and encode to md5 to avoid special characters in the file name.
	filename := fmt.Sprintf("%x", md5.Sum([]byte(base64.StdEncoding.EncodeToString([]byte(key)))))

	return c.path + filename
}
