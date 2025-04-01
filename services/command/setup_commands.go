package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (c *cmd) SetupCommands(msg tgbotapi.MessageConfig, isAdmin bool) tgbotapi.MessageConfig {
	var commands []tgbotapi.BotCommand

	commands = append(commands,
		tgbotapi.BotCommand{
			Command:     "schedule",
			Description: "Узнать расписание",
		},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Помощь",
		},
	)

	if isAdmin {
		adminCommands := []tgbotapi.BotCommand{
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

		commands = append(commands, adminCommands...)
	}

	msg.ReplyMarkup = tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeChat(msg.ChatID),
		commands...,
	)
	return msg
}
