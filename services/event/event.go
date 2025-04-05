package event

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

type event struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.Event {
	return &event{db: db, log: log}
}

func (s *event) GetAll(ctx context.Context, offset, limit int) ([]models.Event, int, error) {
	log := logger.SetupLogger(s.log,
		"service_event", "GetAll",
		"offset", offset,
		"limit", limit,
	)
	log.Info("getting all events")

	schedules, total, err := s.db.GetAllEvents(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundEvent) {
			log.Warn("events not found")
			return nil, 0, err
		}
		log.Error("failed to get events")
		return nil, 0, err
	}

	log.Info("events retrieved successfully")
	return schedules, total, err
}

func (s *event) Update(ctx context.Context, event models.Event) error {
	log := logger.SetupLogger(s.log,
		"service_event", "Update",
		"event_id", event.ID,
	)
	log.Info("updating all events")

	err := s.db.UpdateEvent(ctx, event)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundEvent) {
			log.Warn("event not found")
			return err
		}
		log.Error("failed to update event")
		return err
	}

	log.Info("event updated successfully")
	return err
}

func (s *event) Create(ctx context.Context, event models.Event) error {
	log := logger.SetupLogger(s.log,
		"service_event", "Create",
		"title", event.Title,
	)
	log.Info("creating all events")

	err := s.db.CreateEvent(ctx, event)
	if err != nil {
		log.Error("failed to create event")
		return err
	}

	log.Info("event created successfully")
	return err
}

func (s *event) Delete(ctx context.Context, eventId int) error {
	log := logger.SetupLogger(s.log,
		"service_event", "Delete",
		"event_id", eventId,
	)
	log.Info("deleting event")

	err := s.db.DeleteEvent(ctx, eventId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundEvent) {
			log.Warn("event not found")
			return err
		}
		log.Error("failed to delete event", "error", err)
		return err
	}

	log.Info("event deleted successfully")
	return err
}
