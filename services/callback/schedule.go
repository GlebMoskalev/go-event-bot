package callback

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
)

func (c *callback) ScheduleAll(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	schedules, total, err := c.scheduleService.GetAll(ctx, 0, models.ItemsPerPage)
	if err != nil {
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}
	if len(schedules) == 0 {
		return tgbotapi.NewCallback(query.ID, messages.ScheduleEmpty())
	}

	maxPage := int(math.Ceil(float64(total) / float64(models.ItemsPerPage)))

	return tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, messages.AllSchedules(schedules), keyboards.PaginationScheduleInline([]models.CallbackButton{
		models.PaginationSchedule(1, maxPage, models.Prev),
		models.PageNumber(1, maxPage),
		models.PaginationSchedule(1, maxPage, models.Next),
	}))
}
