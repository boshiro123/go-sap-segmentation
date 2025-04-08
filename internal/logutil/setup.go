package logutil

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go-test/pkg/logger/slogpretty"
)

func SetupLogger(env, logDir, logFileName string) (*slog.Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	var logger *slog.Logger

	switch env {
	case "local":
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := opts.NewPrettyHandler(mw)
		logger = slog.New(handler)
	case "dev":
		logger = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		logger = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return logger, nil
}
