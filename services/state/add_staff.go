package state

import (
	"context"
	"encoding/json"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"strings"
)

func (s *state) StartAddStaff(ctx context.Context, chatID int64) error {
	log := logger.SetupLogger(s.log, "service_state", "StartAddStaff", "chat_id", chatID)
	log.Info("starting staff registration process")

	err := s.db.SetState(ctx, chatID, models.StateStaffRegisterPhoneNumber, []byte("{}"))
	if err != nil {
		log.Error("failed to set initial state", "error", err)
		return err
	}

	log.Info("staff registration process started successfully")
	return err
}

func (s *state) RegisterStaffFullName(ctx context.Context, chatID int64, firstName, lastName, patronymic string) error {
	log := logger.SetupLogger(s.log,
		"service_state", "RegisterStaffFullName",
		"chat_id", chatID,
	)
	log.Info("registering staff full name")

	_, data, err := s.db.GetStateAndData(ctx, chatID)
	if err != nil {
		log.Error("failed to get current state and data", "error", err)
		return err
	}
	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		log.Error("failed to unmarshal staff data", "error", err)
		return err
	}
	staff.FirstName = firstName
	staff.LastName = lastName
	staff.Patronymic = patronymic

	data, err = json.Marshal(staff)
	if err != nil {
		log.Error("failed to marshal staff data", "error", err)
		return err
	}
	err = s.db.SetState(ctx, chatID, models.StateStaffRegisterConfirm, data)
	if err != nil {
		log.Error("failed to set state to phone number registration", "error", err)
	}

	log.Error("staff full name registered successfully")
	return nil
}

func (s *state) RegisterStaffNumberPhone(ctx context.Context, chatID int64, phoneNumber string) error {
	log := logger.SetupLogger(s.log,
		"service_state", "RegisterStaffNumberPhone",
		"chat_id", chatID,
		"phone", phoneNumber,
	)
	log.Info("registering staff phone number")

	_, data, err := s.db.GetStateAndData(ctx, chatID)
	if err != nil {
		log.Error("failed to get current state and data", "error", err)
		return err
	}
	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		log.Error("failed to unmarshal staff data", "error", err)
		return err
	}

	phoneNumber = strings.TrimSpace(phoneNumber)
	staff.PhoneNumber = phoneNumber

	updateData, err := json.Marshal(staff)
	if err != nil {
		log.Error("failed to marshal updated staff data", "error", err)
		return err
	}

	err = s.db.SetState(ctx, chatID, models.StateStaffRegisterFullName, updateData)
	if err != nil {
		log.Error("failed to set state to confirmation", "error", err)
		return err
	}

	log.Info("staff phone number registered successfully")
	return nil
}

func (s *state) ConfirmAddStaff(ctx context.Context, chatID int64) error {
	log := logger.SetupLogger(s.log, "service_state", "ConfirmAddStaff", "chat_id", chatID)
	log.Info("confirming staff registration")

	err := s.db.RemoveState(ctx, chatID)
	if err != nil {
		log.Error("failed to remove state after confirmation", "error", err)
		return err
	}

	log.Info("staff registration completed successfully")
	return nil
}
