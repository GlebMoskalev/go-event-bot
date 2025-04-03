package services

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Start(ctx context.Context, msg tgbotapi.MessageConfig, telegramID int64) (tgbotapi.MessageConfig, error)
	Event(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig
}
