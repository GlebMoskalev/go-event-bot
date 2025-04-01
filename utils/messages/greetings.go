package messages

import "fmt"

func Welcome(firstName, patronymic string) string {
	return fmt.Sprintf("%s %s, приветствуем на нашем мероприятии!", firstName, patronymic)
}

func RequestContact() string {
	return "Привет! Для дальнейшей работы нужен твой контакт!"
}
