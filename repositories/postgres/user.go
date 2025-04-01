package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
)

func (p *postgres) GetUser(ctx context.Context, telegramID int64) (models.User, error) {
	log := p.log.With("layer", "repository_user", "operation", "Get", "telegram_id", telegramID)
	log.Info("fetching user data")

	query := `
	SELECT
		firstname,
		lastname,
		patronymic,
		role
	FROM users
	WHERE telegram_id = $1
`
	var user models.User
	err := p.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user not found")
			return models.User{}, apperrors.ErrNotFoundStaff
		}

		log.Error("failed to fetch user from database", "error", err)
	}

	log.Info("user retrieved successfully")
	user.TelegramID = telegramID
	return user, err
}

func (p *postgres) CreateUser(ctx context.Context, user models.User) error {
	log := p.log.With("layer", "repository_user", "operation", "Create", "telegram_id", user.TelegramID)
	log.Info("creating new user")

	query := `
	INSERT INTO users 
		(telegram_id, firstname, lastname, patronymic, role) 
	VALUES 
		($1, $2, $3, $4, $5)
`

	_, err := p.db.ExecContext(ctx, query, user.TelegramID, user.FirstName, user.LastName, user.Patronymic, user.Role)
	if err != nil {
		log.Error("failed to create user in database", "error", err)
		return err
	}

	log.Info("user created successfully")
	return nil
}

func (p *postgres) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := p.log.With("layer", "repository_user", "operation", "ExistsUserByTelegramID", "telegram_id", telegramID)
	log.Info("checking user existence by telegram_id")

	query := `
	SELECT 
	EXISTS(
		SELECT 1 FROM users WHERE telegram_id = $1
	)
`
	var exists bool
	err := p.db.QueryRowContext(ctx, query, telegramID).Scan(&exists)

	if err != nil {
		log.Error("failed to verify user existence using telegram_id", "error", err)
		return false, err
	}

	log.Info("user existence by telegram_id checked successfully")
	return exists, err
}
