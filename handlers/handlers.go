package handlers

import (
	"context"
	"github.com/GlebMoskalev/go-event-bot/config"
)

type Bot interface {
	Start(ctx context.Context, cfg config.App, debugMode bool) error
}
