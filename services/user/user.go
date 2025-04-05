package user

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
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
	log := logger.SetupLogger(u.log, "service_user", "Get", "telegram_id", telegramID)
	log.Info("getting user")

	user, err := u.db.GetUser(ctx, telegramID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundUser) {
			log.Warn("user not found")
			return models.User{}, err
		}
		log.Error("failed to retrieve user", "error", err)
		return models.User{}, err
	}

	log.Info("user retrieved successfully")
	return user, err
}

func (u *user) Create(ctx context.Context, user models.User) error {
	log := logger.SetupLogger(u.log, "service_user", "Create", "telegram_id", user.TelegramID)
	log.Info("creating user")

	err := u.db.CreateUser(ctx, user)
	if err != nil {
		log.Error("failed to created user", "error", err)
		return err
	}

	log.Info("user created successfully")
	return nil
}

func (u *user) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := logger.SetupLogger(u.log, "service_user", "ExistsUserByTelegramID", "telegram_id", telegramID)
	log.Info("checking user existence")

	exists, err := u.db.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Error("failed to check user existence", "error", err)
		return false, err
	}

	log.Info("user existence checked successfully", "exists", exists)
	return exists, nil
}

func (u *user) HasRole(ctx context.Context, telegramID int64, role models.Role) (bool, error) {
	log := logger.SetupLogger(u.log,
		"service_user", "HasRole",
		"telegram_id", telegramID,
		"role", role,
	)
	log.Info("checking user role")

	if role == models.RoleGuest {
		log.Info("guest role detected, returning true")
		return true, nil
	}

	usr, err := u.db.GetUser(ctx, telegramID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundUser) {
			log.Warn("user not found")
			return false, err
		}
		log.Warn("failed to get user", "error", err)
		return false, err
	}

	log.Info("user role checked successfully")
	return usr.HasRole(role), err
}
