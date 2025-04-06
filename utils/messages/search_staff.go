package messages

import (
	"fmt"
	"github.com/GlebMoskalev/go-event-bot/models"
)

func StaffSearchMethod() string {
	return "Выберите способ поиска сотрудника"
}

func StaffNodFoundSearch() string {
	return "Сотрудников нет"
}

func StaffList(staff []models.Staff) string {
	res := ""
	for _, s := range staff {
		res += fmt.Sprintf("%s %s %s\n\t%s\n\n", s.LastName, s.FirstName, s.Patronymic, s.PhoneNumber)
	}
	return res
}

func LastNameTooShort() string {
	return "Фамилия слишком короткая. Укажите хотя бы 2 символа."
}
