package staff

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
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

func (s *staff) Create(ctx context.Context, staff models.Staff) error {
	log := s.log.With("layer", "service_staff", "operation", "Create", "phone_number", staff.PhoneNumber)
	log.Info("creating new staff member")

	if err := s.db.CreateStaff(ctx, staff); err != nil {
		log.Error("failed to create staff in repository", "error", err)
		return err
	}

	log.Info("staff created successfully")
	return nil
}

func (s *staff) Update(ctx context.Context, staff models.Staff) error {
	log := s.log.With("layer", "service_staff", "operation", "Update", "phone_number", staff.PhoneNumber)
	log.Info("updating staff member")

	if err := s.db.UpdateStaff(ctx, staff); err != nil {
		log.Error("failed to update staff in repository", "error", err)
		return err
	}

	log.Info("staff updated successfully")
	return nil
}
