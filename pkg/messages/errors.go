package messages

func Error() string {
	return "Произошла ошибка!"
}

func InvalidPhoneNumber() string {
	return "Ваш номер телефона не соответствует нашему формату"
}

func StaffNotFound() string {
	return "Вас нет в списках, попросите администратора добавить вас"
}

func ContactExists() string {
	return "Твой контакт у нас уже есть"
}

func UnknownCommand() string {
	return "Неизвестная команда"
}

func AccessDenied() string {
	return "У вас нет доступа к этой команде\n\n"
}
