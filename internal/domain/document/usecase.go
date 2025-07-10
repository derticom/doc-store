package document

import "context"

type UseCase interface {
	List(ctx context.Context, userID string) ([]*Document, error)
	Get(ctx context.Context, id string, userID string) (*Document, []byte, error)
}
