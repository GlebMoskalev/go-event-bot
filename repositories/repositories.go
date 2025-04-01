package repositories

import (
	"context"
	models2 "github.com/GlebMoskalev/go-event-bot/models"
)

type DB interface {
	User
	Schedule
	Staff
	Close() error
}

type User interface {
	GetUser(ctx context.Context, telegramID int64) (models2.User, error)
	CreateUser(ctx context.Context, user models2.User) error
	ExistsUserByTelegramID(ctx context.Context, telegramID int64) (bool, error)
	IsAdmin(ctx context.Context, telegramID int64) (bool, error)
}

type Schedule interface {
	GetAllSchedules(ctx context.Context, offset, limit int) ([]models2.Schedule, int, error)
	UpdateSchedule(ctx context.Context, schedule models2.Schedule) error
	CreateSchedule(ctx context.Context, schedule models2.Schedule) error
	DeleteSchedule(ctx context.Context, scheduleId int) error
}

type Staff interface {
	GetStaffByPhoneNumber(ctx context.Context, phoneNumber string) (models2.Staff, error)
}
