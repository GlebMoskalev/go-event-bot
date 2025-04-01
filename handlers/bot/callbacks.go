package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (h *handler) Callbacks(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	split := strings.Split(update.CallbackQuery.Data, ":")
	if split[0] == "pager" && split[1] == "schedule" {
		msg := h.callback.PagerSchedule(ctx, update.CallbackQuery, split[2:]...)
		tgbot.Send(msg)
	}
}
