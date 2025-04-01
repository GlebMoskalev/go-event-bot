package command

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) AdminPanel(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.Text = "Админ панель"
	return msg
}
