package urls

import (
	"context"
	"testing"
	"time"
)

const cacheTTL = 3 * time.Second

func TestRedisInsert(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		err := redisInstance.insert(context.TODO(), "", "any-value", time.Second)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestRedisRetrieve(t *testing.T) {
	t.Run("check expiration", func(t *testing.T) {
		miniredisInstance.Set("key-1", "value-1")
		miniredisInstance.SetTTL("key-1", cacheTTL)

		err := redisInstance.insert(context.TODO(), "key-expiration", "value-expiration", cacheTTL)
		if err != nil {
			t.Error(err)
		}
	})
}
