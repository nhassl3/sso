package slogdiscard

import (
	"context"
	"io"
	"log"
	"log/slog"
)

type Handler interface {
	Enabled(context.Context, slog.Level) bool
	Handle(ctx context.Context, r slog.Record) error
	WithAttrs(attrs []slog.Attr) Handler
	WithGroup(name string) Handler
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	return &PrettyHandler{}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}
