package urls

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

const cacheTTL = 3 * time.Second

func TestRedisInsert(t *testing.T) {
	var (
		sampleID  = "sample-id"
		sampleURL = "sample-url"
	)

	t.Run("empty parameters", func(t *testing.T) {
		t.Run("empty id", func(t *testing.T) {
			err := redisInstance.insert(context.TODO(), "", sampleURL, cacheTTL)
			if !errors.Is(err, errInvalidInsertParameters) {
				t.Error(err)
			}
		})

		t.Run("empty url", func(t *testing.T) {
			err := redisInstance.insert(context.TODO(), sampleID, "", cacheTTL)
			if !errors.Is(err, errInvalidInsertParameters) {
				t.Error(err)
			}
		})
	})

	t.Run("valid insert", func(t *testing.T) {
		err := redisInstance.insert(context.TODO(), sampleID, sampleURL, cacheTTL)
		if err != nil {
			t.Error(err)
		}

		url, err := miniredisInstance.Get(sampleID)
		if err != nil {
			t.Error(err)
		}

		if url != sampleURL {
			t.Error("invalid url has been returned")
		}
	})

	t.Run("check ttl", func(t *testing.T) {
		err := redisInstance.insert(context.TODO(), sampleID, sampleURL, cacheTTL)
		if err != nil {
			t.Error(err)
		}

		_, err = miniredisInstance.Get(sampleID)
		if err != nil {
			t.Error(err)
		}

		miniredisInstance.FastForward(cacheTTL)

		_, err = miniredisInstance.Get(sampleID)
		if !errors.Is(err, miniredis.ErrKeyNotFound) {
			t.Errorf("expecting miniredis ErrKeyNotFound error but got something else: %v", err)
		}
	})
}

func TestRedisRetrieve(t *testing.T) {
	var (
		sampleID  = "sample-id"
		sampleURL = "sample-url"
	)

	t.Run("with empty id", func(t *testing.T) {
		_, err := redisInstance.retrieve(context.TODO(), "")
		if !errors.Is(err, errInvalidRetrieveParameters) {
			t.Error(err)
		}
	})

	t.Run("valid result", func(t *testing.T) {
		// miniredisInstance.FlushAll()
		miniredisInstance.Set(sampleID, sampleURL)
		miniredisInstance.SetTTL(sampleID, cacheTTL)

		url, err := redisInstance.retrieve(context.TODO(), sampleID)
		if err != nil {
			t.Error(err)
		}

		if url != sampleURL {
			t.Error("invalid url has been returned")
		}
	})

	t.Run("check ttl", func(t *testing.T) {
		miniredisInstance.Set(sampleID, sampleURL)
		miniredisInstance.SetTTL(sampleID, cacheTTL)

		miniredisInstance.FastForward(cacheTTL)

		_, err := redisInstance.retrieve(context.TODO(), sampleID)
		if !errors.Is(err, errIDNotFound) {
			t.Errorf("expecting errIDNotFound error but got something else: %v", err)
		}
	})
}
