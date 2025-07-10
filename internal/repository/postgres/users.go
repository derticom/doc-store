package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/derticom/doc-store/internal/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	const query = `INSERT INTO users (id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Login, u.PasswordHash, u.CreatedAt)
	return err
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	const query = `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := r.db.QueryRowContext(ctx, query, login)

	var u user.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	const query = `SELECT id, login, password_hash, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return &u, err
}
