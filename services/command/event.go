package command

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Event(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(c.log, "service_command", "Event", "chat_id", msg.ChatID)
	log.Info("processing event command")

	msg.Text = messages.EventTitle()
	msg.ReplyMarkup = keyboards.EventInline()

	log.Info("event message prepared successfully")
	return msg
}
