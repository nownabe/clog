package clog_test

import (
	"context"
	"runtime"
	"strconv"
	"testing"

	"github.com/nownabe/clog"
)

func Test_sourceLocation(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)
	ctx := context.Background()

	pc, _, _, _ := runtime.Caller(0) // This must be called before calling l.Info
	l.Info(ctx, "foo")

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	want := buildWantLog("INFO", "foo")
	want["logging.googleapis.com/sourceLocation"] = map[string]any{
		"file":     frame.File,
		"line":     strconv.Itoa(frame.Line + 1),
		"function": frame.Function,
	}

	w.assertLog(t, want)
}
