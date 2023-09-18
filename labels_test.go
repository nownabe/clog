package clog_test

import (
	"context"
	"testing"

	"go.nownabe.dev/clog"
)

func Test_ContextWithLabels(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)

	ctx := context.Background()

	l.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	func(ctx context.Context) {
		ctx, removeLabel := clog.ContextWithLabel(ctx, "lk1", "lv1")
		defer removeLabel()

		l.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2", keyLabels, map[string]any{"lk1": "lv1"}))

		func(ctx context.Context) {
			ctx, removeLabel := clog.ContextWithLabel(ctx, "lk2", "lv2")
			defer removeLabel()

			l.Info(ctx, "msg3")
			w.assertLog(t, buildWantLog("INFO", "msg3", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))
		}(ctx)

		l.Info(ctx, "msg4")
		w.assertLog(t, buildWantLog("INFO", "msg4", keyLabels, map[string]any{"lk1": "lv1"}))
	}(ctx)

	l.Info(ctx, "msg5")
	w.assertLog(t, buildWantLog("INFO", "msg5"))
}

func Test_DefaultLables(t *testing.T) {
	t.Parallel()

	defaultLabels := map[string]string{
		"lk1": "lv1",
		"lk2": "lv2",
	}

	l, w := newLogger(clog.SeverityInfo, clog.WithLabels(defaultLabels))

	ctx := context.Background()

	l.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))

	func(ctx context.Context) {
		ctx, removeLabel := clog.ContextWithLabel(ctx, "lk3", "lv3")
		defer removeLabel()

		l.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2", "lk3": "lv3"}))

		func(ctx context.Context) {
			ctx, removeLabel := clog.ContextWithLabel(ctx, "lk2", "LV2")
			defer removeLabel()

			l.Info(ctx, "msg3")
			w.assertLog(t, buildWantLog("INFO", "msg3", keyLabels, map[string]any{"lk1": "lv1", "lk2": "LV2", "lk3": "lv3"}))
		}(ctx)

		l.Info(ctx, "msg4")
		w.assertLog(t, buildWantLog("INFO", "msg4", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2", "lk3": "lv3"}))
	}(ctx)

	l.Info(ctx, "msg5")
	w.assertLog(t, buildWantLog("INFO", "msg5", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))
}
