package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message interface {
	Contact(ctx context.Context, msg tgbotapi.MessageConfig,
		contact *tgbotapi.Contact) (tgbotapi.MessageConfig, models.Role, error)
	State(ctx context.Context, msg tgbotapi.MessageConfig, state models.State) tgbotapi.MessageConfig
}
