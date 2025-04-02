package keyboards

import (
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
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

func ScheduleInline(schedules []models.Schedule, paginationButtons []models.CallbackButton) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, s := range schedules {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(s.Title, strconv.Itoa(s.ID)),
		}
		rows = append(rows, row)
	}
	var extraRow []tgbotapi.InlineKeyboardButton
	for _, button := range paginationButtons {
		extraRow = append(extraRow, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
	}
	rows = append(rows, extraRow)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
