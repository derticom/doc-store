package user

import "context"

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
