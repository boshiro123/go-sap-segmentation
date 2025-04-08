package main

import (
	"log"
	"os"

	_ "go-test/docs/generated" // Импорт сгенерированной Swagger документации
	"go-test/internal/api"
	"go-test/internal/logutil"
	"go-test/internal/sap"
	"go-test/internal/storage"
	"go-test/model"
	"go-test/pkg/config"
)

const (
	logDir      = "log"
	logFileName = "segmentation_import.log"
)

// @title SAP Segmentation API
// @version 1.0
// @description API для импорта и доступа к данным сегментации из SAP

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Настраиваем временный логгер для процесса инициализации
	initLogger := log.New(os.Stdout, "INIT: ", log.LstdFlags)

	// Загружаем конфигурацию
	cfg := config.MustLoad(nil) // Временно передаем nil, так как основной логгер еще не создан

	// Настраиваем основной логгер
	logger, err := logutil.SetupLogger(cfg.Env, logDir, logFileName)
	if err != nil {
		initLogger.Fatalf("failed to setup logger: %v", err)
	}

	// Очищаем старые логи
	if err := logutil.CleanupOldLogs(logDir, cfg.Import.LogCleanupMaxAge, logger); err != nil {
		logger.Error("failed to cleanup old logs", "error", err.Error())
	}

	// Устанавливаем соединение с базой данных
	db, err := storage.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Создаем репозиторий для работы с сегментацией
	segmentationRepo := model.NewSegmentationRepository(db)

	// Создаем клиент для работы с SAP API
	sapClient := sap.NewClient(cfg, logger)

	// Запускаем импорт в фоновом режиме при старте, если установлен флаг
	runImportOnStart := os.Getenv("RUN_IMPORT_ON_START") == "true"
	if runImportOnStart {
		go func() {
			logger.Info("starting initial import from SAP API")
			segments, err := sapClient.FetchSegmentation()
			if err != nil {
				logger.Error("failed to fetch segmentation data", "error", err.Error())
				return
			}

			if len(segments) == 0 {
				logger.Info("no segmentation data to import")
				return
			}

			logger.Info("importing segmentation data to database", "count", len(segments))
			if err := segmentationRepo.InsertOrUpdate(segments); err != nil {
				logger.Error("failed to import segmentation data", "error", err.Error())
				return
			}

			logger.Info("initial import completed successfully", "total_imported", len(segments))
		}()
	}

	// Создаем и запускаем HTTP сервер
	server := api.NewServer(cfg, logger, sapClient, segmentationRepo)
	if err := server.Run(":" + cfg.App.Port); err != nil {
		logger.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
