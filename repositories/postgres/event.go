package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
)

func (p *postgres) GetAllEvents(ctx context.Context, offset, limit int) ([]models.Event, int, error) {
	log := logger.SetupLogger(p.log, "repository_event", "GetAllEvents")
	log.Info("fetching all events")

	query := `
	SELECT
	    id,
	    title,
	    speaker,
	    auditorium,
	    status,
	    date
	FROM event
	WHERE status != 'completed'
	ORDER BY 
	    CASE WHEN status = 'ongoing' THEN 0 ELSE 1 END,
	    date
	LIMIT $1
	OFFSET $2
`

	rows, err := p.dbBot.QueryContext(ctx, query, limit, offset)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error("failed to close rows", "error", err)
		}
	}()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("no events found")
			return nil, 0, apperrors.ErrNotFoundEvent
		}
		log.Error("failed to fetch events", "error", err)
		return nil, 0, err
	}

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Speaker, &event.Auditorium, &event.Status, &event.Date)
		if err != nil {
			log.Error("failed to scan row", "error", err)
			return nil, 0, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Info("row processing error", "error", err)
		return nil, 0, err
	}

	query = `
	SELECT COUNT(*)
	FROM event
	WHERE status != 'completed'
`

	var total int
	err = p.dbBot.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		p.log.Error("failed to get events count", "error", err)
		return nil, 0, err
	}

	log.Info("successfully retrieved events", "count", len(events), "total", total)
	return events, total, nil
}

func (p *postgres) UpdateEvent(ctx context.Context, event models.Event) error {
	log := logger.SetupLogger(p.log, "repository_event", "Update", "event_id", event.ID)
	log.Info("updating event")

	query := `
	UPDATE event
	SET
	    title = $2,
	    speaker = $3,
	    auditorium = $4,
	    status = $5,
	    date = $6
	WHERE id = $1
`
	result, err := p.dbBot.ExecContext(ctx, query, event.ID, event.Title, event.Speaker, event.Auditorium, event.Status, event.Date)

	if err != nil {
		log.Error("failed to update event", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("event not found")
		return apperrors.ErrNotFoundEvent
	}

	log.Info("event updated successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) CreateEvent(ctx context.Context, event models.Event) error {
	log := logger.SetupLogger(p.log, "repository_event", "Create", "title", event.Title)
	log.Info("creating event")

	query := `
	INSERT INTO event 
	    (title, speaker, auditorium, status, date) 
	VALUES 
		($1, $2, $3, $4)
`
	result, err := p.dbBot.ExecContext(ctx, query, event.Title, event.Speaker, event.Auditorium, event.Status, event.Date)
	if err != nil {
		log.Error("failed to create event", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	log.Info("event created successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) DeleteEvent(ctx context.Context, eventID int) error {
	log := logger.SetupLogger(p.log, "repository_event", "Delete", "event_id", eventID)
	log.Info("deleting event")

	query := `
	DELETE FROM event
	WHERE id = $1
`

	result, err := p.dbBot.ExecContext(ctx, query, eventID)

	if err != nil {
		log.Error("failed to delete event", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("event not found")
		return apperrors.ErrNotFoundEvent
	}

	log.Info("event deleted successfully", "rows_affected", rowsAffected)
	return nil
}
