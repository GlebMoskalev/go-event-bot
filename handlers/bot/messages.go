package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update, state models.State) {
	log := logger.SetupLogger(h.log,
		"bot_handler", "Message",
		"chat_id", update.Message.Chat.ID,
		"message_id", update.Message.MessageID,
		"state", state,
	)

	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, update.Message.Text)

	if state != "" {
		log.Info("handling message with state")
		msg = h.message.State(ctx, msg, state)
		h.SendMessage(tgbot, msg)
	} else if update.Message.Contact != nil {
		log.Info("handling contact message", "phone_number", update.Message.Contact.PhoneNumber)
		var err error
		var role models.Role
		msg, role, err = h.message.Contact(ctx, msg, update.Message.Contact)
		menuCommands := commands.GetMenuCommands(role)
		_, err = tgbot.Request(tgbotapi.NewSetMyCommandsWithScope(
			tgbotapi.NewBotCommandScopeChat(chatID),
			menuCommands...,
		))
		if err != nil {
			h.log.Error("failed to set menu commands", "error", err)
		}
		h.SendMessage(tgbot, msg)
	} else {
		log.Info("handling plain text message", "text", update.Message.Text)
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
