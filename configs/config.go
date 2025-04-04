package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type (
	App struct {
		AppEnv        string `env:"APP_ENV,notEmpty"`
		Bot           Bot
		BotPostgres   BotPostgres
		StaffPostgres StaffPostgres
	}
	BotPostgres struct {
		Host     string `env:"POSTGRES_BOT_HOST,notEmpty"`
		Password string `env:"POSTGRES_BOT_PASSWORD,notEmpty"`
		Name     string `env:"POSTGRES_BOT_DATABASE_NAME,notEmpty"`
		Port     string `env:"POSTGRES_BOT_PORT,notEmpty"`
		User     string `env:"POSTGRES_BOT_USER,notEmpty"`
	}
	StaffPostgres struct {
		Host     string `env:"POSTGRES_STAFF_HOST,notEmpty"`
		Password string `env:"POSTGRES_STAFF_PASSWORD,notEmpty"`
		Name     string `env:"POSTGRES_STAFF_DATABASE_NAME,notEmpty"`
		Port     string `env:"POSTGRES_STAFF_PORT,notEmpty"`
		User     string `env:"POSTGRES_STAFF_USER,notEmpty"`
	}

	Bot struct {
		Token         string `env:"BOT_TOKEN,notEmpty"`
		UpdateTimeout int    `env:"UPDATE_TIMEOUT,notEmpty"`
	}
)

func init() {
	if err := godotenv.Load("configs/.env"); err != nil {
		slog.Error("No .env file not found")
		os.Exit(1)
	}
}

func LoadConfig() (App, error) {
	cfg := App{}
	err := env.Parse(&cfg)
	return cfg, err
}
