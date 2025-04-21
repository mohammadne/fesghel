package urls

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/mohammadne/fesghel/internal/entities"
	postgres_pkg "github.com/mohammadne/fesghel/pkg/databases/postgres"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

type Postgres interface {
	insert(ctx context.Context, id, url string, timestamp time.Time) (err error)
	retrieve(ctx context.Context, id string) (url string, err error)
}

type postgres struct {
	instance *postgres_pkg.Postgres
}

func NewPostgres(cfg *postgres_pkg.Config) (Postgres, error) {
	instance, err := postgres_pkg.Open(cfg, entities.Namespace, entities.System)
	if err != nil {
		return nil, err
	}
	return &postgres{instance: instance}, nil
}

var (
	errUniqueConstraintViolated = errors.New("error duplicate key")
	errInsertingURL             = errors.New("error inserting url")
)

const (
	queryInsert = `
	INSERT INTO urls (id, url, created_at)
	VALUES ($1, $2, $3)`
)

func (s *postgres) insert(ctx context.Context, id, url string, timestamp time.Time) (err error) {
	defer func(start time.Time) {
		if err != nil {
			s.instance.Vectors.Counter.IncrementVector("urls", "insert", metrics_pkg.StatusFailure)
			return
		}
		s.instance.Vectors.Counter.IncrementVector("urls", "insert", metrics_pkg.StatusSuccess)
		s.instance.Vectors.Histogram.ObserveResponseTime(start, "urls", "insert")
	}(time.Now())

	_, err = s.instance.ExecContext(ctx, queryInsert, id, url, timestamp)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errUniqueConstraintViolated
		}
		return errors.Join(errInsertingURL, err)
	}

	return nil
}

var (
	ErrIDNotExists     = errors.New("error id not exists")
	ErrRetreivingValue = errors.New("error retreiving url")
)

const (
	queryRetrieve = `
	SELECT url
	FROM urls
	WHERE id = $1`
)

func (s *postgres) retrieve(ctx context.Context, id string) (url string, err error) {
	defer func(start time.Time) {
		if err != nil {
			s.instance.Vectors.Counter.IncrementVector("urls", "retrieve", metrics_pkg.StatusFailure)
			return
		}
		s.instance.Vectors.Counter.IncrementVector("urls", "retrieve", metrics_pkg.StatusSuccess)
		s.instance.Vectors.Histogram.ObserveResponseTime(start, "urls", "retrieve")
	}(time.Now())

	err = s.instance.QueryRowContext(ctx, queryRetrieve, id).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrIDNotExists
		}
		return "", errors.Join(ErrRetreivingValue, err)
	}

	return url, nil
}
