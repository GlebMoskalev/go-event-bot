package command

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	messages2 "github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *cmd) Schedule(ctx context.Context, msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	log := c.log.With("layer", "service_command", "operation", "Schedule", "chat_id", msg.ChatID)
	schedules, total, err := c.scheduleService.GetAll(ctx, 0, 5)
	if err != nil {
		log.Error("failed to get schedules")
		msg.Text = messages2.Error()
		return msg
	}

	if len(schedules) == 0 {
		msg.Text = messages2.ScheduleEmpty()
	} else {
		msg.Text = messages2.ScheduleTitle()
		msg.ReplyMarkup = keyboards.ScheduleInline(
			schedules, []models.CallbackButton{
				models.PaginationSchedule(1, total/5, models.Prev),
				models.PageNumber(1, total/5),
				models.PaginationSchedule(1, total/5, models.Next),
			},
		)
	}

	return msg
}
