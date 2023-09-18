package clog_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/nownabe/clog"
)

func Test_WithHandleFunc(t *testing.T) {
	t.Parallel()

	type ctxkey struct{}

	f := func(next clog.HandleFunc) clog.HandleFunc {
		return (func(ctx context.Context, r slog.Record) error {
			if userID, ok := ctx.Value(ctxkey{}).(string); ok {
				r.AddAttrs(slog.String("user_id", userID))
			}
			return next(ctx, r)
		})
	}

	l, w := newLogger(clog.SeverityInfo, clog.WithHandleFunc(f))

	ctx := context.WithValue(context.Background(), ctxkey{}, "user1")
	l.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1", "user_id", "user1"))
}
