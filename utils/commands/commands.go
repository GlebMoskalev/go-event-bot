package commands

import (
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CmdStart       = "start"
	CmdEvent       = "event"
	CmdAdminPanel  = "admin_panel"
	CmdChangeEvent = "change_event"
	CmdAddStaff    = "add_staff"
	CmdSearchStaff = "search_staff"
)

var CommandAccess = map[string]models.Role{
	CmdStart:       models.RoleGuest,
	CmdEvent:       models.RoleStaff,
	CmdAdminPanel:  models.RoleAdmin,
	CmdChangeEvent: models.RoleAdmin,
	CmdAddStaff:    models.RoleAdmin,
	CmdSearchStaff: models.RoleAdmin,
}

var BotCommands = map[models.Role][]tgbotapi.BotCommand{
	models.RoleGuest: {
		{
			Command:     CmdStart,
			Description: "Начала работы",
		},
	},
	models.RoleStaff: {
		{
			Command:     CmdEvent,
			Description: "Расписание мероприятий",
		},
	},
	models.RoleAdmin: {
		{
			Command:     CmdAddStaff,
			Description: "Добавить сотрудника",
		},
		{
			Command:     CmdSearchStaff,
			Description: "Поиск сотрудника",
		},
	},
}

func GetMenuCommands(role models.Role) []tgbotapi.BotCommand {
	var commands []tgbotapi.BotCommand
	switch role {
	case models.RoleGuest:
		commands = append(commands, BotCommands[role]...)
	case models.RoleStaff:
		commands = append(commands, BotCommands[role]...)
	case models.RoleAdmin:
		commands = append(commands, BotCommands[models.RoleStaff]...)
		commands = append(commands, BotCommands[role]...)
	}

	return commands
}
