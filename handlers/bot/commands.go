package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Commands(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	cmd := update.Message.Command()

	existsUser, existsErr := h.user.ExistsUserByTelegramID(ctx, update.Message.From.ID)
	isAdmin, adminErr := h.user.IsAdmin(ctx, update.Message.From.ID)

	if existsErr != nil {
		h.log.Error("failed check user by telegram_id", "error", existsErr)
		msg.Text = "Произошла ошибка!"
	} else if adminErr != nil {
		h.log.Error("failed check user by telegram_id", "error", existsErr)
		msg.Text = "Произошла ошибка!"
	} else if !existsUser && cmd != "start" {
		msg.Text = "Нужно выполнить команду /start"
	} else if cmd == "start" {
		msg = h.command.Start(ctx, msg, update.Message.From.ID)
	} else if cmd == "schedule" {
		msg = h.command.Schedule(ctx, msg)
	} else if isAdmin && cmd == "admin_panel" {
		msg = h.command.AdminPanel(ctx, msg)
	} else if isAdmin && cmd == "change_event" {
		msg = h.command.AdminChangeEvent(ctx, msg)
	} else if isAdmin && cmd == "add_staff" {
		msg = h.command.AdminAddStaff(ctx, msg)
	} else {
		msg.Text = "Неизвестная команда"
	}

	_, err := tgbot.Send(msg)
	if err != nil {
		h.log.Error("failed to send message: %v", err)
	}
}
