package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (h *handler) Callbacks(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log := logger.SetupLogger(h.log,
		"bot_handler", "Callbacks",
		"chat_id", update.CallbackQuery.Message.Chat.ID,
		"query_id", update.CallbackQuery.ID,
		"callback_data", update.CallbackQuery.Data,
	)
	log.Info("processing pagination request")

	split := strings.Split(update.CallbackQuery.Data, ":")
	if split[0] == models.PaginationPrefix && split[1] == models.EventContext {
		log.Info("handling pagination event callback")
		msg := h.callback.PagerEvent(ctx, update.CallbackQuery, split[2:]...)
		_, err := tgbot.Request(msg)
		if err != nil {
			log.Error("failed to send pagination message", "error", err)
			return
		}
	} else if split[0] == "event" && split[1] == "all" {
		log.Info("handling event all callback")
		msg := h.callback.EventAll(ctx, update.CallbackQuery)
		_, err := tgbot.Request(msg)
		if err != nil {
			log.Error("failed to send event all message", "error", err)
			return
		}
	} else if split[0] == "back" && split[1] == "event" {
		log.Info("handling back to events callback")
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
			"Мероприятия:", keyboards.EventInline(),
		)
		_, err := tgbot.Request(msg)
		if err != nil {
			log.Error("failed to send back to events message", "error", err)
			return
		}
	} else {
		log.Warn("unknown callback data received")
		msg := tgbotapi.NewCallback(update.CallbackQuery.ID, "Функция не реализована")
		_, err := tgbot.Request(msg)
		if err != nil {
			log.Error("failed to send unknown callback", "error", err)
			return
		}
	}
}
