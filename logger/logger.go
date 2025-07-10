package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func SetupLogger(logLevel string) (*slog.Logger, error) {
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", logLevel)
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
