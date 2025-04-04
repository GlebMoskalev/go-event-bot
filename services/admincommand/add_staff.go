package admincommand

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *adminCmd) AddStaff(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	err := a.state.StartAddStaff(ctx, msg.ChatID)
	if err != nil {
		msg.Text = "Произошла ошибка"
	} else {
		msg.Text = "Введите ФИО сотрудника в формате: Иван Иванович Иванов"
	}
	return msg
}
