package admincommand

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type adminCmd struct {
	db              repositories.DB
	staffService    services.Staff
	userService     services.User
	scheduleService services.Event
	log             *slog.Logger
}

func New(
	db repositories.DB,
	staffSvc services.Staff,
	userSvc services.User,
	scheduleSvc services.Event,
	log *slog.Logger) services.AdminCommand {
	return &adminCmd{
		db:              db,
		staffService:    staffSvc,
		userService:     userSvc,
		scheduleService: scheduleSvc,
		log:             log,
	}
}
