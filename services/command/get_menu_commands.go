package command

import (
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var DefaultCommands = []tgbotapi.BotCommand{
	{
		Command:     "start",
		Description: "Начало работы",
	},
}

var staffCommands = []tgbotapi.BotCommand{
	{
		Command:     "schedule",
		Description: "Расписание мероприятий",
	},
}

var adminCommands = []tgbotapi.BotCommand{
	{
		Command:     "admin_panel",
		Description: "Панель администратора",
	},
	{
		Command:     "change_event",
		Description: "Изменить расписание",
	},
	{
		Command:     "add_staff",
		Description: "Добавить сотрудника",
	},
}

func (c *cmd) GetMenuCommands(role models.Role) []tgbotapi.BotCommand {
	var commands []tgbotapi.BotCommand

	switch role {
	case models.RoleGuest:
		commands = append(commands, DefaultCommands...)
	case models.RoleStaff:
		commands = append(commands, staffCommands...)
	case models.RoleAdmin:
		commands = append(commands, staffCommands...)
		commands = append(commands, adminCommands...)
	}

	return commands
}
