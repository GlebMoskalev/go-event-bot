package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/configs"
	"github.com/GlebMoskalev/go-event-bot/handlers"
	"github.com/GlebMoskalev/go-event-bot/models"
	"github.com/GlebMoskalev/go-event-bot/pkg/logger"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services/admincommand"
	"github.com/GlebMoskalev/go-event-bot/services/callback"
	"github.com/GlebMoskalev/go-event-bot/services/command"
	"github.com/GlebMoskalev/go-event-bot/services/event"
	"github.com/GlebMoskalev/go-event-bot/services/message"
	"github.com/GlebMoskalev/go-event-bot/services/staff"
	"github.com/GlebMoskalev/go-event-bot/services/state"
	"github.com/GlebMoskalev/go-event-bot/services/user"
	"github.com/GlebMoskalev/go-event-bot/utils/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type bot struct {
	handler *handler
	log     *slog.Logger
}

func New(db repositories.DB, log *slog.Logger) handlers.Bot {

	usr := user.New(db, log)
	stf := staff.New(db, log)
	sched := event.New(db, log)
	cmd := command.New(db, usr, sched, log)
	cbk := callback.New(db, usr, sched, log)
	st := state.New(db, log)
	adminCmd := admincommand.New(db, stf, usr, sched, st, log)
	msg := message.New(db, stf, usr, cmd, st, log)
	handler := &handler{
		user:         usr,
		staff:        stf,
		command:      cmd,
		callback:     cbk,
		adminCommand: adminCmd,
		message:      msg,
		schedule:     sched,
		state:        st,
		log:          log,
	}

	return &bot{handler: handler, log: log}
}

func (b *bot) Start(ctx context.Context, cfg config.App, debugMode bool) error {
	log := logger.SetupLogger(b.log, "bot", "Start")
	log.Info("bot starting", "debug_mode", debugMode)

	tgbot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		b.log.Error("failed to create bot: %v", err)
		return err
	}

	tgbot.Debug = debugMode
	log.Info("bot api initialized successfully")

	_, err = tgbot.Request(tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeDefault(),
		commands.GetMenuCommands(models.RoleGuest)...,
	))
	if err != nil {
		log.Error("failed to set default commands", "error", err)
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Bot.UpdateTimeout

	updates := tgbot.GetUpdatesChan(u)
	for update := range updates {
		var st models.State
		if update.Message != nil && update.Message.Chat != nil {
			st, _ = b.handler.state.Get(ctx, update.Message.Chat.ID)
		}
		if st != "" {
			go b.handler.Message(ctx, tgbot, update, st)
		} else if update.CallbackQuery != nil {
			go b.handler.Callbacks(ctx, tgbot, update)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				go b.handler.Commands(ctx, tgbot, update)
			} else {
				go b.handler.Message(ctx, tgbot, update, "")
			}
		}
	}

	return nil
}
