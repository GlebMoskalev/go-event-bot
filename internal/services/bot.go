package services

import (
	"context"
	"errors"
	"github.com/GlebMoskalev/go-event-bot/internal/models"
	"github.com/GlebMoskalev/go-event-bot/internal/utils/apperrors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type BotService interface {
	Start(bot *tgbotapi.BotAPI, update tgbotapi.Update)
	RequestContact(bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

type botService struct {
	log          *slog.Logger
	staffService StaffService
	userService  UserService
}

func NewBotService(staffService StaffService, userService UserService, log *slog.Logger) BotService {
	return &botService{staffService: staffService, userService: userService, log: log}
}

func (b *botService) Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)
	log.Info("sta")

	msg := tgbotapi.NewMessage(chatID, "Привет! Для дальнейшей работы нужен твой контакт!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
	})
	_, err := bot.Send(msg)

	if err != nil {
		log.Error("failed to message", "command", "Start", "chat_id", chatID)
	}
	//b.activeContactRequests = append(b.activeContactRequests, chatID)
}

func (b *botService) RequestContact(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)

	exists, err := b.userService.ExistsUserByChatID(context.Background(), chatID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать некорректно")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error("failed to message")
		}
		return
	}

	if exists {
		msg := tgbotapi.NewMessage(chatID, "Твой контакт у нас уже есть")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		_, err = bot.Send(msg)
		if err != nil {
			log.Error("failed to message")
		}
		return
	}

	phoneNumber := update.Message.Contact.PhoneNumber
	staff, err := b.staffService.GetByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Ваш номер телефона не соответствует нашему формату")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			_, err = bot.Send(msg)
			if err != nil {
				log.Error("failed to message")
			}
			return
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Вас нет в списках, попросите администратора добавить вас")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			_, err = bot.Send(msg)
			if err != nil {
				log.Error("failed to message")
			}
			return
		}

		log.Error("failed to get staff in service", "error", err)

		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать некорректно")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error("failed to message")
		}
		return
	}

	err = b.userService.Create(context.Background(), models.User{
		TelegramID: update.Message.Contact.UserID,
		ChatID:     chatID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
	})

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать некорректно")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error("failed to message")
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Приветствуем на нашем мероприятии!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	_, err = bot.Send(msg)
	if err != nil {
		log.Error("failed to message")
	}
	return
}
