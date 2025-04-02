package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
)

func (p *postgres) GetStaffByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error) {
	log := p.log.With("layer", "repository_staff", "operation", "GetByPhoneNumber", "phone_number", phoneNumber)
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

	err := p.dbStaff.QueryRowContext(ctx, query, phoneNumber).Scan(
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
