package main

import (
	"github.com/GlebMoskalev/go-event-bot/internal/config"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.AppEnv)
	log.Info("Starting app")
}
