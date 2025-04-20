package urls

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	postgres_pkg "github.com/mohammadne/fesghel/pkg/databases/postgres"
	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

var (
	mockDatabase     sqlmock.Sqlmock
	postgresInstacne *postgres

	miniredisInstance *miniredis.Miniredis
	redisInstance     Redis

	serviceInstance *service
)

func TestMain(m *testing.M) {
	var err error

	{ // postgres
		var sqlDB *sql.DB

		sqlDB, mockDatabase, err = sqlmock.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start sqlmock: %v\n", err)
			os.Exit(1) // Exit with a non-zero status code
		}
		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		vectors := postgres_pkg.Vectors{
			Counter:   metrics_pkg.RegisterCounterNoop(),
			Histogram: metrics_pkg.RegisterHistogramNoop(),
		}

		postgresInstacne = &postgres{
			instance: &postgres_pkg.Postgres{DB: sqlxDB, Vectors: &vectors},
		}
	}

	{ // redis
		miniredisInstance, err = miniredis.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start miniredis: %v\n", err)
			os.Exit(1) // Exit with a non-zero status code
		}
		defer miniredisInstance.Close()

		cfg := redis_pkg.Config{Address: miniredisInstance.Addr(), Timeout: time.Second * 2}
		redisInstance, err = NewRedis(&cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open redis: %v\n", err)
			os.Exit(1) // Exit with a non-zero status code
		}
	}

	// service
	serviceInstance = &service{
		config: &Config{ShortURLLength: 6,
			MaxRetriesOnCollision: 3,
			CacheExpiration:       time.Second * 10,
		},
		logger:   zap.NewNop(),
		metrics:  newMetricsNoop(),
		postgres: postgresInstacne,
		redis:    redisInstance,
	}

	m.Run()
}
