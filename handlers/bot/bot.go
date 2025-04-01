package bot

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/configs"
	"github.com/GlebMoskalev/go-event-bot/handlers"
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services/admincommand"
	"github.com/GlebMoskalev/go-event-bot/services/command"
	"github.com/GlebMoskalev/go-event-bot/services/message"
	"github.com/GlebMoskalev/go-event-bot/services/schedule"
	"github.com/GlebMoskalev/go-event-bot/services/staff"
	"github.com/GlebMoskalev/go-event-bot/services/user"
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
	sched := schedule.New(db, log)
	cmd := command.New(db, stf, usr, sched, log)
	adminCmd := admincommand.New(db, stf, usr, sched, log)
	msg := message.New(db, stf, usr, cmd, log)
	handler := &handler{
		user:         usr,
		staff:        stf,
		command:      cmd,
		adminCommand: adminCmd,
		message:      msg,
		schedule:     sched,
		log:          log,
	}

	return &bot{handler: handler, log: log}
}

func (b *bot) Start(ctx context.Context, cfg config.App, debugMode bool) error {
	b.log.Info("bot starting...")

	tgbot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		b.log.Error("failed to create bot: %v", err)
		return err
	}

	tgbot.Debug = debugMode

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Bot.UpdateTimeout

	updates := tgbot.GetUpdatesChan(u)
	for update := range updates {
		msg := update.Message
		if msg == nil {
			continue
		} else if update.Message.IsCommand() {
			go b.handler.Commands(ctx, tgbot, update)
		} else {
			go b.handler.Message(ctx, tgbot, update)
		}
	}

	return nil
}
