package user

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/apperrors"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type user struct {
	db  repositories.DB
	log *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) services.User {
	return &user{db: db, log: log}
}

func (u *user) Get(ctx context.Context, telegramID int64) (models.User, error) {
	log := u.log.With("layer", "service_user", "operation", "Get", "telegram_id", telegramID)
	log.Info("getting user by telegram id")

	user, err := u.db.GetUser(ctx, telegramID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundUser) {
			log.Warn("user not found in repository")
			return models.User{}, err
		}
		log.Error("failed to to retrieve user from repository", "error", err)
		return models.User{}, err
	}

	log.Info("user retrieved successfully")
	return user, err
}

func (u *user) Create(ctx context.Context, user models.User) error {
	log := u.log.With("layer", "service_user", "operation", "Create", "telegram_id", user.TelegramID)
	log.Info("creating user")

	err := u.db.CreateUser(ctx, user)
	if err != nil {
		log.Error("failed to created user in repository", "error", err)
		return err
	}
	return nil
}

func (u *user) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := u.log.With("layer", "service_user", "operation", "ExistsUserByTelegramID", "telegram_id", telegramID)
	log.Info("checking user existence by telegram_id")

	exists, err := u.db.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Error("failed to check user in repository by telegram_id", "error", err)
		return false, err
	}
	return exists, nil
}

func (u *user) HasRole(ctx context.Context, telegramID int64, role models.Role) (bool, error) {
	log := u.log.With("layer", "service_user", "operation", "HasRole", "telegram_id", telegramID, "role", role)
	log.Info("checking user role by telegram_id")

	if role == models.RoleGuest {
		return true, nil
	}

	usr, err := u.db.GetUser(ctx, telegramID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundUser) {
			log.Warn("user not found")
			return false, err
		}
		log.Warn("failed to get user in repository by telegram_id", "error", err)
		return false, err
	}

	return usr.HasRole(role), err
}
