package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Callbacks(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	fmt.Println(update.CallbackQuery.Data)
}
