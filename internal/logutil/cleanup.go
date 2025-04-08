package logutil

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func CleanupOldLogs(logDir string, maxAgeInDays int, logger *slog.Logger) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -maxAgeInDays)

	logger.Info("starting log cleanup",
		"log_dir", logDir,
		"max_age_days", maxAgeInDays,
		"cutoff_time", cutoffTime.Format(time.RFC3339),
	)

	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var removedCount int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(logDir, file.Name())

		fileInfo, err := file.Info()
		if err != nil {
			logger.Error("failed to get file info", "file", filePath, "error", err.Error())
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			logger.Info("removing old log file",
				"file", filePath,
				"mod_time", fileInfo.ModTime().Format(time.RFC3339),
			)

			if err := os.Remove(filePath); err != nil {
				logger.Error("failed to remove log file", "file", filePath, "error", err.Error())
			} else {
				removedCount++
			}
		}
	}

	logger.Info("log cleanup completed", "removed_files", removedCount)
	return nil
}
