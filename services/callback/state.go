package callback

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/apperrors"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *callback) CancelAddStaff(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	log := logger.SetupLogger(c.log,
		"service_callback", "CancelAddStaff",
		"chat_id", query.Message.Chat.ID,
		"query_id", query.ID,
	)
	log.Info("starting cancellation of staff addition process")

	_, err := c.stateService.Get(ctx, query.Message.Chat.ID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFoundState) {
			log.Warn("state not found")
			return tgbotapi.NewCallback(query.ID, messages.StaffAdditionMissing())
		}
		log.Error("failed to get state", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}

	err = c.stateService.RemoveState(ctx, query.Message.Chat.ID)
	if err != nil {
		log.Error("failed to remove state", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}

	log.Info("staff addition cancelled successfully")
	return tgbotapi.NewEditMessageTextAndMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		messages.StaffAdditionCancelled(),
		keyboards.EmptyInlineKeyboard(),
	)
}

func (c *callback) AppendStaff(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	log := logger.SetupLogger(c.log,
		"service_callback", "AppendStaff",
		"chat_id", query.Message.Chat.ID,
		"query_id", query.ID,
	)
	log.Info("starting staff addition confirmation process")

	_, data, err := c.stateService.GetWithData(ctx, query.Message.Chat.ID)
	if err != nil {
		log.Error("failed to retrieve staff data", "error", err)
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}
	var staff models.Staff
	err = json.Unmarshal(data, &staff)
	if err != nil {
		log.Error("failed to unmarshal staff data", "error", err)
		return tgbotapi.NewCallback(query.ID, "Произошла ошибка")
	}

	err = c.staffService.Create(ctx, staff)
	if err != nil {
		log.Error("failed to create staff", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}

	err = c.stateService.ConfirmAddStaff(ctx, query.Message.Chat.ID)
	if err != nil {
		log.Error("failed to reset state", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.Error())
	}

	log.Info("staff member created successfully")
	return tgbotapi.NewEditMessageTextAndMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		messages.StaffAdded(),
		keyboards.EmptyInlineKeyboard(),
	)
}
