package clog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"testing"

	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
)

const (
	keyLabels         = "logging.googleapis.com/labels"
	keyOperation      = "logging.googleapis.com/operation"
	keySourceLocation = "logging.googleapis.com/sourceLocation"
	keySpanID         = "logging.googleapis.com/spanId"
	keyTrace          = "logging.googleapis.com/trace"
	keyTraceSampled   = "logging.googleapis.com/trace_sampled"
)

type (
	anyVal    struct{}
	anyString struct{}
	anyNonNil struct{}
)

type writer struct {
	*bytes.Buffer
}

func (w *writer) assertLog(t *testing.T, want map[string]any) {
	t.Helper()

	written, err := w.ReadBytes('\n')
	if err != nil {
		if want == nil && err == io.EOF {
			return
		}
		t.Fatalf("w.ReadBytes('\\n') got error %v", err)
	}

	got := map[string]any{}
	if err := json.Unmarshal(written, &got); err != nil {
		t.Fatalf("json.Unmarshal(l, &got) got error %v", err)
	}

	assertEqual(t, want, got)
}

func assertEqual(t *testing.T, want, got map[string]any) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) got %d: %+v, want %d: %+v", len(got), got, len(want), want)
	}

	for k, wantVal := range want {
		gotRawVal, ok := got[k]
		if !ok {
			t.Errorf("got[%q] not found: %+v", k, got)
		}

		switch wantVal := wantVal.(type) {
		case anyVal:
			continue
		case anyString:
			if _, ok := gotRawVal.(string); !ok {
				t.Errorf("got[%q] got %#v (%T), want string value: %#v", k, gotRawVal, gotRawVal, got)
			}
		case anyNonNil:
			if gotRawVal == nil {
				t.Errorf("got[%q] got nil, want non-nil value: %#v", k, got)
			}
		case map[string]any:
			if gotMap, ok := gotRawVal.(map[string]any); ok {
				assertEqual(t, wantVal, gotMap)
			} else {
				t.Errorf("got[%q] got %#v (%T), want map[string]any value: %#v", k, gotRawVal, gotRawVal, got)
			}
		case *regexp.Regexp:
			if gotVal, ok := gotRawVal.(string); ok {
				if !wantVal.MatchString(gotVal) {
					t.Errorf("got[%q] got %q, want match regexp %v: %#v", k, gotVal, wantVal, got)
				}
			} else {
				t.Errorf("got[%q] got %#v (%T), want string and should match %v: %#v", k, gotRawVal, gotVal, wantVal, got)
			}
		case string:
			if gotVal, ok := gotRawVal.(string); ok {
				if wantVal != gotVal {
					t.Errorf("got[%q] got %q, want %q: %#v", k, gotVal, wantVal, got)
				}
			} else {
				t.Errorf("got[%q] got %#v (%T), want string value %q: %#v", k, gotRawVal, gotRawVal, wantVal, got)
			}
		case bool:
			if gotVal, ok := gotRawVal.(bool); ok {
				if wantVal != gotVal {
					t.Errorf("got[%q] got %t, want %t: %#v", k, gotVal, wantVal, got)
				}
			} else {
				t.Errorf("got[%q] got %#v (%T), want bool value %t: %#v", k, gotRawVal, gotRawVal, wantVal, got)
			}
		case int:
			// json.Unmarshal converts numbers to float64.
			if gotVal, ok := gotRawVal.(float64); ok {
				if float64(wantVal) != gotVal {
					t.Errorf("got[%q] got %f, want %d: %#v", k, gotVal, wantVal, got)
				}
			} else {
				t.Errorf("got[%q] got %#v (%T), want int value %d: %#v", k, gotRawVal, gotRawVal, wantVal, got)
			}
		default:
			panic(fmt.Sprintf("unexpected want value %#v (%T)", wantVal, wantVal))
		}
	}
}

func newLogger(s clog.Severity, opts ...clog.Option) (*clog.Logger, *writer) {
	w := &writer{&bytes.Buffer{}}
	return clog.New(w, s, true, opts...), w
}

