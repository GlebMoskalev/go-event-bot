package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
)

func (p *postgres) GetState(ctx context.Context, chatID int64) (models.State, error) {
	log := p.log.With("layer", "repository_state", "operation", "GetState", "chat_id", chatID)
	log.Info("fetching state")

	query := `
	SELECT state
	FROM states
	WHERE chat_id = $1
`
	var state models.State
	err := p.dbBot.QueryRowContext(ctx, query, chatID).Scan(&state)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("state not found")
			return "", apperrors.ErrNotFoundState
		}
		log.Error("failed to fetch state", "error", err)
		return "", err
	}

	log.Info("state retrieved successfully")
	return state, nil
}

func (p *postgres) GetStateAndData(ctx context.Context, chatID int64) (models.State, []byte, error) {
	log := p.log.With("layer", "repository_state", "operation", "GetStateAndData", "chat_id", chatID)
	log.Info("fetching state and data")

	query := `
	SELECT state, data
	FROM states
	WHERE chat_id = $1
`
	var data []byte
	var state models.State
	err := p.dbBot.QueryRowContext(ctx, query, chatID).Scan(&state, &data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("state and data not found")
			return "", nil, apperrors.ErrNotFoundState
		}
		log.Error("failed to fetch state and data", "error", err)
		return "", nil, err
	}

	log.Info("state and data retrieved successfully")
	return state, data, nil
}

func (p *postgres) SetState(ctx context.Context, chatID int64, state models.State, dataJSON []byte) error {
	log := p.log.With("layer", "repository_state", "operation", "SetState", "chat_id", chatID, "state", state)
	log.Info("setting state")

	query := `
	INSERT INTO states (chat_id, state, data)
	VALUES ($1, $2, $3)
	ON CONFLICT (chat_id)
	DO UPDATE SET state = $2, data = $3
`
	_, err := p.dbBot.ExecContext(ctx, query, chatID, state, dataJSON)
	if err != nil {
		log.Error("failed to set state", "error", err)
		return err
	}

	log.Info("state set successfully")
	return nil
}

func (p *postgres) RemoveState(ctx context.Context, chatID int64) error {
	log := p.log.With("layer", "repository_state", "operation", "RemoveState", "chat_id", chatID)
	log.Info("removing state")

	query := `
	DELETE FROM states
	WHERE chat_id = $1
`
	result, err := p.dbBot.ExecContext(ctx, query, chatID)
	if err != nil {
		log.Error("failed to remove state", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("state not found")
		return apperrors.ErrNotFoundState
	}

	log.Info("state removed successfully")
	return nil
}
