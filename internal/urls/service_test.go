package urls

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadne/fesghel/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceShorten(t *testing.T) {
	var (
		url = "https://example.com"
	)

	t.Run("success", func(t *testing.T) {
		{ // prepare the mocks
			postgresMock.
				On("insert", mock.Anything, mock.Anything, url, mock.Anything).
				Return(nil).Once()

			redisMock.
				On("insert", mock.Anything, mock.Anything, url, serviceInstance.config.CacheExpiration).
				Return(nil).Once()
		}

		id, err := serviceInstance.Shorten(context.TODO(), entities.URL(url))
		assert.NoError(t, err)
		assert.NotNil(t, id)
		postgresMock.AssertExpectations(t)
		redisMock.AssertExpectations(t)
	})

	t.Run("check collision", func(t *testing.T) {
		t.Run("retry for one collision", func(t *testing.T) {
			{ // prepare the mocks
				postgresMock.
					On("insert", mock.Anything, mock.Anything, url, mock.Anything).
					Return(errUniqueConstraintViolated).Once()

				postgresMock.
					On("insert", mock.Anything, mock.Anything, url, mock.Anything).
					Return(nil).Once()

				redisMock.
					On("insert", mock.Anything, mock.Anything, url, serviceInstance.config.CacheExpiration).
					Return(nil).Once()
			}

			id, err := serviceInstance.Shorten(context.TODO(), entities.URL(url))
			assert.NoError(t, err)
			assert.NotNil(t, id)
			postgresMock.AssertExpectations(t)
			redisMock.AssertExpectations(t)
		})

		t.Run("max retry exceeded", func(t *testing.T) {
			{ // prepare the mocks
				for range serviceInstance.config.MaxRetriesOnCollision {
					postgresMock.
						On("insert", mock.Anything, mock.Anything, url, mock.Anything).
						Return(errUniqueConstraintViolated).Once()
				}
			}

			_, err := serviceInstance.Shorten(context.TODO(), entities.URL(url))
			if !errors.Is(err, errUniqueConstraintViolated) {
				t.Errorf("expect errUniqueConstraintViolated error %v", err)
			}
			postgresMock.AssertExpectations(t)
		})
	})

	t.Run("postgres error", func(t *testing.T) {
		{ // prepare the mocks
			postgresMock.
				On("insert", mock.Anything, mock.Anything, url, mock.Anything).
				Return(errInsertingURL).Once()
		}

		_, err := serviceInstance.Shorten(context.TODO(), entities.URL(url))
		if !errors.Is(err, ErrInsertingIntoPostgres) {
			t.Errorf("expect ErrInsertingIntoPostgres error %v", err)
		}
		postgresMock.AssertExpectations(t)
	})
}

func TestGenerateKey(t *testing.T) {
	key := serviceInstance.generateKey("anything", time.Now())

	expected := false
	for i := range 3 {
		if len(key) == serviceInstance.config.ShortURLLength+i {
			expected = true
		}
	}

	if !expected {
		t.Errorf("invalid key length %d", len(key))
	}
}

func TestServiceRetrieve(t *testing.T) {
	var (
		sampleURL = "id"
		sampleID  = "id"
	)

	t.Run("no cache (error) and no postgres", func(t *testing.T) {
		{ // prepare the mocks
			redisMock.
				On("retrieve", mock.Anything, sampleID).
				Return("", errIDNotFound).Once()

			postgresMock.
				On("retrieve", mock.Anything, sampleID).
				Return("", ErrIDNotExists).Once()
		}

		_, err := serviceInstance.Retrieve(context.TODO(), sampleID)
		if !errors.Is(err, ErrShortenIDNotExists) {
			t.Errorf("expect ErrShortenIDNotExists error %v", err)
		}
		postgresMock.AssertExpectations(t)
	})

	t.Run("no cache (error) and postgres error", func(t *testing.T) {
		{ // prepare the mocks
			redisMock.
				On("retrieve", mock.Anything, sampleID).
				Return("", errIDNotFound).Once()

			postgresMock.
				On("retrieve", mock.Anything, sampleID).
				Return("", ErrRetreivingValue).Once()
		}

		_, err := serviceInstance.Retrieve(context.TODO(), sampleID)
		if !errors.Is(err, ErrRetreivingDataFromDatabase) {
			t.Errorf("expect ErrRetreivingDataFromDatabase error %v", err)
		}
		postgresMock.AssertExpectations(t)
	})

	t.Run("success with cache", func(t *testing.T) {
		{ // prepare the mocks
			redisMock.
				On("retrieve", mock.Anything, sampleID).
				Return(sampleURL, nil).Once()
		}

		url, err := serviceInstance.Retrieve(context.TODO(), sampleID)
		assert.NoError(t, err)
		assert.Equal(t, sampleURL, string(url))
		postgresMock.AssertExpectations(t)
		redisMock.AssertExpectations(t)
	})

	t.Run("success on no cache (error)", func(t *testing.T) {
		{ // prepare the mocks
			redisMock.
				On("retrieve", mock.Anything, sampleID).
				Return("", errIDNotFound).Once()

			postgresMock.
				On("retrieve", mock.Anything, sampleID).
				Return(sampleURL, nil).Once()

			redisMock.
				On("insert", mock.Anything, mock.Anything, sampleURL, serviceInstance.config.CacheExpiration).
				Return(nil).Once()
		}

		url, err := serviceInstance.Retrieve(context.TODO(), sampleID)
		assert.NoError(t, err)
		assert.Equal(t, entities.URL(sampleURL), url)
		postgresMock.AssertExpectations(t)
		redisMock.AssertExpectations(t)
	})
}
