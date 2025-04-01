package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	contact := update.Message.Contact

	existsUser, err := h.user.ExistsUserByTelegramID(ctx, update.Message.From.ID)

	if err != nil {
		h.log.Error("failed to check user by telegram_id", "error", err)
	} else if !existsUser && contact == nil {
		msg.Text = "Нужно выполнить команду /start"
	} else if contact != nil {
		msg = h.message.Contact(ctx, msg, update.Message.Contact)
	} else {
		msg.Text = "Привет"
	}

	_, err = tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
