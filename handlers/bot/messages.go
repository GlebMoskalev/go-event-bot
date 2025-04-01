package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")

	if update.Message.Contact != nil {
		var err error
		msg, err = h.message.Contact(ctx, msg, update.Message.Contact)
		if err == nil {
			h.log.Info("menuCommands")
			menuCommands := h.command.GetMenuCommands(models.RoleStaff)
			_, err := tgbot.Request(tgbotapi.NewSetMyCommandsWithScope(
				tgbotapi.NewBotCommandScopeChat(chatID),
				menuCommands...,
			))
			if err != nil {
				h.log.Error("failed to set menu commands", "error", err)
			}
		}
		h.SendMessage(tgbot, msg)
		return
	} else {
		msg.Text = update.Message.Text
		h.SendMessage(tgbot, msg)
		return
	}
}

func (h *handler) SendMessage(tgbot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	_, err := tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
