package bot

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Commands(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log := logger.SetupLogger(h.log,
		"bot_handler", "Commands",
		"chat_id", update.Message.Chat.ID,
		"telegram_id", update.Message.From.ID,
		"command", update.Message.Command(),
	)

	telegramID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	cmd := update.Message.Command()

	msg := tgbotapi.NewMessage(chatID, "")

	requiredRole, exists := commands.CommandAccess[cmd]
	if !exists {
		log.Warn("unknown command received")
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
		log.Info("access denied for user")
		msg.Text = messages.AccessDenied()
		h.SendMessage(tgbot, msg)
		return
	}

	switch cmd {
	case "start":
		log.Info("handling start command")
		msg, err = h.command.Start(ctx, msg, telegramID)
		h.log.Info("start", "err", err)
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
	case "event":
		log.Info("handling event command")
		msg = h.command.Event(ctx, msg)
	case "admin_panel":
		log.Info("handling admin panel command")
		msg = h.adminCommand.Panel(ctx, msg)
	case "change_event":
		log.Info("handling change event command")
		msg = h.adminCommand.ChangeEvent(ctx, msg)
	case "add_staff":
		log.Info("handling add staff command")
		msg = h.adminCommand.AddStaff(ctx, msg)
	}
	h.SendMessage(tgbot, msg)
}
