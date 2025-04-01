package handlers

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/configs"
)

type Bot interface {
	Start(ctx context.Context, cfg config.App, debugMode bool) error
}
