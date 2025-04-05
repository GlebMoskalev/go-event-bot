package callback

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
)

func (c *callback) EventAll(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	log := logger.SetupLogger(c.log,
		"service_callback", "EventAll",
		"chat_id", query.Message.Chat.ID,
		"query_id", query.ID,
	)
	log.Info("processing request to get all events")

	schedules, total, err := c.scheduleService.GetAll(ctx, 0, models.ItemsPerPage)
	if err != nil {
		log.Error("failed to fetch events", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}

	if len(schedules) == 0 {
		log.Info("no events found")
		return tgbotapi.NewCallback(query.ID, messages.EventEmpty())
	}

	maxPage := int(math.Ceil(float64(total) / float64(models.ItemsPerPage)))

	log.Info("events retrieved successfully",
		"events_count", len(schedules),
		"total", total,
		"max_page", maxPage,
	)
	return tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, messages.AllEvents(schedules), keyboards.PaginationScheduleInline([]models.CallbackButton{
		models.PaginationSchedule(1, maxPage, models.Prev),
		models.PageNumber(1, maxPage),
		models.PaginationSchedule(1, maxPage, models.Next),
	}))
}
