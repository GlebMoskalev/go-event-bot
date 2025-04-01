package schedule

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type schedule struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.Schedule {
	return &schedule{db: db, log: log}
}

func (s *schedule) GetAll(ctx context.Context, offset, limit int) ([]models.Schedule, int, error) {
	log := s.log.With("layer", "service_schedule", "operation", "GetAll")
	log.Info("getting all schedule")
	schedules, total, err := s.db.GetAllSchedules(ctx, offset, limit)
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

func (s *schedule) Update(ctx context.Context, schedule models.Schedule) error {
	log := s.log.With("layer", "service_schedule", "operation", "Update")
	log.Info("updating all schedule")
	err := s.db.UpdateSchedule(ctx, schedule)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundSchedule) {
			log.Warn("schedules not found in repository")
			return err
		}
		log.Error("failed to update schedule from repository")
		return err
	}
	return err
}

func (s *schedule) Create(ctx context.Context, schedule models.Schedule) error {
	log := s.log.With("layer", "service_schedule", "operation", "Create")
	log.Info("creating all schedule")
	err := s.db.CreateSchedule(ctx, schedule)
	if err != nil {
		log.Error("failed to create schedule from repository")
		return err
	}
	return err
}

func (s *schedule) Delete(ctx context.Context, scheduleId int) error {
	log := s.log.With("layer", "service_schedule", "operation", "Delete")
	log.Info("deleting all schedule")
	err := s.db.DeleteSchedule(ctx, scheduleId)
	if err != nil {
		log.Error("failed to delete schedule from repository")
		return err
	}
	return err
}
