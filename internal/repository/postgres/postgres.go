package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context, dsn string, migrationsPath string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to sql.Open: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to PingContext: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to goose.SetDialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to goose.Up: %w", err)
	}

	return db, nil
}
