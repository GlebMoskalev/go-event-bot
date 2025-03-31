package handlers

import (
	"github.com/GlebMoskalev/go-event-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type AdminHandler struct {
	botService services.BotService
	log        *slog.Logger
}

func NewAdminHandler(botService services.BotService, log *slog.Logger) *AdminHandler {
	return &AdminHandler{botService: botService, log: log}
}

func (h *AdminHandler) CheckAdmin(update tgbotapi.Update) bool {
	isAdmin, err := h.botService.IsAdmin(update)
	if err != nil {
		h.log.Error("failed to check admin status", "error", err)
		return false
	}

	if !isAdmin {
		h.botService.SendMessage(update.Message.Chat.ID, "У вас нет доступа к этой команде", nil)
		return false
	}

	return true
}

func (h *AdminHandler) HandleAdmin(update tgbotapi.Update) {
	if !h.CheckAdmin(update) {
		return
	}

	h.botService.SendMessage(update.Message.Chat.ID, "Панель администратора", nil)
}

func (h *AdminHandler) HandleChangeEvent(update tgbotapi.Update) {
	if !h.CheckAdmin(update) {
		return
	}

	h.botService.SendMessage(update.Message.Chat.ID, "Изменение мероприятия", nil)
}

func (h *AdminHandler) HandleAddStaff(update tgbotapi.Update) {
	if !h.CheckAdmin(update) {
		return
	}

	h.botService.SendMessage(update.Message.Chat.ID, "Изменение мероприятия", nil)
}
