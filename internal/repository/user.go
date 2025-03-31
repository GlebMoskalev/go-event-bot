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
	ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error)
	IsAdmin(ctx context.Context, telegramID int64) (bool, error)
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
		firstname,
		lastname,
		patronymic,
		is_admin
	FROM users
	WHERE telegram_id = $1
`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.IsAdmin,
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

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	log := r.log.With("layer", "repository_user", "operation", "Create", "telegram_id", user.TelegramID)
	log.Info("creating new user")

	query := `
	INSERT INTO users 
		(telegram_id, firstname, lastname, patronymic, is_admin) 
	VALUES 
		($1, $2, $3, $4, $5)
`

	_, err := r.db.ExecContext(ctx, query, user.TelegramID, user.FirstName, user.LastName, user.Patronymic, user.IsAdmin)
	if err != nil {
		log.Error("failed to create user in database", "error", err)
		return err
	}

	log.Info("user created successfully")
	return nil
}

func (r *userRepo) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := r.log.With("layer", "repository_user", "operation", "ExistsUserByTelegramID", "telegram_id", telegramID)
	log.Info("checking user existence by telegram_id")

	query := `
	SELECT 
	EXISTS(
		SELECT 1 FROM users WHERE telegram_id = $1
	)
`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(&exists)

	if err != nil {
		log.Error("failed to verify user existence using telegram_id", "error", err)
		return false, err
	}

	log.Info("user existence by telegram_id checked successfully")
	return exists, err
}

func (r *userRepo) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	log := r.log.With("layer", "repository_user", "operation", "IsAdmin", "telegram_id", telegramID)
	log.Info("checking user admin by telegram_id")

	var isAdmin bool
	query := `
	SELECT 
	    is_admin
	FROM users
	WHERE telegram_id = $1
`

	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(&isAdmin)
	if err != nil {
		log.Error("failed to check user admin by telegram_id", "error", err)
		return false, err
	}
	return isAdmin, nil
}
