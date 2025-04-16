package urls

import (
	"time"

	"github.com/mohammadne/fesghel/pkg/databases/postgres"
	"github.com/mohammadne/fesghel/pkg/databases/redis"
)

type Config struct {
	Redis                 *redis.Config    `required:"true"`
	Postgres              *postgres.Config `required:"true"`
	ShortURLLength        int              `required:"true" split_words:"true"`
	MaxRetriesOnCollision int              `required:"true" split_words:"true"`
	CacheExpiration       time.Duration    `required:"true" split_words:"true"`
}
