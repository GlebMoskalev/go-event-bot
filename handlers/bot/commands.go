package bot

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
	"github.com/GlebMoskalev/go-event-bot/pkg/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CommandAccess = map[string]models.Role{
	"start":        models.RoleGuest,
	"schedule":     models.RoleStaff,
	"admin_panel":  models.RoleAdmin,
	"change_event": models.RoleAdmin,
	"add_staff":    models.RoleAdmin,
}

func (h *handler) Commands(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	cmd := update.Message.Command()

	msg := tgbotapi.NewMessage(chatID, "")

	requiredRole, exists := CommandAccess[cmd]
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
	case "schedule":
		msg = h.command.Schedule(ctx, msg)
	case "admin_panel":
		msg = h.adminCommand.Panel(ctx, msg)
	case "change_event":
		msg = h.adminCommand.ChangeEvent(ctx, msg)
	case "add_staff":
		msg = h.adminCommand.AddStaff(ctx, msg)
	}
	h.SendMessage(tgbot, msg)
}
