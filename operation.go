package clog

import (
	"context"
	"log/slog"

	"go.nownabe.dev/clog/internal/keys"
)

type ctxKeyOperation struct{}

type operation struct {
	id       string
	producer string
}

type operationHandler struct {
	slog.Handler
}

func newOperationHandler(h slog.Handler) slog.Handler {
	return &operationHandler{h}
}

func (h *operationHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *operationHandler) Handle(ctx context.Context, r slog.Record) error {
	if op, ok := ctx.Value(ctxKeyOperation{}).(*operation); ok {
		r.AddAttrs(slog.Group(keys.Operation, "id", op.id, "producer", op.producer))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *operationHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &operationHandler{h.Handler.WithAttrs(attrs)}
}

func (h *operationHandler) WithGroup(group string) slog.Handler {
	return &operationHandler{h.Handler.WithGroup(group)}
}
