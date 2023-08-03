# Cache - GO

A simple cache library for golang.

## Installation

```bash
go get github.com/farzai/cache-go
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/farzai/cache-go"
)


func main() {
	cache := cache.NewCache(cache.NewRedisDriver("localhost:6379", "", 0, 10*time.Second))
	// Or
	// cache := cache.NewCache(cache.NewMemoryDriver(10*time.Second))
	// cache := cache.NewCache(cache.NewLocalFileDriver("/storage/cache", 10*time.Second))

	err := cache.Set("key", "value")
	if err != nil {
		fmt.Println("Error setting value in cache:", err)
		return
	}

	value, err := cache.Get("key")
	if err != nil {
		fmt.Println("Error getting value from cache:", err)
		return
	}

	fmt.Println("Value from cache:", value)
}
```

## License
Please see the [LICENSE](LICENSE) file for more information.
