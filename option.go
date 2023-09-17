package clog

import "log/slog"

type Option interface {
	apply(h slog.Handler) slog.Handler
}

type optionFunc func(h slog.Handler) slog.Handler

func (f optionFunc) apply(h slog.Handler) slog.Handler {
	return f(h)
}
