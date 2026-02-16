package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	} `yaml:"http_server"`
	Storage struct {
		Postgres struct {
			DSN            string
			MaxOpenConns   int           `yaml:"max_open_conns"`
			MaxIdleConns   int           `yaml:"max_idle_conns"`
			MaxIdleTime    time.Duration `yaml:"max_idle_time"`
			ConnAttempts   int           `yaml:"conn_attempts"`
			BaseRetryDelay time.Duration `yaml:"base_retry_delay"`
			MaxRetryDelay  time.Duration `yaml:"max_retry_delay"`
		} `yaml:"postgres"`
	} `yaml:"storage"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("Load config path is failed")
	}

	//panic("Load config failed")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file %s does not exists", configPath))

	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("Parse config is failed: %s", err.Error()))
	}

	cfg.Storage.Postgres.DSN = os.Getenv("DSN")
	if cfg.Storage.Postgres.DSN == "" {
		panic("Load DSN is failed")
	}

	return &cfg
}
