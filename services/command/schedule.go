package command

import (
	"context"
	"fmt"
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
			schedules, []models.PaginationButton{
				{
					Text: fmt.Sprintf("1 / %d", total/5),
					Data: "number:pages",
				},
				{
					Text: "Next",
					Data: fmt.Sprintf("pager:schedule:next:%d:%d", 1, total/5),
				},
			},
		)
	}

	return msg
}
