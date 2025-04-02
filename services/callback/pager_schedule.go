package callback

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (c *callback) PagerSchedule(ctx context.Context, query *tgbotapi.CallbackQuery, data ...string) tgbotapi.Chattable {
	pagerType := data[0]
	currentPage, _ := strconv.Atoi(data[1])
	maxPage, _ := strconv.Atoi(data[2])

	var nextPage int
	var schedules []models.Schedule
	//var total int
	var err error

	if pagerType == "next" {
		nextPage = currentPage + 1
		if nextPage > maxPage {
			return tgbotapi.NewCallback(query.ID, "Это последняя страница")
		}
		schedules, _, err = c.scheduleService.GetAll(ctx, (nextPage-1)*models.ItemsPerPage, models.ItemsPerPage)
	} else if pagerType == "prev" {
		nextPage = currentPage - 1
		if nextPage < 1 {
			return tgbotapi.NewCallback(query.ID, "Это первая страница")
		}
		schedules, _, err = c.scheduleService.GetAll(
			ctx, (nextPage-1)*models.ItemsPerPage, models.ItemsPerPage,
		)
	}
	if err != nil {
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	return tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, messages.AllSchedules(schedules), keyboards.PaginationScheduleInline([]models.CallbackButton{
		models.PaginationSchedule(nextPage, maxPage, models.Prev),
		models.PageNumber(nextPage, maxPage),
		models.PaginationSchedule(nextPage, maxPage, models.Next),
	}))
}
