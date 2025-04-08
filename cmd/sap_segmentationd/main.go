package main

import (
	"log"
	"os"

	_ "go-test/docs/generated"
	"go-test/internal/api"
	"go-test/internal/logutil"
	"go-test/internal/repository"
	"go-test/internal/sap"
	"go-test/internal/storage"
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
	initLogger := log.New(os.Stdout, "INIT: ", log.LstdFlags)

	cfg := config.MustLoad(nil)

	logger, err := logutil.SetupLogger(cfg.Env, logDir, logFileName)
	if err != nil {
		initLogger.Fatalf("failed to setup logger: %v", err)
	}

	if err := logutil.CleanupOldLogs(logDir, cfg.Import.LogCleanupMaxAge, logger); err != nil {
		logger.Error("failed to cleanup old logs", "error", err.Error())
	}

	db, err := storage.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	segmentationRepo := repository.NewSegmentationRepository(db)

	sapClient := sap.NewClient(cfg, logger)

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

	server := api.NewServer(cfg, logger, sapClient, segmentationRepo)
	if err := server.Run(":" + cfg.App.Port); err != nil {
		logger.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
