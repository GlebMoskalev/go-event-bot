package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type Event interface {
	GetAll(ctx context.Context, offset, limit int) ([]models.Event, int, error)
	Update(ctx context.Context, event models.Event) error
	Create(ctx context.Context, event models.Event) error
	Delete(ctx context.Context, eventId int) error
}
