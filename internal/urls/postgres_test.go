package urls

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	postgres_pkg "github.com/mohammadne/fesghel/pkg/databases/postgres"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

func TestPostgres(t *testing.T) {
	var (
		mockDB sqlmock.Sqlmock
		p      *postgres
	)

	{ // initialization
		var err error
		var sqlDB *sql.DB

		sqlDB, mockDB, err = sqlmock.New()
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

		p = &postgres{instance: &postgres_pkg.Postgres{DB: sqlxDB, Vectors: &vectors}}
	}

	t.Run("insert", func(t *testing.T) {

	})

	t.Run("retrieve", func(t *testing.T) {

	})
}
