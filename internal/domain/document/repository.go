package document

import "context"

type Repository interface {
	List(ctx context.Context, userID string) ([]*Document, error)
	GetByID(ctx context.Context, id string) (*Document, error)
}
