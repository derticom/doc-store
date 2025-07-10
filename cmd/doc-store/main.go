package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/derticom/doc-store/config"
	"github.com/derticom/doc-store/internal/app"
	"github.com/derticom/doc-store/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	log, err := logger.SetupLogger(cfg.LogLevel)
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
