package user

import "context"

type UseCase interface {
	Register(ctx context.Context, adminToken, login, password string) error
	Authenticate(ctx context.Context, login, password string) (string, error)
	Logout(ctx context.Context, token string) error
}
