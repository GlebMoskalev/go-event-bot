package staff

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
	"regexp"
	"strings"
)

type staff struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.Staff {
	return &staff{db: db, log: log}
}

func (s *staff) GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error) {
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

	staff, err := s.db.GetStaffByPhoneNumber(ctx, phoneNumber)
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
