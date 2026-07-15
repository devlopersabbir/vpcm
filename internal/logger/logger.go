package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// Level names
const (
	LevelTrace = slog.Level(-8)
)

// Init initializes the global logger based on environment/config settings.
func Init(levelStr string, format string, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	var level slog.Level
	switch levelStr {
	case "trace":
		level = LevelTrace
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = NewPrettyHandler(out, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// PrettyHandler is a custom slog handler for pretty human-readable outputs.
type PrettyHandler struct {
	slog.Handler
	out io.Writer
}

func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewTextHandler(out, opts),
		out:     out,
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelAttr := ""
	switch r.Level {
	case LevelTrace:
		levelAttr = "[TRACE]"
	case slog.LevelDebug:
		levelAttr = "\033[36m[DEBUG]\033[0m"
	case slog.LevelInfo:
		levelAttr = "\033[32m[INFO]\033[0m"
	case slog.LevelWarn:
		levelAttr = "\033[33m[WARN]\033[0m"
	case slog.LevelError:
		levelAttr = "\033[31m[ERROR]\033[0m"
	}

	timeStr := r.Time.Format(time.RFC3339)
	msg := r.Message

	_, err := io.WriteString(h.out, timeStr+" "+levelAttr+" "+msg+"\n")
	return err
}
