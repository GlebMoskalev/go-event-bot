package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update, state models.State) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, update.Message.Text)

	if state != "" {
		msg = h.message.State(ctx, msg, state)
		h.SendMessage(tgbot, msg)
	} else if update.Message.Contact != nil {
		var err error
		msg, err = h.message.Contact(ctx, msg, update.Message.Contact)
		if err == nil {
			menuCommands := commands.GetMenuCommands(models.RoleStaff)
			_, err := tgbot.Request(tgbotapi.NewSetMyCommandsWithScope(
				tgbotapi.NewBotCommandScopeChat(chatID),
				menuCommands...,
			))
			if err != nil {
				h.log.Error("failed to set menu commands", "error", err)
			}
		}
		_ = err
		h.SendMessage(tgbot, msg)
	} else {
		msg.Text = update.Message.Text
		h.SendMessage(tgbot, msg)
	}
}

func (h *handler) SendMessage(tgbot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	_, err := tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
