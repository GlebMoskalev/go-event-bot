package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type CommandManager struct {
	bot *tgbotapi.BotAPI
	log *slog.Logger
}

func NewCommandManger(bot *tgbotapi.BotAPI, log *slog.Logger) *CommandManager {
	return &CommandManager{bot: bot, log: log}
}

func (m *CommandManager) SetupCommands(update tgbotapi.Update, isAdmin bool) error {
	chatID := update.Message.Chat.ID

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
				Command:     "admin",
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

	cmgCfgWithScope := tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeChat(chatID),
		commands...,
	)

	_, err := m.bot.Send(cmgCfgWithScope)
	if err != nil {
		m.log.Error("failed to set commands", "error", err)
		return err
	}

	return nil
}
