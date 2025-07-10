package document

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/derticom/doc-store/internal/domain/document"
)

type FileStorage interface {
	Download(ctx context.Context, path string) ([]byte, error)
	Upload(ctx context.Context, path string, content []byte, mime string) error
	Delete(ctx context.Context, path string) error
}

type DocUseCase struct {
	repo    document.Repository
	storage FileStorage
	cache   document.Cache
}

func NewDocUseCase(repo document.Repository, storage FileStorage, cache document.Cache) document.UseCase {
	return &DocUseCase{
		repo:    repo,
		storage: storage,
		cache:   cache,
	}
}

func (u *DocUseCase) List(ctx context.Context, userID string) ([]*document.Document, error) {
	return u.repo.List(ctx, userID)
}

func (u *DocUseCase) Get(ctx context.Context, id, userID string) (*document.Document, []byte, error) {
	doc, ok := u.cache.Get(id)
	if !ok {
		var err error
		doc, err = u.repo.GetByID(ctx, id)
		if err != nil {
			return nil, nil, err
		}
		u.cache.Set(id, doc)
	}

	if !doc.Public && doc.OwnerID != userID && slices.Contains(doc.Grant, userID) { // проверка доступа
		return nil, nil, errors.New("access denied")
	}

	var data []byte
	var err error
	if doc.File {
		path := fmt.Sprintf("%s/%s", doc.OwnerID, doc.ID)
		data, err = u.storage.Download(ctx, path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to storage.Download: %w", err)
		}
	}

	u.cache.Set(id, doc)
	return doc, data, nil
}
