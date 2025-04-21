package urls

import (
	"context"
	"crypto/sha256"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/mohammadne/fesghel/internal/entities"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

type Service interface {
	// Shorten shortenes the url by giving a url string, then returns the shortened id
	Shorten(ctx context.Context, url entities.URL) (string, error)

	// Retrieve returns the actual url by giving url's shortened id
	Retrieve(ctx context.Context, id string) (entities.URL, error)
}

type service struct {
	config   *Config
	logger   *zap.Logger
	metrics  *metrics
	postgres Postgres
	redis    Redis
}

func NewService(cfg *Config, l *zap.Logger) (Service, error) {
	var svc = service{config: cfg, logger: l}

	metrics, err := newMetrics()
	if err != nil {
		l.Panic("error loading Postgres instance", zap.Error(err))
	}
	svc.metrics = metrics

	postgres, err := NewPostgres(cfg.Postgres)
	if err != nil {
		l.Panic("error loading Postgres instance", zap.Error(err))
	}
	svc.postgres = postgres

	redis, err := NewRedis(cfg.Redis)
	if err != nil {
		l.Panic("error initializing Redis cache", zap.Error(err))
	}
	svc.redis = redis

	return &svc, nil
}

var (
	ErrInsertingIntoPostgres  = errors.New("error inserting value into postgres")
	ErrMaxRetriesForCollision = errors.New("max retries exceeded while generating unique key")
)

// Shorten stores the data and returns the shorten key
// 1. generate key
// 2. store on oracle
// 3. retry on conflicts
func (s *service) Shorten(ctx context.Context, url entities.URL) (key string, err error) {
	defer func(start time.Time) {
		var status = metrics_pkg.StatusFailure
		if err == nil {
			s.metrics.Histogram.ObserveResponseTime(start, "shorten")
			status = metrics_pkg.StatusSuccess
		}
		s.metrics.Counter.IncrementVector("shorten", status)
	}(time.Now())

	for attempt := 1; attempt <= s.config.MaxRetriesOnCollision; attempt++ {
		timestamp := time.Now()

		key = s.generateKey(string(url), timestamp)

		err = s.postgres.insert(ctx, key, string(url), timestamp)
		if err == nil {
			_ = s.redis.insert(ctx, key, string(url), s.config.CacheExpiration)
			return key, nil // success
		}

		if errors.Is(err, errUniqueConstraintViolated) {
			// Collision: retry with new key
			continue
		}

		// Some other DB error
		return "", errors.Join(ErrInsertingIntoPostgres, err)
	}

	return "", ErrMaxRetriesForCollision
}

// generateKey generates a random key
// 1. calculate current epoch timestamp
// 2. generate a hash via sha256 from timestamp and the value
// 3. calculate base62 of the trunicated hash
func (s *service) generateKey(seed string, timestamp time.Time) string {
	epoch := timestamp.UnixNano()
	salt := strconv.FormatInt(epoch, 10)

	// TODO: use mobile-number or account-id instead of value
	hash := sha256.Sum256([]byte(seed + salt))
	shortHash := hash[:s.config.ShortURLLength]

	return encodeToBase62(shortHash)
}

// Base62 charset
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encodeToBase62(input []byte) string {
	num := new(big.Int).SetBytes(input)
	var result strings.Builder
	base := big.NewInt(62)
	mod := new(big.Int)

	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base, mod)
		result.WriteByte(base62Chars[mod.Int64()])
	}

	// reverse result
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

var (
	ErrShortenIDNotExists         = errors.New("ErrShortenIDNotExists")
	ErrRetreivingDataFromDatabase = errors.New("error retreiving data from database")
)

// Retrieve retrieves the key's value from the database
func (s *service) Retrieve(ctx context.Context, id string) (value entities.URL, err error) {
	defer func(start time.Time) {
		var status = metrics_pkg.StatusFailure
		if err == nil {
			s.metrics.Histogram.ObserveResponseTime(start, "retrieve")
			status = metrics_pkg.StatusSuccess
		}
		s.metrics.Counter.IncrementVector("retrieve", status)
	}(time.Now())

	urlString, err := s.redis.retrieve(ctx, id)
	if err == nil {
		return entities.URL(urlString), nil
	}
	// todo: just log the error

	urlString, err = s.postgres.retrieve(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIDNotExists) {
			return "", ErrShortenIDNotExists
		}
		return "", errors.Join(ErrRetreivingDataFromDatabase, err)
	}
	_ = s.redis.insert(ctx, id, urlString, s.config.CacheExpiration)

	return entities.URL(urlString), nil
}
