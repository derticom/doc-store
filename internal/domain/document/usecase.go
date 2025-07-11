package document

import "context"

type UseCase interface {
	List(ctx context.Context, userID string) ([]*Document, error)
	Get(ctx context.Context, id string, userID string) (*Document, []byte, error)
	Upload(ctx context.Context, doc *Document, file []byte) error
	Delete(ctx context.Context, id string, userID string) error
}
