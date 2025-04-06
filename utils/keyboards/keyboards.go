package keyboards

import (
	"fmt"
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

func PaginationEventInline(paginationButtons []models.CallbackButton) tgbotapi.InlineKeyboardMarkup {
	var extraRow []tgbotapi.InlineKeyboardButton
	for _, button := range paginationButtons {
		extraRow = append(extraRow, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
	}
	rowBack := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
		"Вернутся обратно", fmt.Sprintf("%s:%s", models.EventContext, models.BackContext)))
	return tgbotapi.NewInlineKeyboardMarkup(extraRow, rowBack)
}

func EventInline() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание всех мероприятий",
				fmt.Sprintf("%s:%s", models.EventContext, models.AllContext)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание по дням",
				fmt.Sprintf("%s:%s", models.EventContext, models.DaysContext)),
		),
	)
}

func CancelAddStaff() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Отменить добавление",
			fmt.Sprintf("%s:%s:%s", models.StaffContext, models.AddContext, models.CancelContext),
		),
	))
}

func EmptyInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
}

func AgreeStaff() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"Добавить",
			fmt.Sprintf("%s:%s:%s", models.StaffContext, models.AddContext, models.AppendContext),
		),
		tgbotapi.NewInlineKeyboardButtonData(
			"Отменить",
			fmt.Sprintf("%s:%s:%s", models.StaffContext, models.AddContext, models.CancelContext),
		),
	))
}
