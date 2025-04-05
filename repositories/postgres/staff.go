package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
)

func (p *postgres) GetStaffByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error) {
	log := logger.SetupLogger(p.log,
		"repository_staff", "GetByPhoneNumber",
		"phone_number", phoneNumber,
	)
	log.Info("fetching staff data")

	var staff models.Staff

	query := `
	SELECT
		firstname, 
		lastname, 
		patronymic, 
		phone_number
	FROM staffs
	WHERE phone_number = $1
`

	err := p.dbStaff.QueryRowContext(ctx, query, phoneNumber).Scan(
		&staff.FirstName,
		&staff.LastName,
		&staff.Patronymic,
		&staff.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("staff not found")
			return models.Staff{}, apperrors.ErrNotFoundStaff
		}
		log.Error("failed to fetch staff from", "error", err)
		return models.Staff{}, err
	}

	log.Info("staff retrieved successfully")
	return staff, nil
}

func (p *postgres) CreateStaff(ctx context.Context, staff models.Staff) error {
	log := logger.SetupLogger(p.log,
		"repository_staff", "CreateStaff",
		"first_name", staff.FirstName,
		"last_name", staff.LastName,
		"phone_number", staff.PhoneNumber,
	)
	log.Info("creating new staff")

	query := `
	INSERT INTO staffs 
	    (firstname, lastname, patronymic, phone_number)
	VALUES 
	    ($1, $2, $3, $4)
`
	result, err := p.dbStaff.ExecContext(ctx, query, staff.FirstName, staff.LastName, staff.Patronymic, staff.PhoneNumber)
	if err != nil {
		log.Error("failed to create staff", "error", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	log.Info("staff created successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) UpdateStaff(ctx context.Context, staff models.Staff) error {
	log := logger.SetupLogger(p.log,
		"repository_staff", "UpdateStaff",
		"phone_number", staff.PhoneNumber,
	)
	log.Info("updating staff record")

	query := `
	UPDATE staffs
	SET firstname = $1, 
		lastname = $2, 
		patronymic = $3
	WHERE phone_number = $4
`

	result, err := p.dbStaff.ExecContext(ctx, query, staff.FirstName, staff.LastName, staff.Patronymic, staff.PhoneNumber)
	if err != nil {
		log.Error("failed to update staff", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("staff not found")
		return apperrors.ErrNotFoundStaff
	}

	log.Info("staff updated successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) DeleteStaff(ctx context.Context, phoneNumber string) error {
	log := logger.SetupLogger(p.log, "repository_staff", "DeleteStaff", "phone_number", phoneNumber)
	log.Info("deleting staff")

	query := `
	DELETE FROM staffs
	WHERE phone_number = $1
	`

	result, err := p.dbStaff.ExecContext(ctx, query, phoneNumber)
	if err != nil {
		log.Error("failed to delete staff", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("staff not found")
		return apperrors.ErrNotFoundStaff
	}

	log.Info("staff record deleted successfully", "rows_affected", rowsAffected)
	return nil
}
