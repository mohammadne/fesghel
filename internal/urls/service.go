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
	"github.com/mohammadne/fesghel/pkg/databases/postgres"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

type Service interface {
	// Shorten shortenes the url by giving a url string, then returns the shortened id
	Shorten(ctx context.Context, url entities.URL) (string, error)

	// Retrieve returns the actual url by giving url's shortened id
	Retrieve(ctx context.Context, id string) (entities.URL, error)
}

type service struct {
	config *Config
	logger *zap.Logger
	m      *metrics
	p      *postgres.Postgres
	r      Redis
}

func Initialize(cfg *Config, l *zap.Logger) (Service, error) {
	var svc = service{config: cfg, logger: l}

	metrics, err := newMetrics()
	if err != nil {
		l.Panic("error loading Postgres instance", zap.Error(err))
	}
	svc.m = metrics

	postgresInstance, err := postgres.Open(cfg.Postgres, entities.Namespace, entities.System)
	if err != nil {
		l.Panic("error loading Postgres instance", zap.Error(err))
	}
	svc.p = postgresInstance

	redis, err := newRedis(cfg.Redis)
	if err != nil {
		l.Panic("error initializing Redis cache", zap.Error(err))
	}
	svc.r = redis

	return &svc, nil
}

var (
	ErrGenerateKey            = errors.New("error generate key")
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
			s.m.Histogram.ObserveResponseTime(start, "shorten")
			status = metrics_pkg.StatusSuccess
		}
		s.m.Counter.IncrementVector("shorten", status)
	}(time.Now())

	for attempt := 1; attempt <= s.config.MaxRetriesOnCollision; attempt++ {
		timestamp := time.Now()

		key, err = s.generateKey(string(url), timestamp)
		if err != nil {
			return "", errors.Join(ErrGenerateKey, err)
		}

		err = s.insertIntoPostgres(ctx, key, string(url), timestamp)
		if err == nil {
			_ = s.r.insert(ctx, key, string(url), s.config.CacheExpiration)
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
func (s *service) generateKey(url string, timestamp time.Time) (string, error) {
	epoch := timestamp.UnixNano()
	salt := strconv.FormatInt(epoch, 10)

	// TODO: use mobile-number or account-id instead of value
	hash := sha256.Sum256([]byte(url + salt))
	shortHash := hash[:s.config.ShortURLLength]

	return encodeToBase62(shortHash), nil
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
	ErrRetreivingDataFromDatabase = errors.New("error retreiving data from database")
)

// Retrieve retrieves the key's value from the database
func (s *service) Retrieve(ctx context.Context, id string) (value entities.URL, err error) {
	defer func(start time.Time) {
		var status = metrics_pkg.StatusFailure
		if err == nil {
			s.m.Histogram.ObserveResponseTime(start, "retrieve")
			status = metrics_pkg.StatusSuccess
		}
		s.m.Counter.IncrementVector("retrieve", status)
	}(time.Now())

	urlString, err := s.r.retrieve(ctx, id)
	if err == nil {
		return entities.URL(urlString), nil
	}
	// todo: just log the error

	urlString, err = s.retrieveFromOracle(ctx, id)
	if err != nil {
		return "", errors.Join(ErrRetreivingDataFromDatabase, err)
	}
	_ = s.r.insert(ctx, id, urlString, s.config.CacheExpiration)

	return entities.URL(urlString), nil
}
