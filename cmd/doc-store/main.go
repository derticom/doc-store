package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/derticom/doc-store/config"
	"github.com/derticom/doc-store/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	cfg := config.NewConfig()

	log, err := setupLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger: %+v", err))
	}

	go func() {
		if err = app.Run(ctx, cfg, log); err != nil {
			log.Error("critical service error", "error", err)
			stop()
			return
		}
	}()

	<-ctx.Done()

	log.Info("shutdown service ...")
}

func setupLogger(cfg *config.Config) (*slog.Logger, error) {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		),
	)

	return logger, nil
}
