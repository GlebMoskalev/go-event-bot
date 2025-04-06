package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

func (m *msg) State(ctx context.Context, msg tgbotapi.MessageConfig, state models.State) tgbotapi.MessageConfig {
	log := logger.SetupLogger(m.log,
		"service_message",
		"State", "chat_id",
		msg.ChatID, "state", state,
	)
	log.Info("processing message based on state")

	switch state {
	case models.StateStaffRegisterFullName:
		return m.stateStaffRegisterFullName(ctx, msg)
	case models.StateStaffRegisterPhoneNumber:
		return m.stateStaffRegisterPhoneNumber(ctx, msg)
	case models.StateStaffRegisterConfirm:
		return m.stateStaffRegisterConfirm(ctx, msg)
	default:
		log.Warn("unknown state encountered")
		msg.Text = "Произошла ошибка"
		return msg
	}
}

func (m *msg) stateStaffRegisterFullName(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(m.log,
		"service_message", "StateStaffRegisterFullName",
		"chat_id", msg.ChatID,
		"text", msg.Text,
	)
	log.Info("processing full name registration")

	fullNameSplit := strings.Split(msg.Text, " ")
	if len(fullNameSplit) < 3 {
		log.Warn("invalid full name format", "name_parts", len(fullNameSplit))
		msg.Text = "Неправильный формат. Введите ФИО полностью (Фамилия Имя Отчество)"
		return msg
	}

	err := m.stateService.RegisterStaffFullName(ctx, msg.ChatID, fullNameSplit[0], fullNameSplit[1], fullNameSplit[2])
	if err != nil {
		log.Error("failed to register staff full name", "error", err)
		msg.Text = "Произошла ошибка"
		return msg
	}

	_, data, err := m.stateService.GetWithData(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to retrieve staff data", "error", err)
		msg.Text = "Произошла ошибка"
		return msg
	}

	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		log.Error("failed to unmarshal staff data", "error", err)
		msg.Text = "Произошла ошибка"
		return msg
	}

	msg.Text = fmt.Sprintf(
		"Подтвердите, что вы хотите добавить сотрудника:\n%s %s %s\n%s\n\n Введите 'да', если подтверждаете",
		staff.LastName, staff.FirstName, staff.Patronymic,
		staff.PhoneNumber,
	)

	log.Info("staff full name registered successfully")
	return msg
}

func (m *msg) stateStaffRegisterPhoneNumber(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(m.log,
		"service_message", "StateStaffRegisterPhoneNumber",
		"chat_id", msg.ChatID,
		"phone_number", msg.Text,
	)
	log.Info("processing phone number registration")

	re := regexp.MustCompile(`^[1-9]\d{1,14}$`)
	if !re.MatchString(msg.Text) {
		log.Warn("invalid phone number format")
		msg.Text = "Неправильный формат номера телефона. Введите в формате 79532597271"
		return msg
	}

	_, err := m.staffService.GetByPhoneNumber(ctx, msg.Text)
	if !errors.Is(err, apperrors.ErrNotFoundStaff) {
		err = m.stateService.RemoveState(ctx, msg.ChatID)
		if err != nil {
			log.Error("failed to remove state", "error", err)
		}
		log.Warn("staff with this phone number already exists")
		msg.Text = "Сотрудник с таким номером телефона существует"
		return msg
	}

	err = m.stateService.RegisterStaffNumberPhone(ctx, msg.ChatID, msg.Text)
	if err != nil {
		log.Error("failed to register staff phone number", "error", err)
		msg.Text = "Неправильный формат номера телефона. Введите в формате 79137777777"
		return msg
	}

	msg.Text = "Введите ФИО сотрудника в формате:\n Иван Иванов Иванович"

	log.Info("phone number registered successfully")
	return msg
}

func (m *msg) stateStaffRegisterConfirm(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(m.log,
		"service_message", "StateStaffRegisterConfirm",
		"chat_id", msg.ChatID,
		"response", msg.Text,
	)
	log.Info("processing staff registration confirmation")

	if msg.Text == "да" {
		_, data, err := m.stateService.GetWithData(ctx, msg.ChatID)
		if err != nil {
			log.Error("failed to retrieve staff data", "error", err)
			msg.Text = "Произошла ошибка"
			return msg
		}
		var staff models.Staff
		err = json.Unmarshal(data, &staff)
		if err != nil {
			log.Error("failed to unmarshal staff data", "error", err)
			msg.Text = "Произошла ошибка"
			return msg
		}

		err = m.staffService.Create(ctx, staff)
		if err != nil {
			log.Error("failed to create staff", "error", err)
			msg.Text = "Произошла ошибка"
			return msg
		}

		log.Info("staff member created successfully")
		msg.Text = "Сотрудник добавлен"
	} else {
		log.Info("staff creation cancelled")
		msg.Text = "Добавление сотрудника отменено"
	}

	err := m.stateService.ConfirmAddStaff(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to reset state", "error", err)
		msg.Text = "Произошла ошибка"
		return msg
	}
	return msg
}
