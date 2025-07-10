package auth

import (
	"context"
	"time"
)

type SessionStore interface {
	Save(ctx context.Context, token, userID string, ttl time.Duration) error
	GetUserID(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
