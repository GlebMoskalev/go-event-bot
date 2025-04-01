package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var AdminCommands = map[string]bool{
	"admin_panel":  true,
	"change_event": true,
	"add_staff":    true,
}

func (h *handler) Commands(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	cmd := update.Message.Command()

	msg := tgbotapi.NewMessage(chatID, "")

	if cmd == "start" {
		msg = h.command.Start(ctx, msg, telegramID)
	} else {
		requiresAdmin := AdminCommands[cmd]

		textError, err := h.AuthorizeUser(ctx, telegramID, requiresAdmin)
		if err != nil {
			msg.Text = textError
			goto sendMessage
		}

		switch cmd {
		case "schedule":
			msg = h.command.Schedule(ctx, msg)
		case "admin_panel":
			msg = h.adminCommand.Panel(ctx, msg)
		case "change_event":
			msg = h.adminCommand.ChangeEvent(ctx, msg)
		case "add_staff":
			msg = h.adminCommand.AddStaff(ctx, msg)
		default:
			msg.Text = "Неизвестная команда"
		}
	}

sendMessage:
	_, err := tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
