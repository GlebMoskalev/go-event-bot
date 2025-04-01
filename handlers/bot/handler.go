package bot

import (
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type handler struct {
	user         services.User
	staff        services.Staff
	command      services.Command
	adminCommand services.AdminCommand
	message      services.Message
	schedule     services.Schedule
	log          *slog.Logger
}
