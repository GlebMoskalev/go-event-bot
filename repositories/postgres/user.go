package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
)

func (p *postgres) GetUser(ctx context.Context, telegramID int64) (models.User, error) {
	log := logger.SetupLogger(p.log, "repository_user", "GetUser", "telegram_id", telegramID)
	log.Info("fetching user data")

	query := `
	SELECT
		firstname,
		lastname,
		patronymic,
		role,
		chat_id
	FROM users
	WHERE telegram_id = $1
`
	var user models.User
	err := p.dbBot.QueryRowContext(ctx, query, telegramID).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.Role,
		&user.ChatID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user not found")
			return models.User{}, apperrors.ErrNotFoundStaff
		}

		log.Error("failed to fetch user", "error", err)
	}

	user.TelegramID = telegramID
	log.Info("user retrieved successfully")
	return user, nil
}

func (p *postgres) CreateUser(ctx context.Context, user models.User) error {
	log := logger.SetupLogger(p.log,
		"repository_user", "CreateUser",
		"telegram_id", user.TelegramID,
	)
	log.Info("creating new user")

	query := `
	INSERT INTO users 
		(telegram_id, firstname, lastname, patronymic, role, chat_id) 
	VALUES 
		($1, $2, $3, $4, $5, $6)
`

	result, err := p.dbBot.ExecContext(ctx, query, user.TelegramID, user.FirstName, user.LastName, user.Patronymic,
		user.Role, user.ChatID)
	if err != nil {
		log.Error("failed to create user", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", "error", err)
		return err
	}

	log.Info("user created successfully", "rows_affected", rowsAffected)
	return nil
}

func (p *postgres) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := logger.SetupLogger(p.log,
		"repository_user", "ExistsUserByTelegramID",
		"telegram_id", telegramID,
	)
	log.Info("checking user existence")

	query := `
	SELECT 
	EXISTS(
		SELECT 1 FROM users WHERE telegram_id = $1
	)
`
	var exists bool
	err := p.dbBot.QueryRowContext(ctx, query, telegramID).Scan(&exists)

	if err != nil {
		log.Error("failed to check user existence", "error", err)
		return false, err
	}

	log.Info("user existence checked successfully", "exists", exists)
	return exists, nil
}
