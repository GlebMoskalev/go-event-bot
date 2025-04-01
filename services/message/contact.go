package message

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	messages2 "github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *msg) Contact(ctx context.Context, msg tgbotapi.MessageConfig, contact *tgbotapi.Contact) (tgbotapi.MessageConfig, error) {
	log := m.log.With("layer", "service_message", "operation", "Contact", "chat_id", msg.ChatID)

	exists, err := m.userService.ExistsUserByTelegramID(ctx, contact.UserID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		msg.Text = messages2.Error()
		return msg, err
	}

	if exists {
		msg.Text = messages2.ContactExists()
		return msg, nil
	}

	phoneNumber := contact.PhoneNumber
	staff, err := m.staffService.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = messages2.InvalidPhoneNumber()
			return msg, err
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = messages2.StaffNotFound()
			return msg, err
		}

		log.Error("failed to get staff in service", "error", err)

		msg.Text = messages2.Error()
		return msg, err
	}

	err = m.userService.Create(ctx, models.User{
		TelegramID: contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
		Role:       models.RoleStaff,
	})

	if err != nil {
		msg.Text = messages2.Error()
		return msg, err
	}

	msg.Text = messages2.Welcome(
		staff.FirstName,
		staff.Patronymic)
	msg.ReplyMarkup = keyboards.RemoveKeyboard()
	return msg, nil
}
