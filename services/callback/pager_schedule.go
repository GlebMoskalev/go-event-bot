package callback

import (
	"context"
	"fmt"
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
			pages := make(map[string]string)
			pages["Prev"] = fmt.Sprintf("pager:schedule:prev:%d:%d", nextPage, total/5)
			pages[fmt.Sprintf("%d / %d", nextPage, total/5)] = "number:pages"
			if nextPage < maxPage {
				pages["Next"] = fmt.Sprintf("pager:schedule:next:%d:%d", nextPage, total/5)
			}

			msg := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID,
				keyboards.ScheduleInline(schedules, pages))
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
			pages := make(map[string]string)
			if nextPage > 1 {
				pages["Prev"] = fmt.Sprintf("pager:schedule:next:%d:%d", nextPage, total/5)
			}
			pages[fmt.Sprintf("%d / %d", nextPage, total/5)] = "number:pages"
			pages["Next"] = fmt.Sprintf("pager:schedule:next:%d:%d", nextPage, total/5)

			msg := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID,
				keyboards.ScheduleInline(schedules, pages))
			return msg
		}
	}
	return tgbotapi.NewMessage(query.Message.Chat.ID, "Нет доступных страниц")
}
