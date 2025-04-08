package storage

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"go-test/pkg/config"
)

func NewPostgresDB(cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	logger.Info("connecting to database",
		"host", cfg.DB.Host,
		"port", cfg.DB.Port,
		"dbname", cfg.DB.Name,
		"user", cfg.DB.User,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("successfully connected to database")

	return db, nil
}
