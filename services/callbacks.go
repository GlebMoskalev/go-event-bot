package services

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Callback interface {
	PagerSchedule(ctx context.Context, query *tgbotapi.CallbackQuery, data ...string) tgbotapi.Chattable
	ScheduleAll(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
}
