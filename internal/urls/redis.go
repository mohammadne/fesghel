package urls

import (
	"context"
	"errors"
	"time"

	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Redis interface {
	insert(ctx context.Context, id, url string, expiration time.Duration) error
	retrieve(ctx context.Context, id string) (url string, err error)
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
	errInvalidInsertParameters = errors.New("error Invalid Insert Parameters")
	errInsertURLToRedis        = errors.New("error insert url to redis")
)

func (s *redis) insert(ctx context.Context, id, url string, expiration time.Duration) error {
	if len(id) == 0 || len(url) == 0 {
		return errInvalidInsertParameters
	}

	if err := s.instance.Set(ctx, id, url, expiration).Err(); err != nil {
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return errors.Join(errInsertURLToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return nil
}

var (
	errInvalidRetrieveParameters = errors.New("error Invalid Insert Parameters")
	errIDNotFound                = errors.New("errIDNotFound")
)

func (s *redis) retrieve(ctx context.Context, id string) (string, error) {
	if len(id) == 0 {
		return "", errInvalidRetrieveParameters
	}

	url, err := s.instance.Get(ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis_pkg.Nil) {
			return "", errIDNotFound
		}
		// s.metrics.counter.WithLabelValues("SetInformation", "failure").Inc()
		return "", errors.Join(errInsertURLToRedis, err)
	}

	// s.metrics.counter.WithLabelValues("SetInformation", "success").Inc()
	return url, nil
}
