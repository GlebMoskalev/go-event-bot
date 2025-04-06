package admincommand

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *adminCmd) SearchStaff(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.Text = messages.StaffSearchMethod()
	msg.ReplyMarkup = keyboards.SearchMethodStaff()
	return msg
}
