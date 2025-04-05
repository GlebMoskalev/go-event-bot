package command

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Start(ctx context.Context, msg tgbotapi.MessageConfig, telegramID int64) (tgbotapi.MessageConfig, error) {
	log := logger.SetupLogger(c.log,
		"service_command", "Start",
		"chat_id", msg.ChatID,
		"telegram_id", telegramID,
	)
	log.Info("processing start command")

	existUser, err := c.userService.ExistsUserByTelegramID(ctx, telegramID)
	if err != nil {
		log.Error("failed to check user existence", "error", err)
		msg.Text = messages.Error()
		return msg, err
	}

	if existUser {
		user, err := c.userService.Get(ctx, telegramID)
		if err != nil {
			log.Error("failed to get user", "error", err)
			msg.Text = messages.Error()
			return msg, err
		}
		msg.Text = messages.Welcome(user.FirstName, user.Patronymic)

		log.Info("user exists, sending welcome message")
		return msg, apperrors.ErrUserExists
	}

	msg.Text = messages.RequestContact()
	msg.ReplyMarkup = keyboards.ContactButton()

	log.Info("new user, requesting contact")
	return msg, nil
}
