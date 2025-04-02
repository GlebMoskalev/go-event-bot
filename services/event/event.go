package event

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"log/slog"
)

type event struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.Event {
	return &event{db: db, log: log}
}

func (s *event) GetAll(ctx context.Context, offset, limit int) ([]models.Event, int, error) {
	log := s.log.With("layer", "service_schedule", "operation", "GetAll")
	log.Info("getting all event")
	schedules, total, err := s.db.GetAllEvents(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundSchedule) {
			log.Warn("schedules not found in repository")
			return nil, 0, err
		}
		log.Error("failed to get schedules from repository")
		return nil, 0, err
	}
	return schedules, total, err
}

func (s *event) Update(ctx context.Context, event models.Event) error {
	log := s.log.With("layer", "service_schedule", "operation", "Update")
	log.Info("updating all event")
	err := s.db.UpdateEvent(ctx, event)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundSchedule) {
			log.Warn("schedules not found in repository")
			return err
		}
		log.Error("failed to update event from repository")
		return err
	}
	return err
}

func (s *event) Create(ctx context.Context, event models.Event) error {
	log := s.log.With("layer", "service_schedule", "operation", "Create")
	log.Info("creating all event")
	err := s.db.CreateEvent(ctx, event)
	if err != nil {
		log.Error("failed to create event from repository")
		return err
	}
	return err
}

func (s *event) Delete(ctx context.Context, eventId int) error {
	log := s.log.With("layer", "service_schedule", "operation", "Delete")
	log.Info("deleting all event")
	err := s.db.DeleteEvent(ctx, eventId)
	if err != nil {
		log.Error("failed to delete event from repository")
		return err
	}
	return err
}
