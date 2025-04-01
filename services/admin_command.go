package services

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminCommand interface {
	Panel(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig
	ChangeEvent(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig
	AddStaff(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig
}
