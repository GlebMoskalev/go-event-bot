package callback

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type callback struct {
	db              repositories.DB
	userService     services.User
	scheduleService services.Event
	log             *slog.Logger
}

func New(
	db repositories.DB,
	userSvc services.User,
	scheduleSvc services.Event,
	log *slog.Logger) services.Callback {
	return &callback{
		db:              db,
		userService:     userSvc,
		scheduleService: scheduleSvc,
		log:             log,
	}
}
