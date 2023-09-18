package clog_test

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/nownabe/clog"
)

func Test_Trace(t *testing.T) {
	t.Parallel()

	const projectID = "test-project"
	traceIDStr := "000102030405060708090a0b0c0d0e0f"
	spanIDStr := "0001020304050607"

	l, w := newLogger(clog.SeverityInfo, clog.WithTrace(projectID))

	traceID, err := trace.TraceIDFromHex(traceIDStr)
	if err != nil {
		panic(err)
	}
	spanID, err := trace.SpanIDFromHex(spanIDStr)
	if err != nil {
		panic(err)
	}
	cfg := trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
		Remote:     false,
	}
	spanCtx := trace.NewSpanContext(cfg)

	ctx := context.Background()
	l.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	ctx = trace.ContextWithSpanContext(context.Background(), spanCtx)
	l.Info(ctx, "msg2")
	w.assertLog(t, buildWantLog("INFO", "msg2",
		keySpanID, spanIDStr,
		keyTrace, fmt.Sprintf("projects/%s/traces/%s", projectID, traceIDStr),
		keyTraceSampled, true))
}
