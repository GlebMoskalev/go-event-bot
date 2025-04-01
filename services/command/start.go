package command

import (
	"context"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/pkg/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Start(ctx context.Context, msg tgbotapi.MessageConfig, telegramID int64) tgbotapi.MessageConfig {
	chatID := msg.ChatID
	log := c.log.With("layer", "service_command", "operation", "Start", "chat_id", chatID)
	log.Info("start")

	existUser, err := c.userService.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		c.log.Error("failed to check user")
		msg.Text = "Произошла ошибка!"
		return msg
	}

	if existUser {
		user, err := c.userService.Get(ctx, telegramID)
		if err != nil {
			log.Error("failed to get user")
			msg.Text = "Произошла ошибка!"
			return msg
		}
		msg.Text = fmt.Sprintf("Привет, %s %s", user.FirstName, user.Patronymic)
		return msg
	}

	msg.Text = "Привет! Для дальнейшей работы нужен твой контакт!"
	msg.ReplyMarkup = keyboards.ContactButton()
	return msg
}
