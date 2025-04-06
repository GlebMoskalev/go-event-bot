package message

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
	"unicode/utf8"
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
		msg.ReplyMarkup = keyboards.CancelAddStaff()
		return m.stateStaffRegisterFullName(ctx, msg)
	case models.StateStaffRegisterPhoneNumber:
		msg.ReplyMarkup = keyboards.CancelAddStaff()
		return m.stateStaffRegisterPhoneNumber(ctx, msg)
	case models.StateSearchLastName:
		return m.stateSearchLastName(ctx, msg)
	default:
		log.Warn("unknown state encountered")
		msg.Text = "Произошла ошибка"
		return msg
	}
}

func (m *msg) stateSearchLastName(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(m.log,
		"service_message", "stateSearchLastName",
		"chat_id", msg.ChatID,
		"text", msg.Text,
	)
	log.Info("processing get staff by last name")

	lastName := strings.TrimSpace(msg.Text)
	if utf8.RuneCountInString(lastName) < 2 {
		log.Warn("last name too short", "last_name", lastName)
		msg.Text = messages.LastNameTooShort()
		return msg
	}

	err := m.stateService.RemoveState(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to remove state", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	staffList, err := m.staffService.GetListByPhoneOrLastName(ctx, "", lastName)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("no staff found for last name")
			msg.Text = messages.StaffNodFound()
			return msg
		}
		log.Error("failed to retrieve staff list", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	log.Info("staff list retrieved successfully", "count", len(staffList))
	msg.Text = messages.StaffList(staffList)
	return msg
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
		msg.Text = messages.InvalidFullNameFormat()
		return msg
	}

	err := m.stateService.RegisterStaffFullName(ctx, msg.ChatID, fullNameSplit[1], fullNameSplit[0], fullNameSplit[2])
	if err != nil {
		log.Error("failed to register staff full name", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	_, data, err := m.stateService.GetWithData(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to retrieve staff data", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		log.Error("failed to unmarshal staff data", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	msg.Text = messages.ConfirmStaffAddition(staff.LastName, staff.FirstName, staff.Patronymic, staff.PhoneNumber)
	msg.ReplyMarkup = keyboards.AgreeStaff()
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
		msg.Text = messages.InvalidPhoneFormat()
		return msg
	}

	_, err := m.staffService.GetByPhoneNumber(ctx, msg.Text)
	if !errors.Is(err, apperrors.ErrNotFoundStaff) {
		log.Warn("staff with this phone number already exists")
		msg.Text = messages.PhoneAlreadyExists()
		return msg
	}

	err = m.stateService.RegisterStaffNumberPhone(ctx, msg.ChatID, msg.Text)
	if err != nil {
		log.Error("failed to register staff phone number", "error", err)
		msg.Text = messages.Error()
		return msg
	}

	msg.Text = messages.RequestFullName()

	log.Info("phone number registered successfully")
	return msg
}
