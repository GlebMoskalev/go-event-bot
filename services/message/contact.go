package message

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	messages2 "github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *msg) Contact(ctx context.Context, msg tgbotapi.MessageConfig,
	contact *tgbotapi.Contact) (tgbotapi.MessageConfig, models.Role, error) {
	log := logger.SetupLogger(m.log,
		"service_message", "Contact",
		"chat_id", msg.ChatID,
		"telegram_id", contact.UserID,
	)
	log.Info("processing contact message")

	var role models.Role = models.RoleGuest

	exists, err := m.userService.ExistsUserByTelegramID(ctx, contact.UserID)
	if err != nil {
		log.Error("failed to check user existence", "error", err)
		msg.Text = messages2.Error()
		return msg, role, err
	}

	if exists {
		log.Info("user already exists")
		user, err := m.userService.Get(ctx, contact.UserID)
		if err != nil {
			log.Error("failed to get existing user", "error", err)
			msg.Text = messages2.Error()
			return msg, role, nil
		}

		msg.Text = messages2.ContactExists()
		return msg, user.Role, nil
	}

	phoneNumber := contact.PhoneNumber
	staff, err := m.staffService.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = messages2.InvalidPhoneNumber()
			return msg, role, err
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("staff not found", "phone_number", phoneNumber)
			msg.ReplyMarkup = keyboards.RemoveKeyboard()
			msg.Text = messages2.StaffNotFound()
			return msg, role, err
		}

		log.Error("failed to get staff", "error", err)

		msg.Text = messages2.Error()
		return msg, role, err
	}
	role = staff.Role

	err = m.userService.Create(ctx, models.User{
		TelegramID: contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
		Role:       role,
		ChatID:     msg.ChatID,
	})

	if err != nil {
		log.Error("failed to create user", "error", err)
		msg.Text = messages2.Error()
		return msg, role, err
	}

	msg.Text = messages2.Welcome(
		staff.FirstName,
		staff.Patronymic)
	msg.ReplyMarkup = keyboards.RemoveKeyboard()

	log.Info("contact processed successfully")
	return msg, role, nil
}
