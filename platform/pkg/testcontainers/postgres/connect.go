package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pkg/errors"
)

func connectPostgresDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Errorf("failed to open postgres connection: %v", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, errors.Errorf("failed to ping postgres: %v", err)
	}

	return db, nil
}
