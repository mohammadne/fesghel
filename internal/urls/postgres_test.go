package urls

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var urlColumns = []string{
	// "id",
	"url",
	// "created_at",
}

func TestPostgresInsert(t *testing.T) {
	var (
		sampleId  = ""
		sampleUrl = "https://sample.com"
	)

	t.Run("valid insert", func(t *testing.T) {
		timestamp := time.Now()

		mockDatabase.
			ExpectExec(regexp.QuoteMeta(queryInsert)).
			WithArgs(sampleId, sampleUrl, timestamp).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := postgresInstacne.insert(context.TODO(), sampleId, sampleUrl, timestamp)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if err := mockDatabase.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestPostgresRetrieve(t *testing.T) {
	var (
		sampleId  = ""
		sampleUrl = "https://sample.com"
	)

	t.Run("with empty result", func(t *testing.T) {
		mockDatabase.
			ExpectQuery(regexp.QuoteMeta(queryRetrieve)).
			WithArgs(sampleId).
			WillReturnRows(sqlmock.NewRows(urlColumns))

		_, err := postgresInstacne.retrieve(context.TODO(), sampleId)
		if !errors.Is(err, ErrIDNotExists) {
			t.Errorf("expect ErrIDNotExists error %v", err)
		}

		if err := mockDatabase.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with valid non-empty result", func(t *testing.T) {
		mockDatabase.
			ExpectQuery(regexp.QuoteMeta(queryRetrieve)).
			WithArgs(sampleId).
			WillReturnRows(sqlmock.NewRows(urlColumns).AddRow(sampleUrl))

		url, err := postgresInstacne.retrieve(context.TODO(), sampleId)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if url != sampleUrl {
			t.Error("invalid url has been returned")
		}

		if err := mockDatabase.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
