package message

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type msg struct {
	db             repositories.DB
	commandService services.Command
	staffService   services.Staff
	userService    services.User
	stateService   services.State
	log            *slog.Logger
}

func New(
	db repositories.DB,
	staffSvc services.Staff,
	userSvc services.User,
	commandService services.Command,
	stateService services.State,
	log *slog.Logger) services.Message {
	return &msg{
		db:             db,
		staffService:   staffSvc,
		userService:    userSvc,
		commandService: commandService,
		stateService:   stateService,
		log:            log,
	}
}
