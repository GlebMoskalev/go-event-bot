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
	var total int
	var err error

	if pagerType == "next" {
		nextPage = currentPage + 1
		if nextPage > maxPage {
			return tgbotapi.NewCallback(query.ID, "Это последняя страница")
		}
		schedules, total, err = c.scheduleService.GetAll(ctx, currentPage*5, 5)
	} else if pagerType == "prev" {
		nextPage = currentPage - 1
		if nextPage < 1 {
			return tgbotapi.NewCallback(query.ID, "Это первая страница")
		}
		schedules, total, err = c.scheduleService.GetAll(ctx, currentPage*5-5, 5)
	}

	if err != nil {
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}
	paginationButtons := []models.CallbackButton{
		models.PaginationSchedule(nextPage, total/5, models.Prev),
		models.PageNumber(nextPage, total/5),
		models.PaginationSchedule(nextPage, total/5, models.Next),
	}
	msg := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID,
		keyboards.ScheduleInline(schedules, paginationButtons))
	return msg
}
