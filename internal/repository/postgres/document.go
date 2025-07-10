package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/derticom/doc-store/internal/domain/document"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type Repo struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Repo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to sql.Open: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to db.PingContext: %w", err)
	}

	return &Repo{db: db}, nil
}

func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return fmt.Errorf("failed to Close: %w", err)
	}

	return nil
}

func (r *Repo) Migrate(migrate string) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to goose.SetDialect: %w", err)
	}

	if err := goose.Up(r.db, migrate); err != nil {
		return fmt.Errorf("failed to goose.Up: %w", err)
	}

	return nil
}

func (r *Repo) List(ctx context.Context, userID string) ([]*document.Document, error) {
	const query = `
SELECT id, name, mime, file, public, created_at, owner_id
FROM documents
WHERE owner_id = $1 OR public = true OR $1 = ANY(grant_ids)
ORDER BY created_at DESC
`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to db.QueryContext: %w", err)
	}
	defer rows.Close()

	var docs []*document.Document
	for rows.Next() {
		var d document.Document
		if err := rows.Scan(
			&d.ID, &d.Name, &d.Mime, &d.File, &d.Public, &d.CreatedAt, &d.OwnerID,
		); err != nil {
			return nil, fmt.Errorf("failed to rows.Scan: %w", err)
		}
		docs = append(docs, &d)
	}
	return docs, nil
}

func (r *Repo) GetByID(ctx context.Context, documentID string) (*document.Document, error) {
	const query = `
SELECT id, name, mime, file, public, created_at, owner_id, grant_ids
FROM documents
WHERE id = $1
`
	var d document.Document
	var grantIDs pq.StringArray
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&d.ID,
		&d.Name,
		&d.Mime,
		&d.File,
		&d.Public,
		&d.CreatedAt,
		&d.OwnerID,
		&grantIDs,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}

	d.Grant = grantIDs

	return &d, nil
}
