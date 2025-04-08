package main

import (
	"log"
	"os"

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

	// Получаем данные о сегментации из SAP API
	logger.Info("starting import from SAP API")
	segments, err := sapClient.FetchSegmentation()
	if err != nil {
		logger.Error("failed to fetch segmentation data", "error", err.Error())
		os.Exit(1)
	}

	// Если данных нет, завершаем работу
	if len(segments) == 0 {
		logger.Info("no segmentation data to import")
		return
	}

	// Сохраняем данные в базу
	logger.Info("importing segmentation data to database", "count", len(segments))
	if err := segmentationRepo.InsertOrUpdate(segments); err != nil {
		logger.Error("failed to import segmentation data", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("import completed successfully", "total_imported", len(segments))
}
