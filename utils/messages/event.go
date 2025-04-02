package messages

import (
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
)

func EventEmpty() string {
	return "Расписание отсутствует"
}

func EventTitle() string {
	return "Мероприятия:"
}

func AllEvents(events []models.Event) string {
	res := "Мероприятия:\n\n"
	for _, schedule := range events {
		res += fmt.Sprintf(
			"🌐%s\n🎙️%s\n📍%s\n🗓️%s\n\n",
			schedule.Title,
			schedule.Speaker,
			schedule.Auditorium,
			schedule.Date.Format("02.01.2006 • 15:04"),
		)
	}
	return res
}
