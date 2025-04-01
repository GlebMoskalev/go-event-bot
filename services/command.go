package services

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Start(ctx context.Context, msg tgbotapi.MessageConfig, telegramID int64) tgbotapi.MessageConfig
	GetMenuCommands(role models.Role) []tgbotapi.BotCommand
	Schedule(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig
}
