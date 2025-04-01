package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
	"github.com/GlebMoskalev/go-event-bot/pkg/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *msg) Contact(ctx context.Context, msg tgbotapi.MessageConfig, contact *tgbotapi.Contact) tgbotapi.MessageConfig {
	log := m.log.With("layer", "service_message", "operation", "Contact", "chat_id", msg.ChatID)

	exists, err := m.userService.ExistsUserByTelegramID(ctx, contact.UserID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		msg.Text = "Произошла ошибка!"
		return msg
	}

	if exists {
		//removeKeyboard := keyboards.RemoveKeyboard()
		msg.ReplyMarkup = m.commandService.SetupCommands(msg, false)
		msg.Text = "Твой контакт у нас уже есть"
		return msg
	}

	phoneNumber := contact.PhoneNumber
	staff, err := m.staffService.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = "Ваш номер телефона не соответствует нашему формату"
			return msg
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = "Вас нет в списках, попросите администратора добавить вас"
			return msg
		}

		log.Error("failed to get staff in service", "error", err)

		msg.Text = "Произошла ошибка!"
		return msg
	}

	msg.ReplyMarkup = m.commandService.SetupCommands(msg, false)

	err = m.userService.Create(ctx, models.User{
		TelegramID: contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
		IsAdmin:    false,
	})

	if err != nil {
		msg.Text = "Произошла ошибка!"
		return msg
	}

	msg.Text = fmt.Sprintf("%s %s, приветствуем на нашем мероприятии!",
		staff.FirstName,
		staff.Patronymic)
	msg.ReplyMarkup = keyboards.RemoveKeyboard()
	return msg
}
