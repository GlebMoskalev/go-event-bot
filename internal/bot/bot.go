package bot

import (
	"github.com/GlebMoskalev/go-event-bot/internal/config"
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
)

type Bot struct {
	log        *slog.Logger
	botService services.BotService
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

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Bot.UpdateTimeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		mes := update.Message
		switch {
		case mes.IsCommand():
			b.Commands(bot, update)
		case mes.Contact != nil:
			b.botService.RequestContact(bot, update)
		default:
			b.Messages(bot, update)
		}
	}
}
