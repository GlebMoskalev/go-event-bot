package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

const (
	red         = "\033[31m"
	gray        = "\033[90m"
	cyan        = "\033[36m"
	white       = "\033[97m"
	lightgray   = "\033[37m"
	lightYellow = "\033[93m"
)

func colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, "\033[0m")
}

type prettySlog struct {
	minLevel slog.Level
}

func newPrettySlog(minLevel slog.Level) *prettySlog {
	return &prettySlog{minLevel: minLevel}
}

func (s *prettySlog) Enabled(_ context.Context, level slog.Level) bool {
	return level >= s.minLevel
}

func (s *prettySlog) Handle(_ context.Context, r slog.Record) error {
	var levelColor, levelText string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = lightgray
		levelText = "DEBUG"
	case slog.LevelInfo:
		levelColor = cyan
		levelText = "INFO"
	case slog.LevelWarn:
		levelColor = lightYellow
		levelText = "WARN"
	case slog.LevelError:
		levelColor = red
		levelText = "ERROR"
	}

	timeStr := r.Time.Format("15:04:05.000")

	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	_, err := fmt.Fprintf(os.Stdout, "%s [%s] %s %s\n",
		colorize(gray, timeStr),
		colorize(levelColor, levelText),
		colorize(white, r.Message),
		colorize(lightgray, attrs),
	)
	return err
}

func (s *prettySlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return s
}

func (s *prettySlog) WithGroup(name string) slog.Handler {
	return s
}
