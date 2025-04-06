package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type State interface {
	Get(ctx context.Context, chatID int64) (models.State, error)
	GetWithData(ctx context.Context, chatID int64) (models.State, []byte, error)
	RemoveState(ctx context.Context, chatID int64) error

	StartAddStaff(ctx context.Context, chatID int64) error
	RegisterStaffFullName(ctx context.Context, chatID int64, firstName, lastName, patronymic string) error
	RegisterStaffNumberPhone(ctx context.Context, chatID int64, numberPhone string) error
	ConfirmAddStaff(ctx context.Context, chatID int64) error

	StartSearchByLastName(ctx context.Context, chatID int64) error
}
