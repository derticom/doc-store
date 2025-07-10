package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/derticom/doc-store/internal/domain/user"
	"github.com/derticom/doc-store/internal/usecase/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

type UseCase struct {
	repo       user.Repository
	auth       auth.SessionStore
	adminToken string
}

func NewUserUseCase(repo user.Repository, auth auth.SessionStore, adminToken string) *UseCase {
	return &UseCase{
		repo:       repo,
		auth:       auth,
		adminToken: adminToken,
	}
}

func (u *UseCase) Register(ctx context.Context, adminToken, login, password string) error {
	if adminToken != u.adminToken {
		return errors.New("unauthorized")
	}
	if err := validateLogin(login); err != nil {
		return fmt.Errorf("failed to validateLogin: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return fmt.Errorf("failed to validatePassword: %w", err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return u.repo.Create(ctx, &user.User{
		ID: uuid.NewString(), Login: login, PasswordHash: string(hash), CreatedAt: time.Now(),
	})
}

func (u *UseCase) Authenticate(ctx context.Context, login, password string) (string, error) {
	uData, err := u.repo.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("failed to repo.GetByLogin: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(uData.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("password incorrect")
	}

	token := uuid.NewString()
	if err := u.auth.Save(ctx, token, uData.ID, tokenTTL); err != nil {
		return "", err
	}
	return token, nil
}

func (u *UseCase) Logout(ctx context.Context, token string) error {
	return u.auth.Delete(ctx, token)
}
