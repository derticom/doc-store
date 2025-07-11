package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/derticom/doc-store/internal/domain/document"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type DocRepo struct {
	db *sql.DB
}

func NewDocRepo(db *sql.DB) *DocRepo {
	return &DocRepo{db: db}
}

func (r *DocRepo) List(ctx context.Context, userID string) ([]*document.Document, error) {
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

func (r *DocRepo) GetByID(ctx context.Context, documentID string) (*document.Document, error) {
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

func (r *DocRepo) Create(ctx context.Context, doc *document.Document) error {
	const query = `
INSERT INTO documents (id, name, mime, file, public, created_at, owner_id, grant_ids, json_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.Name, doc.Mime, doc.File, doc.Public,
		doc.CreatedAt, doc.OwnerID, pq.Array(doc.Grant), doc.JSONData)
	return err
}

func (r *DocRepo) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM documents WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
