package logging

import (
	"context"
	"log/slog"
	"os"
)

type processHandler struct {
	slog.Handler
}

func NewProcessHandler(h slog.Handler) slog.Handler {
	return &processHandler{h}
}

var (
	_ slog.Handler = (*processHandler)(nil)
)

func (h *processHandler) Handle(ctx context.Context, r slog.Record) error {
	pid := os.Getpid()
	ppid := os.Getppid()
	attrs := []slog.Attr{slog.Int("pid", pid)}
	if ppid != 0 {
		attrs = append(attrs, slog.Int("ppid", ppid))
	}
	r.AddAttrs(attrs...)
	return h.Handler.Handle(ctx, r)
}

func (h *processHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewProcessHandler(h.Handler.WithAttrs(attrs))
}

func (h *processHandler) WithGroup(name string) slog.Handler {
	return NewProcessHandler(h.Handler.WithGroup(name))
}
