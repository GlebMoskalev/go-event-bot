package messages

import (
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
)

func ScheduleEmpty() string {
	return "Расписание отсутствует"
}

func ScheduleTitle() string {
	return "Расписание:"
}

func AllSchedules(schedules []models.Schedule) string {
	res := "Расписание:\n\n"
	for _, schedule := range schedules {
		res += fmt.Sprintf("🎤%s\n🗓️%s\n\n", schedule.Title, schedule.Date.Format("02.01.2006 • 15:04"))
	}
	return res
}
