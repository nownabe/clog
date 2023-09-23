package clog

import (
	"context"
	"log/slog"
)

type HandleFunc func(context.Context, slog.Record) error

func WithHandleFunc(f func(next HandleFunc) HandleFunc) Option {
	return optionFunc(func(h slog.Handler) slog.Handler {
		return &customHandler{h, f, f(h.Handle)}
	})
}

type customHandler struct {
	slog.Handler

	f func(next HandleFunc) HandleFunc
	h HandleFunc
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.h(ctx, r)
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := h.Handler.WithAttrs(attrs)
	return &customHandler{handler, h.f, h.f(handler.Handle)}
}

func (h *customHandler) WithGroup(group string) slog.Handler {
	handler := h.Handler.WithGroup(group)
	return &customHandler{handler, h.f, h.f(handler.Handle)}
}
