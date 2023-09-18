package clog

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"go.nownabe.dev/clog/internal/keys"
)

// WithTrace returns an Option that sets the trace attributes to the log record.
func WithTrace(projectID string) Option {
	return optionFunc(func(h slog.Handler) slog.Handler {
		return &traceHandler{h, projectID}
	})
}

type traceHandler struct {
	slog.Handler

	projectID string
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		r.AddAttrs(
			slog.String(keys.Trace, fmt.Sprintf("projects/%s/traces/%s", h.projectID, spanCtx.TraceID().String())),
			slog.String(keys.SpanID, spanCtx.SpanID().String()),
			slog.Bool(keys.TraceSampled, spanCtx.IsSampled()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{h.Handler.WithAttrs(attrs), h.projectID}
}

func (h *traceHandler) WithGroup(group string) slog.Handler {
	return &traceHandler{h.Handler.WithGroup(group), h.projectID}
}
