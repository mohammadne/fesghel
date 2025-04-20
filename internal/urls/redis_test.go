package urls

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
)

func TestRedis(t *testing.T) {
	redisInstance, err := miniredis.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start miniredis: %v\n", err)
		os.Exit(1) // Exit with a non-zero status code
	}
	defer redisInstance.Close()

	cfg := redis_pkg.Config{Address: redisInstance.Addr(), Timeout: time.Second * 2}
	redis, err := newRedis(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open redis: %v\n", err)
		os.Exit(1) // Exit with a non-zero status code
	}

	var cacheTTL = 3 * time.Second

	t.Run("insert", func(t *testing.T) {
		t.Run("empty key", func(t *testing.T) {
			err := redis.insert(context.TODO(), "", "any-value", time.Second)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		t.Run("check expiration", func(t *testing.T) {
			redisInstance.Set("key-1", string(marshaledItem))
			redisInstance.SetTTL("key-1", cacheTTL)

			err := redis.insert(context.TODO(), "key-expiration", "value-expiration", cacheTTL)
			if err != nil {
				t.Error(err)
			}
		})
	})
}
