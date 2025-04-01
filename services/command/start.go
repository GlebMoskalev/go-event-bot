package command

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	messages2 "github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Start(ctx context.Context, msg tgbotapi.MessageConfig, telegramID int64) tgbotapi.MessageConfig {
	chatID := msg.ChatID
	log := c.log.With("layer", "service_command", "operation", "Start", "chat_id", chatID)
	log.Info("start")

	existUser, err := c.userService.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		c.log.Error("failed to check user")
		msg.Text = messages2.Error()
		return msg
	}

	if existUser {
		user, err := c.userService.Get(ctx, telegramID)
		if err != nil {
			log.Error("failed to get user")
			msg.Text = messages2.Error()
			return msg
		}
		msg.Text = messages2.Welcome(user.FirstName, user.Patronymic)
		return msg
	}

	msg.Text = messages2.RequestContact()
	msg.ReplyMarkup = keyboards.ContactButton()
	return msg
}
