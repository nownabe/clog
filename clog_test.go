package clog_test

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
)

func setDefault(s clog.Severity, opts ...clog.Option) *writer {
	l, w := newLogger(s, opts...)
	clog.SetDefault(l)
	return w
}

func TestDefaltLogger_SimpleLogFuncs(t *testing.T) {
	type logFn func(ctx context.Context, msg string, args ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Debug, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Info, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Notice, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Warning, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Error, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Critical, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Alert, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Emergency, w
			},
			funcSeverity:   clog.SeverityEmergency,
			severityString: "EMERGENCY",
		},
	}

	type testCase struct {
		loggerSeverity clog.Severity
		msg            string
		args           []any
		want           map[string]any
	}

	makeCases := func(funcSeverity clog.Severity, severityString string) map[string]testCase {
		makeWant := func(loggerSeverity clog.Severity) map[string]any {
			if loggerSeverity > funcSeverity {
				return nil
			}
			return buildWantLog(severityString, "foo")
		}

		cases := map[string]testCase{}

		for _, s := range severities {
			cases[fmt.Sprintf("Severity(%d)", s)] = testCase{
				loggerSeverity: s,
				msg:            "foo",
				args:           []any{},
				want:           makeWant(s),
			}
		}

		cases["with args"] = testCase{
			loggerSeverity: funcSeverity,
			msg:            "foo",
			args:           []any{"k1", "v1", "k2", true},
			want:           buildWantLog(severityString, "foo", "k1", "v1", "k2", true),
		}

		return cases
	}

	ctx := context.Background()

	for testName, test := range tests {
		getFn := test.getFn
		for caseName, tc := range makeCases(test.funcSeverity, test.severityString) {
			tc := tc
			t.Run(testName+"/"+caseName, func(t *testing.T) {
				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.msg, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestDefaultLogger_FormattingLogFuncs(t *testing.T) {
	type logFn func(ctx context.Context, format string, a ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Debugf, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Infof, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Noticef, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Warningf, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Errorf, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Criticalf, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Alertf, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Emergencyf, w
			},
			funcSeverity:   clog.SeverityEmergency,
			severityString: "EMERGENCY",
		},
	}

	type testCase struct {
		loggerSeverity clog.Severity
		msg            string
		args           []any
		want           map[string]any
	}

	makeCases := func(funcSeverity clog.Severity, severityString string) map[string]testCase {
		makeWant := func(loggerSeverity clog.Severity) map[string]any {
			if loggerSeverity > funcSeverity {
				return nil
			}
			return buildWantLog(severityString, "foo")
		}

		cases := map[string]testCase{}

		for _, s := range severities {
			cases[fmt.Sprintf("Severity(%d)", s)] = testCase{
				loggerSeverity: s,
				msg:            "foo",
				args:           []any{},
				want:           makeWant(s),
			}
		}

		cases["with args"] = testCase{
			loggerSeverity: funcSeverity,
			msg:            "foo %q %02d",
			args:           []any{"1", 2},
			want:           buildWantLog(severityString, "foo \"1\" 02"),
		}

		return cases
	}

	ctx := context.Background()

	for testName, test := range tests {
		getFn := test.getFn
		for caseName, tc := range makeCases(test.funcSeverity, test.severityString) {
			tc := tc
			t.Run(testName+"/"+caseName, func(t *testing.T) {
				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.msg, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestDefaultLogger_ErrorLogFuncs(t *testing.T) {
	type logFn func(ctx context.Context, err error, args ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.DebugErr, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.InfoErr, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.NoticeErr, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.WarningErr, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.ErrorErr, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Err": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.Err, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.CriticalErr, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.AlertErr, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				w := setDefault(ls)
				return clog.EmergencyErr, w
			},
			funcSeverity:   clog.SeverityEmergency,
			severityString: "EMERGENCY",
		},
	}

	type testCase struct {
		loggerSeverity clog.Severity
		err            error
		args           []any
		want           map[string]any
	}

	makeCases := func(funcSeverity clog.Severity, severityString string) map[string]testCase {
		err1 := errors.New("err")
		err2 := errors.NewWithoutStack("no stack")

		makeWant := func(loggerSeverity clog.Severity) map[string]any {
			if loggerSeverity > funcSeverity {
				return nil
			}
			return buildWantLog(severityString, "err", "stack_trace", anyString{})
		}

		cases := map[string]testCase{}

		for _, s := range severities {
			cases[fmt.Sprintf("Severity(%d)", s)] = testCase{
				loggerSeverity: s,
				err:            err1,
				args:           []any{},
				want:           makeWant(s),
			}
		}

		cases["with args"] = testCase{
			loggerSeverity: funcSeverity,
			err:            err1,
			args:           []any{"k1", "v1"},
			want:           buildWantLog(severityString, "err", "k1", "v1", "stack_trace", anyString{}),
		}
		cases["error without stack"] = testCase{
			loggerSeverity: funcSeverity,
			err:            err2,
			args:           []any{},
			want:           buildWantLog(severityString, "no stack"),
		}
		cases["without error"] = testCase{
			loggerSeverity: funcSeverity,
			err:            nil,
			args:           []any{},
			want:           nil,
		}

		return cases
	}

	ctx := context.Background()

	for testName, test := range tests {
		getFn := test.getFn
		for caseName, tc := range makeCases(test.funcSeverity, test.severityString) {
			tc := tc
			t.Run(testName+"/"+caseName, func(t *testing.T) {
				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.err, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestDefaultLogger_Enabled(t *testing.T) {
	tests := map[string]struct {
		loggerSeverity clog.Severity
		arg            clog.Severity
		want           bool
	}{
		"Default-Default": {
			loggerSeverity: clog.SeverityDefault,
			arg:            clog.SeverityDefault,
			want:           true,
		},
		"Default-Debug": {
			loggerSeverity: clog.SeverityDefault,
			arg:            clog.SeverityDebug,
			want:           true,
		},
		"Debug-Default": {
			loggerSeverity: clog.SeverityDebug,
			arg:            clog.SeverityDefault,
			want:           false,
		},
		"Debug-Debug": {
			loggerSeverity: clog.SeverityDebug,
			arg:            clog.SeverityDebug,
			want:           true,
		},
		"Debug-Info": {
			loggerSeverity: clog.SeverityDebug,
			arg:            clog.SeverityInfo,
			want:           true,
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			setDefault(tt.loggerSeverity)
			got := clog.Enabled(ctx, tt.arg)
			if got != tt.want {
				t.Errorf("Enabled() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultLogger_Log(t *testing.T) {
	w := setDefault(clog.SeverityInfo)
	clog.Log(context.Background(), clog.SeverityInfo, "msg", "k1", "v1")
	w.assertLog(t, buildWantLog("INFO", "msg", "k1", "v1"))
}

func TestDefaultLogger_With(t *testing.T) {
	w := setDefault(clog.SeverityInfo)
	clog.With("k1", "v1").Info(context.Background(), "msg", "k2", "v2")
	w.assertLog(t, buildWantLog("INFO", "msg", "k1", "v1", "k2", "v2"))
}

func TestDefaultLogger_WithInsertID(t *testing.T) {
	w := setDefault(clog.SeverityInfo)
	clog.WithInsertID("id").Info(context.Background(), "msg")
	w.assertLog(t, buildWantLog("INFO", "msg", "logging.googleapis.com/insertId", "id"))
}

func TestDefaultLogger_HTTPRequest(t *testing.T) {
	sev := "INFO"
	msg := "msg"

	tests := map[string]struct {
		r    *clog.HTTPRequest
		want map[string]any
	}{
		"empty": {
			r:    &clog.HTTPRequest{},
			want: buildWantLog(sev, msg),
		},
		"full": {
			r: &clog.HTTPRequest{
				RequestMethod:                  "GET",
				RequestURL:                     "https://example.com/foo",
				RequestSize:                    123,
				Status:                         200,
				ResponseSize:                   456,
				UserAgent:                      "clog",
				RemoteIP:                       "203.0.113.1",
				ServerIP:                       "203.0.113.2",
				Referer:                        "https://example.com/referer",
				Latency:                        123*time.Second + 456*time.Nanosecond,
				CacheLookup:                    true,
				CacheHit:                       true,
				CacheValidatedWithOriginServer: true,
				CacheFillBytes:                 789,
				Protocol:                       "HTTP/1.1",
			},
			want: buildWantLog(sev, msg, "httpRequest", map[string]any{
				"requestMethod":                  "GET",
				"requestUrl":                     "https://example.com/foo",
				"requestSize":                    "123",
				"status":                         200,
				"responseSize":                   "456",
				"userAgent":                      "clog",
				"remoteIp":                       "203.0.113.1",
				"serverIp":                       "203.0.113.2",
				"referer":                        "https://example.com/referer",
				"latency":                        "123.000000456s",
				"cacheLookup":                    true,
				"cacheHit":                       true,
				"cacheValidatedWithOriginServer": true,
				"cacheFillBytes":                 "789",
				"protocol":                       "HTTP/1.1",
			}),
		},
		"only requestMethod": {
			r:    &clog.HTTPRequest{RequestMethod: "GET"},
			want: buildWantLog(sev, msg, "httpRequest", map[string]any{"requestMethod": "GET"}),
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			w := setDefault(clog.SeverityInfo)
			clog.WithHTTPRequest(tt.r).Info(ctx, msg)
			w.assertLog(t, tt.want)
		})
	}
}

func TestDefaultLogger_ContextWithLabels(t *testing.T) {
	w := setDefault(clog.SeverityInfo)

	ctx := context.Background()

	clog.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	func(ctx context.Context) {
		ctx, removeLabel := clog.ContextWithLabel(ctx, "lk1", "lv1")
		defer removeLabel()

		clog.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2", keyLabels, map[string]any{"lk1": "lv1"}))

		func(ctx context.Context) {
			ctx, removeLabel := clog.ContextWithLabel(ctx, "lk2", "lv2")
			defer removeLabel()

			clog.Info(ctx, "msg3")
			w.assertLog(t, buildWantLog("INFO", "msg3", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))
		}(ctx)

		clog.Info(ctx, "msg4")
		w.assertLog(t, buildWantLog("INFO", "msg4", keyLabels, map[string]any{"lk1": "lv1"}))
	}(ctx)

	clog.Info(ctx, "msg5")
	w.assertLog(t, buildWantLog("INFO", "msg5"))
}

func TestDefaultLogger_DefaultLables(t *testing.T) {
	defaultLabels := map[string]string{
		"lk1": "lv1",
		"lk2": "lv2",
	}

	w := setDefault(clog.SeverityInfo, clog.WithLabels(defaultLabels))

	ctx := context.Background()

	clog.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))

	func(ctx context.Context) {
		ctx, removeLabel := clog.ContextWithLabel(ctx, "lk3", "lv3")
		defer removeLabel()

		clog.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2", "lk3": "lv3"}))

		func(ctx context.Context) {
			ctx, removeLabel := clog.ContextWithLabel(ctx, "lk2", "LV2")
			defer removeLabel()

			clog.Info(ctx, "msg3")
			w.assertLog(t, buildWantLog("INFO", "msg3", keyLabels, map[string]any{"lk1": "lv1", "lk2": "LV2", "lk3": "lv3"}))
		}(ctx)

		clog.Info(ctx, "msg4")
		w.assertLog(t, buildWantLog("INFO", "msg4", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2", "lk3": "lv3"}))
	}(ctx)

	clog.Info(ctx, "msg5")
	w.assertLog(t, buildWantLog("INFO", "msg5", keyLabels, map[string]any{"lk1": "lv1", "lk2": "lv2"}))
}

func TestDefaultLogger_Operation(t *testing.T) {
	w := setDefault(clog.SeverityInfo)

	ctx := context.Background()

	clog.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	func() {
		ctx, end := clog.StartOperation(ctx, clog.SeverityInfo, "start", "id", "producer")
		defer end("end")

		w.assertLog(t, buildWantLog("INFO", "start",
			keyOperation, map[string]any{"id": "id", "producer": "producer", "first": true}))

		clog.Info(ctx, "msg2")
		w.assertLog(t, buildWantLog("INFO", "msg2",
			keyOperation, map[string]any{"id": "id", "producer": "producer"}))
	}()

	w.assertLog(t, buildWantLog("INFO", "end",
		keyOperation, map[string]any{"id": "id", "producer": "producer", "last": true}))

	clog.Info(ctx, "msg3")
	w.assertLog(t, buildWantLog("INFO", "msg3"))
}

func TestDefaultLogger_SourceLocation(t *testing.T) {
	w := setDefault(clog.SeverityInfo)
	ctx := context.Background()

	pc, _, _, _ := runtime.Caller(0) // This must be called before calling clog.Info.
	clog.Info(ctx, "foo")

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	want := buildWantLog("INFO", "foo")
	want["logging.googleapis.com/sourceLocation"] = map[string]any{
		"file":     frame.File,
		"line":     strconv.Itoa(frame.Line + 1),
		"function": frame.Function,
	}

	w.assertLog(t, want)
}

func TestDefaultLogger_Trace(t *testing.T) {
	const projectID = "test-project"
	traceIDStr := "000102030405060708090a0b0c0d0e0f"
	spanIDStr := "0001020304050607"

	w := setDefault(clog.SeverityInfo, clog.WithTrace(projectID))

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
	clog.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1"))

	ctx = trace.ContextWithSpanContext(context.Background(), spanCtx)
	clog.Info(ctx, "msg2")
	w.assertLog(t, buildWantLog("INFO", "msg2",
		keySpanID, spanIDStr,
		keyTrace, fmt.Sprintf("projects/%s/traces/%s", projectID, traceIDStr),
		keyTraceSampled, true))
}

func TestDefaultLogger_WithHandleFunc(t *testing.T) {
	type ctxkey struct{}

	f := func(next clog.HandleFunc) clog.HandleFunc {
		return (func(ctx context.Context, r slog.Record) error {
			if userID, ok := ctx.Value(ctxkey{}).(string); ok {
				r.AddAttrs(slog.String("user_id", userID))
			}
			return next(ctx, r)
		})
	}

	w := setDefault(clog.SeverityInfo, clog.WithHandleFunc(f))

	ctx := context.WithValue(context.Background(), ctxkey{}, "user1")
	clog.Info(ctx, "msg1")
	w.assertLog(t, buildWantLog("INFO", "msg1", "user_id", "user1"))
}
