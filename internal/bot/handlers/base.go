package handlers

import (
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type BaseHandler struct {
	botService services.BotService
	log        *slog.Logger
}

func NewBaseHandlers(botService services.BotService, log *slog.Logger) *BaseHandler {
	return &BaseHandler{botService: botService, log: log}
}

func (h *BaseHandler) HandleStart(update tgbotapi.Update) {
	h.botService.Start(update)
}

func (h *BaseHandler) HandleSchedule(update tgbotapi.Update) {
	h.botService.SendMessage(update.Message.Chat.ID, "Расписнаие мероприятий...", nil)
}
