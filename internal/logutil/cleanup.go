package logutil

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// CleanupOldLogs удаляет лог-файлы старше указанного возраста
func CleanupOldLogs(logDir string, maxAgeInDays int, logger *slog.Logger) error {
	// Создаем директорию для логов, если она не существует
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Получаем текущее время
	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -maxAgeInDays)

	logger.Info("starting log cleanup",
		"log_dir", logDir,
		"max_age_days", maxAgeInDays,
		"cutoff_time", cutoffTime.Format(time.RFC3339),
	)

	// Получаем список файлов в директории логов
	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var removedCount int

	// Проходим по всем файлам
	for _, file := range files {
		if file.IsDir() {
			continue // Пропускаем поддиректории
		}

		// Получаем полный путь к файлу
		filePath := filepath.Join(logDir, file.Name())

		// Получаем информацию о файле
		fileInfo, err := file.Info()
		if err != nil {
			logger.Error("failed to get file info", "file", filePath, "error", err.Error())
			continue
		}

		// Если файл старше указанного возраста, удаляем его
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
