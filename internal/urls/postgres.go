package urls

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

type Postgers interface {
	insertIntoOracle(ctx context.Context, id, value string, timestamp time.Time) (err error)
	retrieveFromOracle(ctx context.Context, id string) (value string, err error)
}

var (
	errUniqueConstraintViolated = errors.New("error duplicate key")
	errInsertingValue           = errors.New("error inserting value")
)

func (s *service) insertIntoPostgres(ctx context.Context, id, url string, timestamp time.Time) (err error) {
	defer func(start time.Time) {
		if err != nil {
			s.p.Vectors.Counter.IncrementVector("urls", "insert", metrics_pkg.StatusFailure)
			return
		}
		s.p.Vectors.Counter.IncrementVector("urls", "insert", metrics_pkg.StatusSuccess)
		s.p.Vectors.Histogram.ObserveResponseTime(start, "urls", "insert")
	}(time.Now())

	query := `
	INSERT INTO urls (id, url, created_at)
	VALUES ($1, $2, $3)`

	_, err = s.p.ExecContext(ctx, query, id, url, timestamp)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errUniqueConstraintViolated
		}
		return errors.Join(errInsertingValue, err)
	}

	return nil
}

var (
	ErrIDNotExists     = errors.New("error id not exists")
	ErrRetreivingValue = errors.New("error retreiving value")
)

func (s *service) retrieveFromOracle(ctx context.Context, id string) (value string, err error) {
	defer func(start time.Time) {
		if err != nil {
			s.p.Vectors.Counter.IncrementVector("urls", "retrieve", metrics_pkg.StatusFailure)
			return
		}
		s.p.Vectors.Counter.IncrementVector("urls", "retrieve", metrics_pkg.StatusSuccess)
		s.p.Vectors.Histogram.ObserveResponseTime(start, "urls", "retrieve")
	}(time.Now())

	query := `
	SELECT url
	FROM urls
	WHERE id = $1`

	err = s.p.QueryRowContext(ctx, query, id).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrIDNotExists
		}
		return "", errors.Join(ErrRetreivingValue, err)
	}

	return value, nil
}
