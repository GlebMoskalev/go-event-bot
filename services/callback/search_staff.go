package callback

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/utils/keyboards"
	"github.com/GlebMoskalev/go-event-bot/utils/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *callback) SearchStaffByLastName(ctx context.Context, query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	log := logger.SetupLogger(c.log,
		"service_callback", "SearchStaffByLastName",
		"chat_id", query.Message.Chat.ID,
		"query_id", query.ID,
	)
	log.Info("initiating staff search by last name")

	err := c.stateService.StartSearchByLastName(ctx, query.Message.Chat.ID)
	if err != nil {
		log.Error("failed to set state", "error", err)
		return tgbotapi.NewCallback(query.ID, messages.StaffAdditionMissing())
	}

	log.Info("state set successfully, requesting last name")
	return tgbotapi.NewEditMessageTextAndMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		messages.RequestFullName(),
		keyboards.EmptyInlineKeyboard(),
	)
}
