package commands

import (
	"github.com/GlebMoskalev/go-event-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CmdStart       = "start"
	CmdSchedule    = "schedule"
	CmdAdminPanel  = "admin_panel"
	CmdChangeEvent = "change_event"
	CmdAddStaff    = "add_staff"
)

var CommandAccess = map[string]models.Role{
	CmdStart:       models.RoleGuest,
	CmdSchedule:    models.RoleStaff,
	CmdAdminPanel:  models.RoleAdmin,
	CmdChangeEvent: models.RoleAdmin,
	CmdAddStaff:    models.RoleAdmin,
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
			Command:     CmdSchedule,
			Description: "Расписание мероприятий",
		},
	},
	models.RoleAdmin: {
		{
			Command:     CmdAdminPanel,
			Description: "Панель администратора",
		},
		{
			Command:     CmdChangeEvent,
			Description: "Изменить расписание",
		},
		{
			Command:     CmdAddStaff,
			Description: "Добавить сотрудника",
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
