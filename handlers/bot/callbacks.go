package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (h *handler) Callbacks(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	split := strings.Split(update.CallbackQuery.Data, ":")
	if split[0] == models.PaginationPrefix && split[1] == models.ScheduleContext {
		msg := h.callback.PagerSchedule(ctx, update.CallbackQuery, split[2:]...)
		tgbot.Send(msg)
	} else if split[0] == "schedule" && split[1] == "all" {
		msg := h.callback.ScheduleAll(ctx, update.CallbackQuery)
		tgbot.Send(msg)
	}
}
