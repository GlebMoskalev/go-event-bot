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
		log.Error("failed to fetch staff from database", "error", err)
		return models.Staff{}, err
	}

	log.Info("staff retrieved successfully")
	return staff, nil
}

func (p *postgres) CreateStaff(ctx context.Context, staff models.Staff) error {
	log := p.log.With("layer", "repository_staff", "operation", "CreateStaff",
		"first_name", staff.FirstName,
		"last_name", staff.LastName,
		"phone_number", staff.PhoneNumber)
	log.Info("creating new staff record")

	query := `
	INSERT INTO staffs 
	    (firstname, lastname, patronymic, phone_number)
	VALUES 
	    ($1, $2, $3, $4)
`
	_, err := p.dbStaff.ExecContext(ctx, query, staff.FirstName, staff.LastName, staff.Patronymic, staff.PhoneNumber)
	if err != nil {
		log.Error("failed to create staff record", "error", err)
		return err
	}

	log.Info("staff record created successfully")
	return nil
}

func (p *postgres) UpdateStaff(ctx context.Context, staff models.Staff) error {
	log := p.log.With("layer", "repository_staff", "operation", "UpdateStaff",
		"first_name", staff.FirstName,
		"last_name", staff.LastName,
		"phone_number", staff.PhoneNumber)
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
		log.Error("failed to update staff record", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("staff not found for update")
		return apperrors.ErrNotFoundStaff
	}

	log.Info("staff record updated successfully")
	return nil
}

func (p *postgres) DeleteStaff(ctx context.Context, phoneNumber string) error {
	log := p.log.With("layer", "repository_staff", "operation", "DeleteStaff", "phone_number", phoneNumber)
	log.Info("deleting staff record")

	query := `
	DELETE FROM staffs
	WHERE phone_number = $1
	`

	result, err := p.dbStaff.ExecContext(ctx, query, phoneNumber)
	if err != nil {
		log.Error("failed to delete staff record", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("staff not found for deletion")
		return apperrors.ErrNotFoundStaff
	}

	log.Info("staff record deleted successfully")
	return nil
}
