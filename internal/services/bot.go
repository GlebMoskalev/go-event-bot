package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type BotService interface {
	Start(bot *tgbotapi.BotAPI, update tgbotapi.Update)
	RequestContact(bot *tgbotapi.BotAPI, update tgbotapi.Update)
	CheckUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool
}

type botService struct {
	log          *slog.Logger
	staffService StaffService
	userService  UserService
}

func NewBotService(staffService StaffService, userService UserService, log *slog.Logger) BotService {
	return &botService{staffService: staffService, userService: userService, log: log}
}

func sendErrorNotification(bot *tgbotapi.BotAPI, chatID int64, log *slog.Logger) {
	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать некорректно")
	sendMessage(bot, msg, log)
}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, log *slog.Logger) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Error("failed to send message", "error", err)
	}
}

func (b *botService) CheckUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "CheckUser", "chat_id", chatID)
	log.Info("check user")

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		sendErrorNotification(bot, chatID, log)
		return false
	}

	return exists
}

func (b *botService) Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)
	log.Info("start")

	if b.CheckUser(bot, update) {
		user, err := b.userService.Get(context.Background(), update.Message.From.ID)
		if err != nil {
			log.Error("failed to get user")
			sendErrorNotification(bot, chatID, log)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Привет, %s %s", user.FirstName, user.Patronymic))
		sendMessage(bot, msg, log)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Привет! Для дальнейшей работы нужен твой контакт!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
	})
	sendMessage(bot, msg, log)
}

func (b *botService) RequestContact(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		sendErrorNotification(bot, chatID, log)
		return
	}

	if exists {
		msg := tgbotapi.NewMessage(chatID, "Твой контакт у нас уже есть")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		sendMessage(bot, msg, log)
		return
	}

	phoneNumber := update.Message.Contact.PhoneNumber
	staff, err := b.staffService.GetByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Ваш номер телефона не соответствует нашему формату")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			sendMessage(bot, msg, log)
			return
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Вас нет в списках, попросите администратора добавить вас")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			sendMessage(bot, msg, log)
			return
		}

		log.Error("failed to get staff in service", "error", err)

		sendErrorNotification(bot, chatID, log)
		return
	}

	err = b.userService.Create(context.Background(), models.User{
		TelegramID: update.Message.Contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
	})

	if err != nil {
		sendErrorNotification(bot, chatID, log)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Приветствуем на нашем мероприятии!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	sendMessage(bot, msg, log)
	return
}
