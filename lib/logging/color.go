package logging

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// define color code
var (
	levelToColor = map[slog.Level]string{
		slog.LevelError: "\x1b[31;20m",
		slog.LevelWarn:  "\x1b[33;20m",
		slog.LevelInfo:  resetColor,
		slog.LevelDebug: "\x1b[35;20m",
	}
	resetColor           = "\x1b[0m"
	enabledTerminalColor bool
)

func init() {
	colorTerminals := []string{
		"xterm",
		"vt100",
		"rxvt",
		"screen",
	}
	if v, ok := os.LookupEnv("TERM"); ok {
		for _, t := range colorTerminals {
			if strings.Contains(v, t) {
				enabledTerminalColor = true
				break
			}
		}
	}
}

func IsEnabledTerminalColor() bool {
	return enabledTerminalColor
}

type colorHandler struct {
	slog.Handler
	buf *bytes.Buffer
	w   io.Writer
}

func NewColorHandler(opts ...Option) slog.Handler {
	buf := &bytes.Buffer{}
	o := defaultOptions(opts...)
	w := o.writer
	o.writer = buf
	return newColorHandler(newJSONHandler(o), buf, w)
}

func newColorHandler(h slog.Handler, buf *bytes.Buffer, w io.Writer) slog.Handler {
	return &colorHandler{h, buf, w}
}

var (
	_ slog.Handler = (*processHandler)(nil)
)

func (h *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	err := h.Handler.Handle(ctx, r)
	if err != nil {
		return err
	}

	// 色付けロジック
	level := r.Level
	color := levelToColor[level]
	str := h.buf.String()
	if str != "" {
		if enabledTerminalColor {
			str = fmt.Sprintf("%s%s%s\n", color, str, resetColor)
		} else {
			str = fmt.Sprintf("%s\n", str)
		}
	}
	_, err = h.w.Write([]byte(str))
	h.buf.Reset()
	return err
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newColorHandler(h.Handler.WithAttrs(attrs), h.buf, h.w)
}

func (h *colorHandler) WithGroup(name string) slog.Handler {
	return newColorHandler(h.Handler.WithGroup(name), h.buf, h.w)
}
