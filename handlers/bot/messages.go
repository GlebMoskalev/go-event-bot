package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")

	if update.Message.Contact != nil {
		msg = h.message.Contact(ctx, msg, update.Message.Contact)
	} else {
		textError, err := h.AuthorizeUser(ctx, telegramID, false)
		if err != nil {
			msg.Text = textError
			goto sendMessage
		}
		msg.Text = "Привет!"
	}

sendMessage:
	_, err := tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
