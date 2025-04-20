package urls

import (
	"context"
	"errors"
	"time"

	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Redis interface {
	insert(ctx context.Context, key, value string, expiration time.Duration) error
	retrieve(ctx context.Context, key string) (string, error)
}

type redis struct {
	instance *redis_pkg.Redis
}

func NewRedis(cfg *redis_pkg.Config) (Redis, error) {
	instance, err := redis_pkg.Open(cfg)
	if err != nil {
		return nil, err
	}
	return &redis{instance: instance}, nil
}

var (
	errInsertDataToRedis = errors.New("error insert data to redis")
)

func (s *redis) insert(ctx context.Context, key, value string, expiration time.Duration) error {
	if err := s.instance.Set(ctx, key, value, expiration).Err(); err != nil {
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return errors.Join(errInsertDataToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return nil
}

var (
	errKeyNotFound = errors.New("errKeyNotFound")
)

func (s *redis) retrieve(ctx context.Context, key string) (string, error) {
	dataString, err := s.instance.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis_pkg.Nil) {
			return "", errKeyNotFound
		}
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return "", errors.Join(errInsertDataToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return dataString, nil
}
