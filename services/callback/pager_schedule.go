package callback

import (
	"context"
	"fmt"
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
	if pagerType == "next" {
		nextPage := currentPage + 1
		if nextPage <= maxPage {
			schedules, total, err := c.scheduleService.GetAll(ctx, currentPage*5, 5)
			if err != nil {
				msg := tgbotapi.NewMessage(query.Message.Chat.ID, messages.Error())
				return msg
			}
			var paginationButtons []models.PaginationButton
			paginationButtons = append(paginationButtons, models.PaginationButton{
				Text: "Prev",
				Data: fmt.Sprintf("pager:schedule:prev:%d:%d", nextPage, total/5),
			})
			paginationButtons = append(paginationButtons, models.PaginationButton{
				Text: fmt.Sprintf("%d / %d", nextPage, total/5),
				Data: "number:pages",
			})
			if nextPage < maxPage {
				paginationButtons = append(paginationButtons, models.PaginationButton{
					Text: "Next",
					Data: fmt.Sprintf("pager:schedule:next:%d:%d", nextPage, total/5),
				})
			}

			msg := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID,
				keyboards.ScheduleInline(schedules, paginationButtons))
			return msg
		}
	} else if pagerType == "prev" {
		nextPage := currentPage - 1
		if nextPage > 0 {
			schedules, total, err := c.scheduleService.GetAll(ctx, currentPage*5-5, 5)
			if err != nil {
				msg := tgbotapi.NewMessage(query.Message.Chat.ID, messages.Error())
				return msg
			}

			var paginationButtons []models.PaginationButton

			if nextPage > 1 {
				paginationButtons = append(paginationButtons, models.PaginationButton{
					Text: "Prev",
					Data: fmt.Sprintf("pager:schedule:prev:%d:%d", nextPage, total/5),
				})
			}
			paginationButtons = append(paginationButtons, models.PaginationButton{
				Text: fmt.Sprintf("%d / %d", nextPage, total/5),
				Data: "number:pages",
			})
			paginationButtons = append(paginationButtons, models.PaginationButton{
				Text: "Next",
				Data: fmt.Sprintf("pager:schedule:next:%d:%d", nextPage, total/5),
			})

			msg := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID,
				keyboards.ScheduleInline(schedules, paginationButtons))
			return msg
		}
	}
	return tgbotapi.NewMessage(query.Message.Chat.ID, "Нет доступных страниц")
}
