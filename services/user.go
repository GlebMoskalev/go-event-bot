package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type User interface {
	Get(ctx context.Context, telegramID int64) (models.User, error)
	Create(ctx context.Context, user models.User) error
	ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error)
	IsAdmin(ctx context.Context, telegramID int64) (bool, error)
}
