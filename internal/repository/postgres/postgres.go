package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context, dsn string, migrationsPath string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to sql.Open: %w", err)
	}

	// Retry Ping
	const maxAttempts = 5
	for i := 1; i <= maxAttempts; i++ {
		err = db.PingContext(ctx)
		if err == nil {
			break
		}

		if i < maxAttempts {
			time.Sleep(time.Second * time.Duration(i))
			continue
		}

		db.Close()
		return nil, fmt.Errorf("failed to PingContext after %d attempts: %w", maxAttempts, err)
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
