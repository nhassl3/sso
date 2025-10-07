package slogpretty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value
		return true
	})

	// TODO: write doc about any method of Value struct
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// TODO: implement method
	_ = name
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h == nil || h.opts.SlogOpts == nil || h.opts.SlogOpts.Level == nil {
		return level >= slog.LevelInfo
	}

	// Don't panic, but throw an error in the level logger options to stderr, if there is one
	defer func() {
		if r := recover(); r != nil {
			_, err := fmt.Fprintf(os.Stderr, "Recovered from panic in Enabled %v\n", r)
			if err != nil {
				return
			}
		}
	}()

	return level >= h.opts.SlogOpts.Level.Level()
}
