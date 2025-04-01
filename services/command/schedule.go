package command

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Schedule(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := c.log.With("layer", "service_command", "operation", "Schedule", "chat_id", msg.ChatID)
	schedules, _, err := c.scheduleService.GetAll(ctx, 0, 5)
	if err != nil {
		log.Error("failed to get schedules")
		msg.Text = "Произошла ошибка!"
		return msg
	}

	if len(schedules) == 0 {
		msg.Text = "Расписание отсутствует"
	} else {
		msg.Text = "Расписание:"
		msg.ReplyMarkup = keyboards.ScheduleInline(schedules)
	}

	return msg
}
