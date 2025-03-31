package bot

import (
	"github.com/GlebMoskalev/go-event-bot/internal/bot/handlers"
	"github.com/GlebMoskalev/go-event-bot/internal/bot/middleware"
	"github.com/GlebMoskalev/go-event-bot/internal/config"
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
)

type Bot struct {
	log          *slog.Logger
	botService   services.BotService
	auth         *middleware.AuthMiddleware
	baseHandler  *handlers.BaseHandler
	adminHandler *handlers.AdminHandler
}

func NewBot(botService services.BotService, log *slog.Logger) *Bot {
	return &Bot{botService: botService, log: log}
}

func (b *Bot) Start(cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		b.log.Error("failed created bot", "error", err)
		os.Exit(1)
	}
	b.log.Info("initial bot")

	b.botService.SetBot(bot)
	b.auth = middleware.NewAuthMiddleWare(b.botService, b.log)
	b.adminHandler = handlers.NewAdminHandler(b.botService, b.log)
	b.baseHandler = handlers.NewBaseHandlers(b.botService, b.log)

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Bot.UpdateTimeout

	updates := bot.GetUpdatesChan(u)
	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Начальная команда",
		},
	)
	bot.Send(cmdCfg)

	for update := range updates {
		if !b.auth.CheckAuth(update) {
			continue
		}

		if update.Message.Contact != nil {
			b.botService.RequestContact(update)
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update)
			continue
		}

	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		b.baseHandler.HandleStart(update)
	case "schedule":
		b.baseHandler.HandleSchedule(update)
	case "admin":
		b.adminHandler.HandleAdmin(update)
	case "change_event":
		b.adminHandler.HandleChangeEvent(update)
	case "add_staff":
		b.adminHandler.HandleAddStaff(update)
	default:
		b.botService.SendMessage(update.Message.Chat.ID, "Неизвестная команда", nil)
	}
}
