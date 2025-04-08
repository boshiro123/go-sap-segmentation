package logutil

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go-test/pkg/logger/slogpretty"
)

// SetupLogger настраивает логгер для вывода в консоль и файл
func SetupLogger(env, logDir, logFileName string) (*slog.Logger, error) {
	// Создаем директорию для логов, если она не существует
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Открываем файл для записи логов
	logFilePath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Создаем мультирайтер для записи в консоль и файл
	mw := io.MultiWriter(os.Stdout, logFile)

	var logger *slog.Logger

	switch env {
	case "local":
		// В локальном окружении используем "красивый" вывод
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := opts.NewPrettyHandler(mw)
		logger = slog.New(handler)
	case "dev":
		// В dev окружении используем JSON формат с уровнем Debug
		logger = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		// В production окружении используем JSON формат с уровнем Info
		logger = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return logger, nil
}
