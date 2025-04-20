package urls

import (
	"context"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	var cacheTTL = 3 * time.Second

	t.Run("insert", func(t *testing.T) {
		t.Run("empty key", func(t *testing.T) {
			err := redisInstance.insert(context.TODO(), "", "any-value", time.Second)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		t.Run("check expiration", func(t *testing.T) {
			miniredisInstance.Set("key-1", "value-1")
			miniredisInstance.SetTTL("key-1", cacheTTL)

			err := redisInstance.insert(context.TODO(), "key-expiration", "value-expiration", cacheTTL)
			if err != nil {
				t.Error(err)
			}
		})
	})
}
