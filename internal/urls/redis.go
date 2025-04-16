package urls

import (
	"context"
	"errors"

	"github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Redis interface {
	insertIntoRedis(ctx context.Context, key, value string) error
	retrieveFromRedis(ctx context.Context, key string) (string, error)
}

var (
	errInsertDataToRedis = errors.New("error insert data to redis")
)

func (s *service) insertIntoRedis(ctx context.Context, key, value string) error {
	if err := s.r.Set(ctx, key, value, s.config.CacheExpiration).Err(); err != nil {
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return errors.Join(errInsertDataToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return nil
}

var (
	errKeyNotFound = errors.New("errKeyNotFound")
)

func (s *service) retrieveFromRedis(ctx context.Context, key string) (string, error) {
	dataString, err := s.r.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errKeyNotFound
		}
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return "", errors.Join(errInsertDataToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return dataString, nil
}
