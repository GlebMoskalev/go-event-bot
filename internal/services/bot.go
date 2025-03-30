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
	CheckUser(update tgbotapi.Update) (bool, error)
	SendMessage(chatID int64, text string, markup any)
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

func (b *botService) sendErrorNotification(chatID int64) {
	b.SendMessage(chatID, "Произошла ошибка! Бот может работать некорректно", nil)
}

func (b *botService) SendMessage(chatID int64, text string, markup any) {
	msg := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	_, err := b.api.Send(msg)
	if err != nil {
		b.log.Error("failed to send message",
			"chat_id", chatID,
			"text", text,
			"error", err)
		return
	}
}

func (b *botService) CheckUser(update tgbotapi.Update) (bool, error) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "CheckUser", "chat_id", chatID)
	log.Info("check user")

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		b.sendErrorNotification(chatID)
		return false, err
	}

	return exists, nil
}

func (b *botService) Start(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)
	log.Info("start")

	existUser, err := b.CheckUser(update)
	if err != nil {
		return
	}

	if existUser {
		user, err := b.userService.Get(context.Background(), update.Message.From.ID)
		if err != nil {
			log.Error("failed to get user")
			b.sendErrorNotification(chatID)
			return
		}
		greeting := fmt.Sprintf("Привет, %s %s", user.FirstName, user.Patronymic)
		b.SendMessage(chatID, greeting, nil)
		return
	}

	contactButton := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
	})
	b.SendMessage(chatID, "Привет! Для дальнейшей работы нужен твой контакт!", contactButton)
}

func (b *botService) RequestContact(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	log := b.log.With("layer", "service_bot", "operation", "RequestContact", "chat_id", chatID)

	exists, err := b.userService.ExistsUserByTelegramID(context.Background(), update.Message.From.ID)
	if err != nil {
		log.Error("failed to check user in service", "error", err)
		b.sendErrorNotification(chatID)
		return
	}

	if exists {
		removeKeyboard := tgbotapi.NewRemoveKeyboard(false)
		b.SendMessage(chatID, "Твой контакт у нас уже есть", removeKeyboard)
		return
	}

	phoneNumber := update.Message.Contact.PhoneNumber
	staff, err := b.staffService.GetByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPhoneNumber) {
			log.Error("invalid phone number", "phone_number", phoneNumber)
			removeKeyboard := tgbotapi.NewRemoveKeyboard(false)
			b.SendMessage(chatID, "Ваш номер телефона не соответствует нашему формату", removeKeyboard)
			return
		}

		if errors.Is(err, apperrors.ErrNotFoundStaff) {
			log.Warn("user not found in staff", "phone_number", phoneNumber)
			removeKeyboard := tgbotapi.NewRemoveKeyboard(false)
			b.SendMessage(chatID, "Вас нет в списках, попросите администратора добавить вас", removeKeyboard)
			return
		}

		log.Error("failed to get staff in service", "error", err)

		b.sendErrorNotification(chatID)
		return
	}

	err = b.userService.Create(context.Background(), models.User{
		TelegramID: update.Message.Contact.UserID,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Patronymic: staff.Patronymic,
	})

	if err != nil {
		b.sendErrorNotification(chatID)
		return
	}

	greeting := fmt.Sprintf("%s %s, приветствуем на нашем мероприятии!",
		staff.FirstName,
		staff.Patronymic)
	removeKeyboard := tgbotapi.NewRemoveKeyboard(false)
	b.SendMessage(chatID, greeting, removeKeyboard)
}
