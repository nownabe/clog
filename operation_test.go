package clog_test

import (
	"context"
	"testing"

	"github.com/nownabe/clog"
)

func Test_Operation(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)

	ctx := context.Background()

	l.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	func() {
		ctx, end := l.StartOperation(ctx, clog.SeverityInfo, "start", "id", "producer")
		defer end("end")

		w.assertLog(t, buildWantLog("INFO", "start", keyOperation, map[string]any{"id": "id", "producer": "producer", "first": true}))

		l.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2", keyOperation, map[string]any{"id": "id", "producer": "producer"}))
	}()
	w.assertLog(t, buildWantLog("INFO", "end", keyOperation, map[string]any{"id": "id", "producer": "producer", "last": true}))

	l.Info(ctx, "msg3")
	w.assertLog(t, buildWantLog("INFO", "msg3"))
}
