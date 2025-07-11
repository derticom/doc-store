package document

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/derticom/doc-store/internal/domain/document"
	"github.com/derticom/doc-store/internal/domain/user"

	"github.com/google/uuid"
)

type FileStorage interface {
	Download(ctx context.Context, path string) ([]byte, error)
	Upload(ctx context.Context, path string, content []byte, mime string) error
	Delete(ctx context.Context, path string) error
}

type DocUseCase struct {
	docRepo  document.Repository
	userRepo user.Repository
	storage  FileStorage
	cache    document.Cache
}

func NewDocUseCase(
	docRepo document.Repository,
	userRepo user.Repository,
	storage FileStorage,
	cache document.Cache,
) document.UseCase {
	return &DocUseCase{
		docRepo:  docRepo,
		userRepo: userRepo,
		storage:  storage,
		cache:    cache,
	}
}

func (u *DocUseCase) List(ctx context.Context, userID string) ([]*document.Document, error) {
	return u.docRepo.List(ctx, userID)
}

func (u *DocUseCase) Get(ctx context.Context, id, userID string) (*document.Document, []byte, error) {
	doc, ok := u.cache.Get(id)
	if !ok {
		var err error
		doc, err = u.docRepo.GetByID(ctx, id)
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

func (u *DocUseCase) Upload(ctx context.Context, meta *document.Document, file []byte) error {
	meta.ID = uuid.NewString()
	meta.CreatedAt = time.Now()

	// проверяем grant логины
	var resolved []string
	for _, login := range meta.Grant {
		user, err := u.userRepo.GetByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("grant login not found: %s", login)
		}
		resolved = append(resolved, user.ID)
	}
	meta.Grant = resolved

	// сохранение файла (если есть)
	if meta.File && len(file) > 0 {
		path := fmt.Sprintf("%s/%s", meta.OwnerID, meta.ID)
		err := u.storage.Upload(ctx, path, file, meta.Mime)
		if err != nil {
			return fmt.Errorf("upload to storage: %w", err)
		}
	}

	// сохранение метаданных в БД
	if err := u.docRepo.Create(ctx, meta); err != nil {
		return fmt.Errorf("save to db: %w", err)
	}

	// инвалидация кеша
	u.cache.Invalidate(meta.ID)

	return nil
}

func (u *DocUseCase) Delete(ctx context.Context, docID, userID string) error {
	doc, err := u.docRepo.GetByID(ctx, docID)
	if err != nil {
		return err
	}

	// только владелец может удалять
	if doc.OwnerID != userID {
		return errors.New("access denied")
	}

	// удалить файл (если есть)
	if doc.File {
		path := fmt.Sprintf("%s/%s", doc.OwnerID, doc.ID)
		_ = u.storage.Delete(ctx, path)
	}

	// удалить из БД
	if err := u.docRepo.Delete(ctx, doc.ID); err != nil {
		return err
	}

	u.cache.Invalidate(doc.ID)

	return nil
}
