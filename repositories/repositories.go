package repositories

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type DB interface {
	User
	Event
	Staff
	State
	Close() error
}

type User interface {
	GetUser(ctx context.Context, telegramID int64) (models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error)
}

type Event interface {
	GetAllEvents(ctx context.Context, offset, limit int) ([]models.Event, int, error)
	UpdateEvent(ctx context.Context, event models.Event) error
	CreateEvent(ctx context.Context, event models.Event) error
	DeleteEvent(ctx context.Context, scheduleId int) error
}

type Staff interface {
	GetStaffByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error)
}

type State interface {
	GetState(ctx context.Context, chatID int64) (models.State, error)
	GetStateAndData(ctx context.Context, chatID int64) (models.State, []byte, error)
	SetState(ctx context.Context, chatID int64, state models.State, dataJSON []byte) error
	RemoveState(ctx context.Context, chatID int64) error
}
