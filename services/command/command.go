package command

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type cmd struct {
	db              repositories.DB
	staffService    services.Staff
	userService     services.User
	scheduleService services.Schedule
	log             *slog.Logger
}

func New(
	db repositories.DB,
	staffSvc services.Staff,
	userSvc services.User,
	scheduleSvc services.Schedule,
	log *slog.Logger) services.Command {
	return &cmd{
		db:              db,
		staffService:    staffSvc,
		userService:     userSvc,
		scheduleService: scheduleSvc,
		log:             log,
	}
}
