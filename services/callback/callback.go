package callback

import (
	"github.com/GlebMoskalev/go-event-bot/repositories"
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type callback struct {
	db           repositories.DB
	userService  services.User
	eventService services.Event
	stateService services.State
	staffService services.Staff
	log          *slog.Logger
}

func New(
	db repositories.DB,
	userSvc services.User,
	scheduleSvc services.Event,
	stateService services.State,
	staffService services.Staff,
	log *slog.Logger) services.Callback {
	return &callback{
		db:           db,
		userService:  userSvc,
		eventService: scheduleSvc,
		stateService: stateService,
		staffService: staffService,
		log:          log,
	}
}
