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
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	documents, err := postgres.New(ctx, cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to postgres.New: %w", err)
	}
	defer documents.Close()

	err = documents.Migrate("./migrations")
	if err != nil {
		return fmt.Errorf("failed to documents.Migrate: %w", err)
	}

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

	cache, err := redis.New(ctx, cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("failed to redis.New: %w", err)

	}

	docUseCase := document.NewDocUseCase(documents, storage, cache)

	srv := server.New(cfg.Server.Address, log, docUseCase)

	err = srv.Run()
	if err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}
