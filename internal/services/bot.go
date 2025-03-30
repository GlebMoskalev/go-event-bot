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
	Start(update tgbotapi.Update)
	RequestContact(update tgbotapi.Update)
	CheckUser(update tgbotapi.Update) bool
	SetBot(bot *tgbotapi.BotAPI)
}

type botService struct {
	log          *slog.Logger
	staffService StaffService
	userService  UserService
	api          *tgbotapi.BotAPI
}

func NewBotService(staffService StaffService, userService UserService, log *slog.Logger) BotService {
	return &botService{staffService: staffService, userService: userService, log: log}
}

func (b *botService) SetBot(bot *tgbotapi.BotAPI) {
	b.api = bot
}

func (b *botService) sendErrorNotification(chatID int64, log *slog.Logger) {
	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать некорректно")
	b.sendMessage(msg, log)
}

func (b *botService) sendMessage(msg tgbotapi.MessageConfig, log *slog.Logger) {
	_, err := b.api.Send(msg)
	if err != nil {
		log.Error("failed to send message", "error", err)
	}
}

func (b *botService) CheckUser(update tgbotapi.Update) bool {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "CheckUser", "chat_id", chatID)
	log.Info("check user")

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		b.sendErrorNotification(chatID, log)
		return false
	}

	return exists
}

func (b *botService) Start(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)
	log.Info("start")

	if b.CheckUser(update) {
		user, err := b.userService.Get(context.Background(), update.Message.From.ID)
		if err != nil {
			log.Error("failed to get user")
			b.sendErrorNotification(chatID, log)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Привет, %s %s", user.FirstName, user.Patronymic))
		b.sendMessage(msg, log)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Привет! Для дальнейшей работы нужен твой контакт!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
	})
	b.sendMessage(msg, log)
}

func (b *botService) RequestContact(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		b.sendErrorNotification(chatID, log)
		return
	}

	if exists {
		msg := tgbotapi.NewMessage(chatID, "Твой контакт у нас уже есть")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		b.sendMessage(msg, log)
		return
	}

	phoneNumber := update.Message.Contact.PhoneNumber
	staff, err := b.staffService.GetByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Ваш номер телефона не соответствует нашему формату")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			b.sendMessage(msg, log)
			return
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			msg := tgbotapi.NewMessage(chatID, "Вас нет в списках, попросите администратора добавить вас")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			b.sendMessage(msg, log)
			return
		}

		log.Error("failed to get staff in service", "error", err)

		b.sendErrorNotification(chatID, log)
		return
	}

	err = b.userService.Create(context.Background(), models.User{
		TelegramID: update.Message.Contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
	})

	if err != nil {
		b.sendErrorNotification(chatID, log)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Приветствуем %s %s на нашем мероприятии!",
		staff.FirstName,
		staff.Patronymic),
	)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	b.sendMessage(msg, log)
	return
}
