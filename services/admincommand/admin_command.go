package admincommand

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type adminCmd struct {
	db           repositories.DB
	staffService services.Staff
	userService  services.User
	eventService services.Event
	state        services.State
	log          *slog.Logger
}

func New(
	db repositories.DB,
	staffSvc services.Staff,
	userSvc services.User,
	scheduleSvc services.Event,
	state services.State,
	log *slog.Logger) services.AdminCommand {
	return &adminCmd{
		db:           db,
		staffService: staffSvc,
		userService:  userSvc,
		eventService: scheduleSvc,
		state:        state,
		log:          log,
	}
}
