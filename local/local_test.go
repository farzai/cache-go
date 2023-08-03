package local_test

import (
	"testing"
	"time"

	"github.com/farzai/cache-go"
)

func TestLocalFileDriver(t *testing.T) {
	t.Run("Should set value as string and get it back successfully", func(t *testing.T) {
		c := cache.NewLocalFileDriver("storage/test/")
		err := c.Set("test", "test", time.Microsecond*100)
		if err != nil {
			t.Error(err)
		}

		var value string
		found, err := c.Get("test", &value)
		if err != nil {
			t.Error(err)
		}

		if !found {
			t.Error("Expected key to be found")
		}

		if value != "test" {
			t.Errorf("Expected value to be 'test', got '%s'", value)
		}
	})
}
