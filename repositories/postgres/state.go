package postgres

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

func (p *postgres) GetState(ctx context.Context, chatID int64) (models.State, error) {
	query := `
	SELECT state
	FROM states
	WHERE chat_id = $1
`
	var state models.State
	err := p.dbBot.QueryRowContext(ctx, query, chatID).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

func (p *postgres) GetStateAndData(ctx context.Context, chatID int64) (models.State, []byte, error) {
	query := `
	SELECT state, data
	FROM states
	WHERE chat_id = $1
`
	var data []byte
	var state models.State
	err := p.dbBot.QueryRowContext(ctx, query, chatID).Scan(&state, &data)
	if err != nil {
		return "", nil, err
	}
	return state, data, nil
}

func (p *postgres) SetState(ctx context.Context, chatID int64, state models.State, dataJSON []byte) error {
	query := `
	INSERT INTO states (chat_id, state, data)
	VALUES ($1, $2, $3)
	ON CONFLICT (chat_id)
	DO UPDATE SET state = $2, data = $3
`
	_, err := p.dbBot.ExecContext(ctx, query, chatID, state, dataJSON)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgres) RemoveState(ctx context.Context, chatID int64) error {
	query := `
	DELETE FROM states
	WHERE chat_id = $1
`
	_, err := p.dbBot.ExecContext(ctx, query, chatID)
	if err != nil {
		return err
	}
	return nil
}
