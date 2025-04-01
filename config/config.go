package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type (
	App struct {
		AppEnv   string `env:"APP_ENV,notEmpty"`
		Bot      Bot
		Postgres Postgres
	}
	Postgres struct {
		Host     string `env:"POSTGRES_HOST,notEmpty"`
		Password string `env:"POSTGRES_PASSWORD,notEmpty"`
		Name     string `env:"POSTGRES_DATABASE_NAME,notEmpty"`
		Port     string `env:"POSTGRES_PORT,notEmpty"`
		User     string `env:"POSTGRES_USER,notEmpty"`
	}

	Bot struct {
		Token         string `env:"BOT_TOKEN,notEmpty"`
		UpdateTimeout int    `env:"UPDATE_TIMEOUT,notEmpty"`
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file not found")
		os.Exit(1)
	}
}

func LoadConfig() (App, error) {
	cfg := App{}
	err := env.Parse(&cfg)
	return cfg, err
}
