package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	AppEnv string `env:"APP_ENV,notEmpty"`
	Bot    struct {
		Token string `env:"BOT_TOKEN,notEmpty"`
	}
	TgBotDb struct {
		Host     string `env:"POSTGRES_TG_BOT_HOST,notEmpty"`
		User     string `env:"POSTGRES_TG_BOT_USER,notEmpty"`
		Password string `env:"POSTGRES_TG_BOT_PASSWORD,notEmpty"`
		Name     string `env:"POSTGRES_TG_BOT_DATABASE_NAME,notEmpty"`
		Port     string `env:"POSTGRES_TG_BOT_PORT,notEmpty"`
	}
	StaffDB struct {
		Host     string `env:"POSTGRES_STAFF_HOST,notEmpty"`
		User     string `env:"POSTGRES_STAFF_USER,notEmpty"`
		Password string `env:"POSTGRES_STAFF_PASSWORD,notEmpty"`
		Name     string `env:"POSTGRES_STAFF_DATABASE_NAME,notEmpty"`
		Port     string `env:"POSTGRES_STAFF_PORT,notEmpty"`
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file not found")
		os.Exit(1)
	}
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
