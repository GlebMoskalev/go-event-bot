package bot

import (
	"github.com/GlebMoskalev/go-event-bot/services"
	"log/slog"
)

type handler struct {
	user         services.User
	staff        services.Staff
	command      services.Command
	callback     services.Callback
	adminCommand services.AdminCommand
	message      services.Message
	schedule     services.Event
	state        services.State
	log          *slog.Logger
}
