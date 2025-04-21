package urls

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	postgres_pkg "github.com/mohammadne/fesghel/pkg/databases/postgres"
	redis_pkg "github.com/mohammadne/fesghel/pkg/databases/redis"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

var (
	mockDatabase     sqlmock.Sqlmock
	postgresInstacne Postgres
	postgresMock     *mockPostgres

	miniredisInstance *miniredis.Miniredis
	redisInstance     Redis
	redisMock         *mockRedis

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
		postgresMock = new(mockPostgres)
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
		redisMock = new(mockRedis)
	}

	// service
	serviceInstance = &service{
		config: &Config{ShortURLLength: 6,
			MaxRetriesOnCollision: 3,
			CacheExpiration:       time.Second * 10,
		},
		logger:   zap.NewNop(),
		metrics:  newMetricsNoop(),
		postgres: postgresMock,
		redis:    redisMock,
	}

	m.Run()
}

type mockPostgres struct{ mock.Mock }

func (m *mockPostgres) insert(ctx context.Context, id, url string, timestamp time.Time) (err error) {
	args := m.Called(ctx, id, url, timestamp)
	return args.Error(0)
}

func (m *mockPostgres) retrieve(ctx context.Context, id string) (url string, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(string), args.Error(1)
}

type mockRedis struct{ mock.Mock }

func (m *mockRedis) insert(ctx context.Context, id, url string, expiration time.Duration) error {
	args := m.Called(ctx, id, url, expiration)
	return args.Error(0)
}

func (m *mockRedis) retrieve(ctx context.Context, id string) (url string, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(string), args.Error(1)
}
