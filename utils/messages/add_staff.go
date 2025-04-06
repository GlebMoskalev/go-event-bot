package messages

import "fmt"

func InvalidFullNameFormat() string {
	return "Неправильный формат. Введите ФИО полностью (Фамилия Имя Отчество)"
}

func InvalidPhoneFormat() string {
	return "Неправильный формат номера телефона.\n Введите в формате 7XXXXXXXXXX"
}

func PhoneAlreadyExists() string {
	return "Сотрудник с таким номером телефона существует, введите другой номер"
}

func RequestPhoneNumber() string {
	return "Введите номер телефона в формате:\n 7XXXXXXXXXX"
}

func RequestFullName() string {
	return "Введите ФИО сотрудника в формате:\n Фамилия Имя Отчество"
}
func ConfirmStaffAddition(lastName, firstName, patronymic, phoneNumber string) string {
	return fmt.Sprintf(
		"Подтвердите, что вы хотите добавить сотрудника:\n%s %s %s\n%s\n\n",
		lastName, firstName, patronymic, phoneNumber,
	)
}
func StaffAdded() string {
	return "Сотрудник добавлен"
}

func StaffAdditionCancelled() string {
	return "Добавление сотрудника отменено"
}

func StaffAdditionMissing() string {
	return "Не найден начатый процесс добавления сотрудника"
}
