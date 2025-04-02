package keyboards

import (
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ContactButton() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
		},
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(false)
}

func PaginationScheduleInline(paginationButtons []models.CallbackButton) tgbotapi.InlineKeyboardMarkup {
	var extraRow []tgbotapi.InlineKeyboardButton
	for _, button := range paginationButtons {
		extraRow = append(extraRow, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
	}
	return tgbotapi.NewInlineKeyboardMarkup(extraRow)
}

func ScheduleInline() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание всех мероприятий", "schedule:all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание по дням", "schedule:days"),
		),
	)
}
