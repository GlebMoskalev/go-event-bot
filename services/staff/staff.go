package staff

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
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
	log := logger.SetupLogger(s.log,
		"service_staff", "GetByPhoneNumber",
		"phone_number", phoneNumber,
	)
	log.Info("getting staff by phone number")

	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !re.MatchString(phoneNumber) {
		log.Warn("invalid phone number format")
		return models.Staff{}, apperrors.ErrInvalidPhoneNumber
	}
	phoneNumber = strings.TrimPrefix(phoneNumber, "+")

	staff, err := s.db.GetStaffByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("staff not found")
			return models.Staff{}, err
		}
		log.Error("failed to retrieve staff", "error", err)
		return models.Staff{}, err
	}

	log.Info("staff retrieved successfully")
	return staff, err
}

func (s *staff) Create(ctx context.Context, staff models.Staff) error {
	log := logger.SetupLogger(s.log,
		"service_staff", "Create",
		"phone_number", staff.PhoneNumber,
	)
	log.Info("creating new staff member")

	if staff.Role == "" {
		staff.Role = models.RoleStaff
	}
	if err := s.db.CreateStaff(ctx, staff); err != nil {
		log.Error("failed to create staff", "error", err)
		return err
	}

	log.Info("staff created successfully")
	return nil
}

func (s *staff) Update(ctx context.Context, staff models.Staff) error {
	log := logger.SetupLogger(s.log, "service_staff", "Update", "phone_number", staff.PhoneNumber)
	log.Info("updating staff member")

	if staff.Role == "" {
		staff.Role = models.RoleStaff
	}
	if err := s.db.UpdateStaff(ctx, staff); err != nil {
		log.Error("failed to update staff", "error", err)
		return err
	}

	log.Info("staff updated successfully")
	return nil
}
