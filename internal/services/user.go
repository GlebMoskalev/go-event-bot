package services

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/repository"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	"log/slog"
)

type UserService interface {
	Get(ctx context.Context, telegramID int64) (models.User, error)
	Create(ctx context.Context, user models.User) error
	ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error)
}

type userService struct {
	repo repository.UserRepository
	log  *slog.Logger
}

func NewUserService(repo repository.UserRepository, log *slog.Logger) UserService {
	return &userService{repo: repo, log: log}
}

func (s *userService) Get(ctx context.Context, telegramID int64) (models.User, error) {
	log := s.log.With("layer", "service_user", "operation", "Get", "telegram_id", telegramID)
	log.Info("getting user by telegram id")

	user, err := s.repo.Get(ctx, telegramID)
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

func (s *userService) Create(ctx context.Context, user models.User) error {
	log := s.log.With("layer", "service_user", "operation", "Create", "telegram_id", user.TelegramID)
	log.Info("creating user")

	err := s.repo.Create(ctx, user)
	if err != nil {
		log.Error("failed to created user in repository", "error", err)
		return err
	}
	return nil
}

func (s *userService) ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	log := s.log.With("layer", "service_user", "operation", "ExistsUserByTelegramID", "telegram_id", telegramID)
	log.Info("checking user existence by telegram_id")

	exists, err := s.repo.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Error("failed to check user in repository by telegram_id", "error", err)
		return false, err
	}
	return exists, nil
}