var (
	timeRE         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)
	sourceLocation = map[string]any{
		"file":     anyString{},
		"line":     anyString{},
		"function": anyString{},
	}
)

func buildWantLog(severity, msg string, args ...any) map[string]any {
	want := map[string]any{
		"severity": severity,
		"message":  msg,
		"time":     timeRE,
	}
	want[keySourceLocation] = sourceLocation

	for i := 0; i < len(args); i += 2 {
		want[args[i].(string)] = args[i+1]
	}

	return want
}

var severities = []clog.Severity{
	clog.SeverityDefault,
	clog.SeverityDebug,
	clog.SeverityInfo,
	clog.SeverityNotice,
	clog.SeverityWarning,
	clog.SeverityError,
	clog.SeverityCritical,
	clog.SeverityAlert,
	clog.SeverityEmergency,
}

func TestLogger_SimpleLogFuncs(t *testing.T) {
	t.Parallel()

	type logFn func(ctx context.Context, msg string, args ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Debug, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Info, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Notice, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Warning, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Error, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Critical, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Alert, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Emergency, w
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
				t.Parallel()

				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.msg, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestLogger_FormattingLogFuncs(t *testing.T) {
	t.Parallel()

	type logFn func(ctx context.Context, format string, a ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Debugf, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Infof, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Noticef, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Warningf, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Errorf, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Criticalf, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Alertf, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Emergencyf, w
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
				t.Parallel()

				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.msg, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestLogger_ErrorLogFuncs(t *testing.T) {
	t.Parallel()

	type logFn func(ctx context.Context, err error, args ...any)

	tests := map[string]struct {
		getFn          func(clog.Severity) (logFn, *writer)
		funcSeverity   clog.Severity
		severityString string
	}{
		"Debug": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.DebugErr, w
			},
			funcSeverity:   clog.SeverityDebug,
			severityString: "DEBUG",
		},
		"Info": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.InfoErr, w
			},
			funcSeverity:   clog.SeverityInfo,
			severityString: "INFO",
		},
		"Notice": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.NoticeErr, w
			},
			funcSeverity:   clog.SeverityNotice,
			severityString: "NOTICE",
		},
		"Warning": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.WarningErr, w
			},
			funcSeverity:   clog.SeverityWarning,
			severityString: "WARNING",
		},
		"Error": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.ErrorErr, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Err": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.Err, w
			},
			funcSeverity:   clog.SeverityError,
			severityString: "ERROR",
		},
		"Critical": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.CriticalErr, w
			},
			funcSeverity:   clog.SeverityCritical,
			severityString: "CRITICAL",
		},
		"Alert": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.AlertErr, w
			},
			funcSeverity:   clog.SeverityAlert,
			severityString: "ALERT",
		},
		"Emergency": {
			getFn: func(ls clog.Severity) (logFn, *writer) {
				l, w := newLogger(ls)
				return l.EmergencyErr, w
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
				t.Parallel()

				fn, w := getFn(tc.loggerSeverity)
				fn(ctx, tc.err, tc.args...)
				w.assertLog(t, tc.want)
			})
		}
	}
}

func TestLogger_Enabled(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			l, _ := newLogger(tt.loggerSeverity)
			got := l.Enabled(ctx, tt.arg)
			if got != tt.want {
				t.Errorf("Enabled() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_Log(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)
	l.Log(context.Background(), clog.SeverityInfo, "msg", "k1", "v1")
	w.assertLog(t, buildWantLog("INFO", "msg", "k1", "v1"))
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)
	l.With("k1", "v1").Info(context.Background(), "msg", "k2", "v2")
	w.assertLog(t, buildWantLog("INFO", "msg", "k1", "v1", "k2", "v2"))
}

func TestLogger_WithInsertID(t *testing.T) {
	t.Parallel()

	l, w := newLogger(clog.SeverityInfo)
	l.WithInsertID("id").Info(context.Background(), "msg")
	w.assertLog(t, buildWantLog("INFO", "msg", "logging.googleapis.com/insertId", "id"))
}
