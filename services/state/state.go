package state

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"log/slog"
)

type state struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.State {
	return &state{db: db, log: log}
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
