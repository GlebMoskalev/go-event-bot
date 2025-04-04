package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *handler) Message(ctx context.Context, tgbot *tgbotapi.BotAPI, update tgbotapi.Update, state models.State) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	h.log.Error("mess", "state", state)
	if state != "" {
		switch state {
		case models.StateStaffRegisterFullName:
			err := h.state.RegisterStaffFullName(ctx, chatID, update.Message.Text)
			h.log.Error("er", "erorr", err)
			msg.Text = "Введите номер телефона"
			tgbot.Send(msg)
			return
		case models.StateStaffRegisterPhoneNumber:
			err := h.state.RegisterStaffNumberPhone(ctx, chatID, update.Message.Text)
			h.log.Error("fss", "error", err)
			msg.Text = "Вы точно уверены"
			tgbot.Send(msg)
			return
		case models.StateStaffRegisterConfirm:
			h.state.ConfirmAddStaff(ctx, chatID)
			msg.Text = "сотрудник добавлен"
			tgbot.Send(msg)
			return
		}
	}

	if update.Message.Contact != nil {
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
