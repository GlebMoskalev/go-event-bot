package main

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/configs"
	"github.com/GlebMoskalev/go-event-bot/handlers/bot"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/repositories/postgres"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.AppEnv)

	ctx := context.TODO()
	db, err := postgres.New(ctx, cfg.BotPostgres, cfg.StaffPostgres, log)
	if err != nil {
		panic(err)
	}

	defer func(db repositories.DB) {
		err := db.Close()
		if err != nil {
			log.Error("failed to close connection database", "error", err)
		}
	}(db)

	botInstance := bot.New(db, log)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := botInstance.Start(ctx, cfg, true)
		if err != nil {
			log.Error("bot failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-sigChan

	log.Info("Received an interrupt, Bot stopped...")
}
