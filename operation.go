package clog

import (
	"context"
	"log/slog"

	"github.com/nownabe/clog/internal/keys"
)

type ctxKeyOperation struct{}

type operation struct {
	id       string
	producer string
}

// StartOperation returns a new context and a function to end the opration, starting the operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
func StartOperation(ctx context.Context, s Severity, msg, id, producer string) (context.Context, func(msg string)) {
	return Default().StartOperation(ctx, s, msg, id, producer)
}

// StartOperation returns a new context and a function to end the opration, starting the operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
func (l *Logger) StartOperation(ctx context.Context, s Severity, msg, id, producer string) (context.Context, func(msg string)) {
	l.logAttrs(ctx, s, msg, slog.Group(keys.Operation, "id", id, "producer", producer, "first", true))

	opCtx := context.WithValue(ctx, ctxKeyOperation{}, &operation{id, producer})

	return opCtx, func(msg string) {
		l.logAttrs(ctx, s, msg, slog.Group(keys.Operation, "id", id, "producer", producer, "last", true))
	}
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
