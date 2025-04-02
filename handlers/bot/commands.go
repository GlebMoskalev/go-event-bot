package bot

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Commands(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	cmd := update.Message.Command()

	msg := tgbotapi.NewMessage(chatID, "")

	requiredRole, exists := commands.CommandAccess[cmd]
	if !exists {
		msg.Text = messages.UnknownCommand()
		h.SendMessage(tgbot, msg)
		return
	}

	hasRole, err := h.user.HasRole(ctx, telegramID, requiredRole)
	if err != nil && errors.Is(err, apperrors.ErrNotFoundUser) {
		msg.Text = messages.Error()
		h.SendMessage(tgbot, msg)
		return
	}

	if !hasRole {
		msg.Text = messages.AccessDenied()
		h.SendMessage(tgbot, msg)
		return
	}

	switch cmd {
	case "start":
		msg = h.command.Start(ctx, msg, telegramID)
		menuCommands := commands.GetMenuCommands(models.RoleGuest)
		_, err := tgbot.Request(tgbotapi.NewSetMyCommandsWithScope(
			tgbotapi.NewBotCommandScopeChat(chatID),
			menuCommands...,
		))
		if err != nil {
			h.log.Error("failed to set menu commands", "error", err)
		}
	case "event":
		msg = h.command.Event(ctx, msg)
	case "admin_panel":
		msg = h.adminCommand.Panel(ctx, msg)
	case "change_event":
		msg = h.adminCommand.ChangeEvent(ctx, msg)
	case "add_staff":
		msg = h.adminCommand.AddStaff(ctx, msg)
	}
	h.SendMessage(tgbot, msg)
}
