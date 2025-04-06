package services

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Callback interface {
	PagerEvent(ctx context.Context, query *tgbotapi.CallbackQuery, data ...string) tgbotapi.Chattable
	EventAll(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
	CancelAddStaff(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
	AppendStaff(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
	SearchStaffByLastName(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
	SearchStaffByPhoneNumber(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable
}
