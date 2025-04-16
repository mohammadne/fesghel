package urls

import (
	"go.uber.org/zap"

	"github.com/mohammadne/fesghel/internal/entities"
	"github.com/mohammadne/fesghel/pkg/databases/postgres"
	"github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Service interface {
}

type service struct {
	p *postgres.Postgres
	r *redis.Redis
}

func Initialize(cfg *Config, l *zap.Logger) (Service, error) {
	var svc = service{}

	postgresInstance, err := postgres.Open(cfg.Postgres, entities.Namespace, entities.System)
	if err != nil {
		l.Panic("error loading Postgres instance", zap.Error(err))
	}
	svc.p = postgresInstance

	redisInstance, err := redis.Open(cfg.Redis)
	if err != nil {
		l.Panic("error initializing Redis cache", zap.Error(err))
	}
	svc.r = redisInstance

	return nil, nil
}
