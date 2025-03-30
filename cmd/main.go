package main

import (
	"database/sql"
	bot2 "github.com/GlebMoskalev/go-event-bot/internal/bot"
	"github.com/GlebMoskalev/go-event-bot/internal/config"
	"github.com/GlebMoskalev/go-event-bot/internal/database"
	"github.com/GlebMoskalev/go-event-bot/internal/repository"
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.AppEnv)

	tgBotDb, err := database.InitDatabase(cfg.TgBotDb.Host, cfg.TgBotDb.Port, cfg.TgBotDb.User, cfg.TgBotDb.Name, cfg.TgBotDb.Password)
	defer func(tgBotDb *sql.DB) {
		err := tgBotDb.Close()
		if err != nil {
			log.Error("failed close tgBot database", "err", err.Error())
		}
	}(tgBotDb)

	if err != nil {
		log.Error("failed initial tgBot database", "err", err.Error())
		os.Exit(1)
	}
	log.Info("initial tgBot database")

	staffDb, err := database.InitDatabase(cfg.StaffDB.Host, cfg.StaffDB.Port, cfg.StaffDB.User, cfg.StaffDB.Name, cfg.StaffDB.Password)
	defer func(tgStaffDb *sql.DB) {
		err := tgStaffDb.Close()
		if err != nil {
			log.Error("failed close staff database", "err", err.Error())
		}
	}(staffDb)

	if err != nil {
		log.Error("failed initial staff database", "err", err.Error())
		os.Exit(1)
	}
	log.Info("initial staff database")

	staffRepo := repository.NewStaffRepository(staffDb, log)
	staffService := services.NewStaffService(staffRepo, log)

	userRepo := repository.NewUserRepository(tgBotDb, log)
	userService := services.NewUserService(userRepo, log)

	botService := services.NewBotService(staffService, userService, log)
	bot := bot2.NewBot(botService, log)
	bot.Start(cfg)
}
