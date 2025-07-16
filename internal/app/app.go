package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/derticom/doc-store/config"
	"github.com/derticom/doc-store/internal/repository/minio"
	"github.com/derticom/doc-store/internal/repository/postgres"
	"github.com/derticom/doc-store/internal/repository/redis"
	"github.com/derticom/doc-store/internal/server"
	"github.com/derticom/doc-store/internal/usecase/document"
	"github.com/derticom/doc-store/internal/usecase/user"
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	pgDB, err := postgres.Connect(ctx, cfg.PostgresURL, "./migrations")
	if err != nil {
		return fmt.Errorf("failed to postgres.Connect: %w", err)
	}

	docRepo := postgres.NewDocRepo(pgDB)
	userRepo := postgres.NewUserRepo(pgDB)

	storage, err := minio.New(
		cfg.Minio.Address,
		cfg.Minio.AccessKey,
		cfg.Minio.SecretKey,
		cfg.Minio.Bucket,
		cfg.Minio.UseSSL,
	)
	if err != nil {
		return fmt.Errorf("failed to minio.New: %w", err)
	}

	cache, err := redis.NewCache(ctx, cfg.RedisURL, 0)
	if err != nil {
		return fmt.Errorf("failed to redis.NewCache: %w", err)

	}

	sessions, err := redis.NewSessions(ctx, cfg.RedisURL, 1)
	if err != nil {
		return fmt.Errorf("failed to redis.NewSessions: %w", err)
	}

	docUseCase := document.NewDocUseCase(docRepo, userRepo, storage, cache)

	userUseCase := user.NewUserUseCase(userRepo, sessions, cfg.AdminToken)

	srv := server.New(cfg.Server.Address, log, docUseCase, userUseCase, sessions)

	err = srv.Run(ctx)
	if err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}
