package urls

import (
	"time"

	Postgres_pkg "github.com/mohammadne/fesghel/pkg/databases/postgres"
	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Config struct {
	Redis                 *redis_pkg.Config    `required:"true"`
	Postgres              *Postgres_pkg.Config `required:"true"`
	ShortURLLength        int                  `required:"true" split_words:"true"`
	MaxRetriesOnCollision int                  `required:"true" split_words:"true"`
	CacheExpiration       time.Duration        `required:"true" split_words:"true"`
	BaseAddress           string               `required:"true" split_words:"true"`
}
