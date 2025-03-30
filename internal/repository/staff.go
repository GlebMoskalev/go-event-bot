package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	"log/slog"
)

type StaffRepository interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error)
}

type staffRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStaffRepository(db *sql.DB, logger *slog.Logger) StaffRepository {
	return &staffRepo{db: db, log: logger}
}

func (r *staffRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error) {
	log := r.log.With("layer", "repository_staff", "operation", "GetByPhoneNumber", "phone_number", phoneNumber)
	log.Info("fetching staff data")

	var staff models.Staff

	query := `
	SELECT
		firstname, 
		lastname, 
		patronymic, 
		email,
		phone_number
	FROM staffs
	WHERE phone_number = $1
`

	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&staff.FirstName,
		&staff.LastName,
		&staff.Patronymic,
		&staff.Email,
		&staff.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("staff not found")
			return models.Staff{}, apperrors.ErrNotFoundStaff
		}
		log.Error("failed to fetch staff from database", "error", err)
		return models.Staff{}, err
	}

	log.Info("staff retrieved successfully")
	return staff, nil
}
