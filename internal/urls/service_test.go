package urls

import (
	"context"
	"errors"
	"testing"

	"github.com/mohammadne/fesghel/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceShorten(t *testing.T) {
	var url = "https://example.com"

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

func TestServiceRetrieve(t *testing.T) {

}
