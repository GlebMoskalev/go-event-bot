package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
)

type Staff interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (models.Staff, error)
	Create(ctx context.Context, staff models.Staff) error
	Update(ctx context.Context, staff models.Staff) error
}
