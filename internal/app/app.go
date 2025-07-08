package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/derticom/doc-store/config"
	"github.com/derticom/doc-store/internal/server"
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	srv := server.New(cfg.Server.Address, log)

	err := srv.Run()
	if err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}
