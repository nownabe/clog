package clog

import (
	"context"
	"log/slog"
	"sync"

	"go.nownabe.dev/clog/internal/keys"
)

type (
	ctxKeyLabels        struct{}
	ctxKeyDefaultLabels struct{}
)

// ContextWithLabel returns a new context with the label that consists of given key and value.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func ContextWithLabel(ctx context.Context, key string, value string) (context.Context, func()) {
	labels, ok := ctx.Value(ctxKeyLabels{}).(*sync.Map)
	if !ok {
		labels = &sync.Map{}
	}
	labels.Store(key, value)

	return context.WithValue(ctx, ctxKeyLabels{}, labels), func() { labels.Delete(key) }
}

type labelsHandler struct {
	slog.Handler
}

func newLabelsHandler(h slog.Handler) slog.Handler {
	return &labelsHandler{h}
}

func (h *labelsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *labelsHandler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []slog.Attr
	labelsKeys := map[string]struct{}{}

	if labels, ok := ctx.Value(ctxKeyLabels{}).(*sync.Map); ok {
		labels.Range(func(key, val any) bool {
			keyStr, keyOK := key.(string)
			valStr, valOK := val.(string)
			if keyOK && valOK {
				labelsKeys[keyStr] = struct{}{}
				attrs = append(attrs, slog.String(keyStr, valStr))
			}
			return true
		})
	}

	if defaultLabels, ok := ctx.Value(ctxKeyDefaultLabels{}).(map[string]string); ok {
		for key, val := range defaultLabels {
			if _, ok := labelsKeys[key]; !ok {
				attrs = append(attrs, slog.String(key, val))
			}
		}
	}

	r.AddAttrs(slog.Any(keys.Labels, slog.GroupValue(attrs...)))

	return h.Handler.Handle(ctx, r)
}

func (h *labelsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &labelsHandler{h.Handler.WithAttrs(attrs)}
}

func (h *labelsHandler) WithGroup(group string) slog.Handler {
	return &labelsHandler{h.Handler.WithGroup(group)}
}

// WithLabels returns an Option that sets the default labels.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func WithLabels(labels map[string]string) Option {
	return optionFunc(func(h slog.Handler) slog.Handler {
		return newDefaultLabelsHandler(h, labels)
	})
}

type defaultLabelsHandler struct {
	slog.Handler

	labels map[string]string
}

func newDefaultLabelsHandler(h slog.Handler, labels map[string]string) slog.Handler {
	return &defaultLabelsHandler{h, labels}
}

func (h *defaultLabelsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *defaultLabelsHandler) Handle(ctx context.Context, r slog.Record) error {
	ctx = context.WithValue(ctx, ctxKeyDefaultLabels{}, h.labels)
	return h.Handler.Handle(ctx, r)
}

func (h *defaultLabelsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &defaultLabelsHandler{h.Handler.WithAttrs(attrs), h.labels}
}

func (h *defaultLabelsHandler) WithGroup(group string) slog.Handler {
	return &defaultLabelsHandler{h.Handler.WithGroup(group), h.labels}
}
