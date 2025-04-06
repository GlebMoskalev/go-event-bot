package state

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
)

func (s *state) StartSearchByLastName(ctx context.Context, chatID int64) error {
	log := logger.SetupLogger(s.log,
		"service_state", "StartSearchByLastName",
		"chat_id", chatID,
	)
	log.Info("starting search by last name")

	err := s.db.SetState(ctx, chatID, models.StateSearchLastName, []byte("{}"))
	if err != nil {
		log.Error("failed to set state for search by last name", "error", err)
		return err
	}

	log.Info("state set successfully for search by last name")
	return nil
}

func (s *state) StartSearchByPhoneNumber(ctx context.Context, chatID int64) error {
	log := logger.SetupLogger(s.log,
		"service_state", "StartSearchByPhoneNumber",
		"chat_id", chatID,
	)
	log.Info("starting search by phone number")

	err := s.db.SetState(ctx, chatID, models.StateSearchPhoneNumber, []byte("{}"))
	if err != nil {
		log.Error("failed to set state for search by phone number", "error", err)
		return err
	}

	log.Info("state set successfully for search by phone number")
	return nil
}
