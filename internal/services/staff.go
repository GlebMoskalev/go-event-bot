package services

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/repository"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	"log/slog"
	"regexp"
	"strings"
)

type StaffService interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error)
}

type staffService struct {
	repo repository.StaffRepository
	log  *slog.Logger
}

func NewStaffService(repo repository.StaffRepository, log *slog.Logger) StaffService {
	return &staffService{repo: repo, log: log}
}

func (s *staffService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error) {
	log := s.log.With("layer", "service_staff", "operation", "GetByPhoneNumber", "phone_number", phoneNumber)
	log.Info("getting staff by phone number")

	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !re.MatchString(phoneNumber) {
		log.Warn(apperrors.ErrInvalidPhoneNumber.Error())
		return models.Staff{}, apperrors.ErrInvalidPhoneNumber
	}

	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	staff, err := s.repo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("staff not found in repository")
			return models.Staff{}, err
		}
		log.Error("failed to retrieve staff from repository", "error", err)
		return models.Staff{}, err
	}

	log.Info("staff retrieved successfully")
	return staff, err
}
