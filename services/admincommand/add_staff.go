package admincommand

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *adminCmd) AddStaff(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := logger.SetupLogger(a.log, "admin_command", "AddStaff", "chat_id", msg.ChatID)
	log.Info("processing request to add staff")

	err := a.state.StartAddStaff(ctx, msg.ChatID)
	if err != nil {
		log.Error("failed to start staff addition process", "error", err)
		msg.Text = "Произошла ошибка"
	} else {
		log.Info("staff addition process started successfully")
		msg.Text = "Введите номер телефона в формате:\n 79137777777"
	}
	return msg
}
