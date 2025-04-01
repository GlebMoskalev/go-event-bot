package services

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message interface {
	Contact(ctx context.Context, msg tgbotapi.MessageConfig, contact *tgbotapi.Contact) (tgbotapi.MessageConfig, error)
}
