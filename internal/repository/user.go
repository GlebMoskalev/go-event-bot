package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	"log/slog"
)

type UserRepository interface {
	Get(ctx context.Context, telegramID int64) (models.User, error)
	Create(ctx context.Context, user models.User) error
	ExistsUserByChatID(ctx context.Context, chatID int64) (bool, error)
}

type userRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewUserRepository(db *sql.DB, log *slog.Logger) UserRepository {
	return &userRepo{db: db, log: log}
}

func (r *userRepo) Get(ctx context.Context, telegramID int64) (models.User, error) {
	log := r.log.With("layer", "repository_user", "operation", "Get", "telegram_id", telegramID)
	log.Info("fetching user data")

	query := `
	SELECT
		chat_id,
		firstname,
		lastname
	FROM users
	WHERE telegram_id = $1
`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ChatID,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user not found")
			return models.User{}, apperrors.ErrNotFoundStaff
		}

		log.Error("failed to fetch user from database", "error", err)
	}

	log.Info("user retrieved successfully")
	return user, err
}

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	log := r.log.With("layer", "repository_user", "operation", "Create", "telegram_id", user.TelegramID)
	log.Info("creating new user")

	query := `
	INSERT INTO users 
		(telegram_id, chat_id, firstname, lastname) 
	VALUES 
		($1, $2, $3, $4)
`

	_, err := r.db.ExecContext(ctx, query, user.TelegramID, user.ChatID, user.FirstName, user.LastName)
	if err != nil {
		log.Error("failed to create user in database", "error", err)
		return err
	}

	log.Info("user created successfully")
	return nil
}

func (r *userRepo) ExistsUserByChatID(ctx context.Context, chatID int64) (bool, error) {
	log := r.log.With("layer", "repository_user", "operation", "FindByChatId", "chat_id", chatID)
	log.Info("checking user existence by chat_id")

	query := `
	SELECT 
	EXISTS(
		SELECT 1 FROM users WHERE chat_id = $1
	)
`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&exists)

	if err != nil {
		log.Error("failed to check user existence", "error", err)
		return false, err
	}

	log.Info("user existence checked successfully")
	return exists, err
}
