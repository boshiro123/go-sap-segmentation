package config

import (
	"log/slog"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит все настройки приложения
type Config struct {
	Env      string        `envconfig:"ENV" default:"local"`
	TokenTTL time.Duration `envconfig:"TOKEN_TTL" default:"1h"`

	DB struct {
		Host     string `envconfig:"DB_HOST" default:"127.0.0.1"`
		Port     string `envconfig:"DB_PORT" default:"5432"`
		Name     string `envconfig:"DB_NAME" default:"mesh_group"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	}

	Connection struct {
		URI          string        `envconfig:"CONN_URI" default:"http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation"`
		AuthLoginPwd string        `envconfig:"CONN_AUTH_LOGIN_PWD" default:"4Dfddf5:jKlljHGH"`
		UserAgent    string        `envconfig:"CONN_USER_AGENT" default:"spacecount-test"`
		Timeout      time.Duration `envconfig:"CONN_TIMEOUT" default:"5s"`
		Interval     time.Duration `envconfig:"CONN_INTERVAL" default:"1500ms"`
	}

	Import struct {
		BatchSize        int  `envconfig:"IMPORT_BATCH_SIZE" default:"50"`
		LogCleanupMaxAge int  `envconfig:"LOG_CLEANUP_MAX_AGE" default:"7"`
		UseTestData      bool `envconfig:"USE_TEST_DATA" default:"true"`
	}

	App struct {
		Port string `envconfig:"APP_PORT" default:"8080"`
	}
}

func MustLoad(logger *slog.Logger) *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("failed to load config", "error", err.Error())
		panic("failed to load config: " + err.Error())
	}

	return &cfg
}
