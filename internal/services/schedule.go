package services

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/repository"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	"log/slog"
)

type ScheduleService interface {
	GetAll(offset, limit int) ([]models.Schedule, int, error)
	Update(schedule models.Schedule) error
	Create(schedule models.Schedule) error
	Delete(scheduleId int) error
}

type scheduleService struct {
	repo repository.ScheduleRepository
	log  *slog.Logger
}

func NewScheduleService(repo repository.ScheduleRepository, log *slog.Logger) ScheduleService {
	return &scheduleService{repo: repo, log: log}
}

func (s *scheduleService) GetAll(offset, limit int) ([]models.Schedule, int, error) {
	log := s.log.With("layer", "service_schedule", "operation", "GetAll")
	log.Info("getting all schedule")
	schedules, total, err := s.repo.GetAll(context.Background(), offset, limit)
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

func (s *scheduleService) Update(schedule models.Schedule) error {
	log := s.log.With("layer", "service_schedule", "operation", "Update")
	log.Info("updating all schedule")
	err := s.repo.Update(context.Background(), schedule)
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

func (s *scheduleService) Create(schedule models.Schedule) error {
	log := s.log.With("layer", "service_schedule", "operation", "Create")
	log.Info("creating all schedule")
	err := s.repo.Create(context.Background(), schedule)
	if err != nil {
		log.Error("failed to create schedule from repository")
		return err
	}
	return err
}

func (s *scheduleService) Delete(scheduleId int) error {
	log := s.log.With("layer", "service_schedule", "operation", "Delete")
	log.Info("deleting all schedule")
	err := s.repo.Delete(context.Background(), scheduleId)
	if err != nil {
		log.Error("failed to delete schedule from repository")
		return err
	}
	return err
}
