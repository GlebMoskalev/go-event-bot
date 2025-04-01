package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type Schedule interface {
	GetAll(ctx context.Context, offset, limit int) ([]models.Schedule, int, error)
	Update(ctx context.Context, schedule models.Schedule) error
	Create(ctx context.Context, schedule models.Schedule) error
	Delete(ctx context.Context, scheduleId int) error
}
