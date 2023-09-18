package clog

import (
	"context"
	"log/slog"
)

type HandleFunc func(context.Context, slog.Record) error

func WithHandleFunc(f func(next HandleFunc) HandleFunc) Option {
	return optionFunc(func(h slog.Handler) slog.Handler {
		return &customHandler{h, f(h.Handle)}
	})
}

type customHandler struct {
	slog.Handler

	f HandleFunc
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.f(ctx, r)
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customHandler{h.Handler.WithAttrs(attrs), h.f}
}

func (h *customHandler) WithGroup(group string) slog.Handler {
	return &customHandler{h.Handler.WithGroup(group), h.f}
}
