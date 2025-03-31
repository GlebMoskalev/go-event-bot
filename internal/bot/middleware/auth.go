package middleware

import (
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type AuthMiddleware struct {
	botService services.BotService
	log        *slog.Logger
}

func NewAuthMiddleWare(botService services.BotService, log *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{botService: botService, log: log}
}

func (m *AuthMiddleware) CheckAuth(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	mes := update.Message

	if mes.IsCommand() && mes.Command() == "start" {
		// commands available to everyone
		return true
	}

	if mes.Contact != nil {
		// contact needed from command start
		return true
	}

	existsUser, err := m.botService.CheckUser(update)
	if err != nil {
		m.log.Error("failed to check user", "error", err)
		return false
	}

	if !existsUser {
		m.botService.SendMessage(update.Message.Chat.ID, "Нужно выполнить команду /start", nil)
		return false
	}

	return true
}
