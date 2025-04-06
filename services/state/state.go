package state

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"log/slog"
	"strings"
)

type state struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.State {
	return &state{db: db, log: log}
}

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

func (s *state) Get(ctx context.Context, chatID int64) (models.State, error) {
	log := logger.SetupLogger(s.log, "service_state", "Get", "chat_id", chatID)
	log.Info("getting current state")

	state, err := s.db.GetState(ctx, chatID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundState) {
			log.Warn("no state found")
			return "", err
		}
		log.Error("failed to get state", "error", err)
		return "", err
	}

	log.Info("state retrieved successfully", "state", state)
	return state, nil
}

func (s *state) GetWithData(ctx context.Context, chatID int64) (models.State, []byte, error) {
	log := logger.SetupLogger(s.log, "service_state", "GetWithData", "chat_id", chatID)
	log.Info("getting current state with data")

	state, data, err := s.db.GetStateAndData(ctx, chatID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundState) {
			log.Warn("no state found for chat")
			return "", nil, err
		}
		log.Error("failed to get state with data", "error", err)
		return "", nil, err
	}

	log.Info("state with data retrieved successfully", "state", state)
	return state, data, nil
}

func (s *state) RemoveState(ctx context.Context, chatID int64) error {
	log := logger.SetupLogger(s.log, "service_state", "RemoveState", "chat_id", chatID)
	log.Info("removing state")

	err := s.db.RemoveState(ctx, chatID)
	if err != nil {
		log.Error("failed to delete state", "error", err)
		return err
	}

	log.Info("state remove successfully")
	return nil
}
