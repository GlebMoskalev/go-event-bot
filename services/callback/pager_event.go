package callback

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (c *callback) PagerEvent(ctx context.Context, query *tgbotapi.CallbackQuery, data ...string) tgbotapi.Chattable {
	log := logger.SetupLogger(c.log,
		"service_callback", "PagerEvent",
		"chat_id", query.Message.Chat.ID,
		"query_id", query.ID,
		"pager_type", data[0],
	)
	log.Info("processing pagination request")

	pagerType := data[0]
	currentPage, err := strconv.Atoi(data[1])
	if err != nil {
		log.Error("failed to parse current page", "error", err, "value", data[1])
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	maxPage, err := strconv.Atoi(data[2])
	if err != nil {
		log.Error("failed to parse max page", "error", err, "value", data[2])
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	var nextPage int
	var events []models.Event

	if pagerType == "next" {
		nextPage = currentPage + 1
		if nextPage > maxPage {
			log.Info("attempt to go beyond last page", "current_page", currentPage, "max_page", maxPage)
			return tgbotapi.NewCallback(query.ID, "Это последняя страница")
		}
		events, _, err = c.eventService.GetAll(ctx, (nextPage-1)*models.ItemsPerPage, models.ItemsPerPage)
	} else if pagerType == "prev" {
		nextPage = currentPage - 1
		if nextPage < 1 {
			log.Info("attempt to go before first page", "current_page", currentPage)
			return tgbotapi.NewCallback(query.ID, "Это первая страница")
		}
		events, _, err = c.eventService.GetAll(
			ctx, (nextPage-1)*models.ItemsPerPage, models.ItemsPerPage,
		)
	} else {
		log.Warn("unknown pager type", "pager_type", pagerType)
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	if err != nil {
		log.Error("failed to fetch events", "error", err, "page", nextPage)
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	log.Info("pagination processed successfully")
	return tgbotapi.NewEditMessageTextAndMarkup(
		query.Message.Chat.ID, query.Message.MessageID, messages.AllEvents(events),
		keyboards.PaginationEventInline([]models.CallbackButton{
			models.PaginationEvent(nextPage, maxPage, models.Prev),
			models.PageNumber(nextPage, maxPage),
			models.PaginationEvent(nextPage, maxPage, models.Next),
		}))
}
