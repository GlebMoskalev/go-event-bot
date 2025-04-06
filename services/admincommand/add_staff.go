package admincommand

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *adminCmd) AddStaff(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(a.log, "admin_command", "AddStaff", "chat_id", msg.ChatID)
	log.Info("processing request to add staff")

	err := a.state.StartAddStaff(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to start staff addition process", "error", err)
		msg.Text = messages.Error()
	} else {
		log.Info("staff addition process started successfully")
		msg.Text = messages.RequestPhoneNumber()
		msg.ReplyMarkup = keyboards.CancelAddStaff()
	}
	return msg
}
