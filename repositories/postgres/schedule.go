package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
)

func (p *postgres) GetAllSchedules(ctx context.Context, offset, limit int) ([]models.Schedule, int, error) {
	log := p.log.With("layer", "repository_schedule", "operation", "GetAll")
	log.Info("fetching all schedule")

	query := `
	SELECT
	    id,
	    title,
	    description,
	    date
	FROM schedule
	LIMIT $1
	OFFSET $2
`

	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error("failed to close rows")
		}
	}()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("schedules not found")
			return nil, 0, apperrors.ErrNotFoundSchedule
		}
		log.Error("failed to fetch all schedule")
		return nil, 0, err
	}

	var schedules []models.Schedule

	for rows.Next() {
		var schedule models.Schedule
		err = rows.Scan(&schedule.ID, &schedule.Title, &schedule.Description, &schedule.Date)
		if err != nil {
			log.Error("failed scan row in database", "error", err)
			return nil, 0, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		log.Info("failed rows error", "error", err)
		return nil, 0, err
	}

	query = `
	SELECT COUNT(*)
	FROM schedule
`

	var total int
	err = p.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		p.log.Error("failed to get count schedule")
		return nil, 0, err
	}

	log.Info("successfully retrieved all schedules")
	return schedules, total, nil
}

func (p *postgres) UpdateSchedule(ctx context.Context, schedule models.Schedule) error {
	log := p.log.With("layer", "repository_schedule", "operation", "Update", "schedule_id", schedule.ID)
	log.Info("updating schedule")

	query := `
	UPDATE schedule
	SET
	    title = $2,
	    description = $3,
	    date = $4
	WHERE id = $1
`
	result, err := p.db.ExecContext(ctx, query, schedule.ID, schedule.Title, schedule.Description, schedule.Date)

	if err != nil {
		log.Error("failed to update schedule", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("not found schedule")
		return apperrors.ErrNotFoundSchedule
	}

	log.Info("schedule updated successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) CreateSchedule(ctx context.Context, schedule models.Schedule) error {
	log := p.log.With("layer", "repository_schedule", "operation", "Create", "schedule_title", schedule.Title)
	log.Info("creating schedule")

	query := `
	INSERT INTO schedule 
	    (title, description, date) 
	VALUES 
		($1, $2, $3)
`
	result, err := p.db.ExecContext(ctx, query, schedule.Title, schedule.Description, schedule.Date)
	if err != nil {
		log.Error("failed to create schedule", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	log.Info("schedule created successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) DeleteSchedule(ctx context.Context, scheduleId int) error {
	log := p.log.With("layer", "repository_schedule", "operation", "Delete", "schedule_id", scheduleId)
	log.Info("deleting schedule")

	query := `
	DELETE FROM schedule
	WHERE id = $1
`

	result, err := p.db.ExecContext(ctx, query, scheduleId)

	if err != nil {
		log.Error("failed to delete schedule", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("not found schedule")
		return apperrors.ErrNotFoundSchedule
	}

	log.Info("schedule deleted successfully", "rows_affected", rowsAffected)
	return nil
}
