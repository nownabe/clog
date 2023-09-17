package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/nownabe/clog/errors"
	"github.com/nownabe/clog/internal/keys"
)

type Logger struct {
	inner *slog.Logger
}

func New(w io.Writer, s Severity, json bool, opts ...Option) *Logger {
	opt := &slog.HandlerOptions{
		Level: s,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			a = replaceLevel(a)
			a = replaceMessage(a)
			return a
		},
	}

	var h slog.Handler
	if json {
		h = slog.NewJSONHandler(w, opt)
	} else {
		h = slog.NewTextHandler(w, opt)
	}

	h = newLabelsHandler(h)
	for _, o := range opts {
		h = o.apply(h)
	}

	return &Logger{slog.New(h)}
}

// Debug logs at SeverityDebug.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityDebug, msg, args...)
}

// Debugf logs formatted in the manner of fmt.Printf at SeverityDebug.
func (l *Logger) Debugf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityDebug, fmt.Sprintf(format, a...))
}

// DebugErr logs an error at SeverityDebug.
func (l *Logger) DebugErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityDebug, err, args...)
}

// Info logs at SeverityInfo.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityInfo, msg, args...)
}

// Infof logs formatted in the manner of fmt.Printf at SeverityInfo.
func (l *Logger) Infof(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityInfo, fmt.Sprintf(format, a...))
}

// InfoErr logs an error at SeverityInfo.
func (l *Logger) InfoErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityInfo, err, args...)
}

// Notice logs at SeverityNotice.
func (l *Logger) Notice(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityNotice, msg, args...)
}

// Noticef logs formatted in the manner of fmt.Printf at SeverityNotice.
func (l *Logger) Noticef(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityNotice, fmt.Sprintf(format, a...))
}

// NoticeErr logs an error at SeverityNotice.
func (l *Logger) NoticeErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityNotice, err, args...)
}

// Warning logs at SeverityWarning.
func (l *Logger) Warning(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityWarning, msg, args...)
}

// Warningf logs formatted in the manner of fmt.Printf at SeverityWarning.
func (l *Logger) Warningf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityWarning, fmt.Sprintf(format, a...))
}

// WarningErr logs an error at SeverityWarning.
func (l *Logger) WarningErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityWarning, err, args...)
}

// Error logs at SeverityError.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityError, msg, args...)
}

// Errorf logs formatted in the manner of fmt.Printf at SeverityError.
func (l *Logger) Errorf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityError, fmt.Sprintf(format, a...))
}

// ErrorErr logs an error at SeverityError.
func (l *Logger) ErrorErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityError, err, args...)
}

// Critical logs at SeverityCritical.
func (l *Logger) Critical(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityCritical, msg, args...)
}

// Criticalf logs formatted in the manner of fmt.Printf at SeverityCritical.
func (l *Logger) Criticalf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityCritical, fmt.Sprintf(format, a...))
}

// CriticalErr logs an error at SeverityCritical.
func (l *Logger) CriticalErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityCritical, err, args...)
}

// Alert logs at SeverityAlert.
func (l *Logger) Alert(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityAlert, msg, args...)
}

// Alertf logs formatted in the manner of fmt.Printf at SeverityAlert.
func (l *Logger) Alertf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityAlert, fmt.Sprintf(format, a...))
}

// AlertErr logs an error at SeverityAlert.
func (l *Logger) AlertErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityAlert, err, args...)
}

// Emergency logs at SeverityEmergency.
func (l *Logger) Emergency(ctx context.Context, msg string, args ...any) {
	l.log(ctx, SeverityEmergency, msg, args...)
}

// Emergencyf logs formatted in the manner of fmt.Printf at SeverityEmergency.
func (l *Logger) Emergencyf(ctx context.Context, format string, a ...any) {
	l.log(ctx, SeverityEmergency, fmt.Sprintf(format, a...))
}

// EmergencyErr logs an error at SeverityEmergency.
func (l *Logger) EmergencyErr(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityEmergency, err, args...)
}

// Enabled reports whether the Logger emits log records at the given context and level.
func (l *Logger) Enabled(ctx context.Context, s Severity) bool {
	return l.inner.Enabled(ctx, s)
}

// Err is a shorthand for ErrorErr.
func (l *Logger) Err(ctx context.Context, err error, args ...any) {
	l.err(ctx, SeverityError, err, args...)
}

// Log emits a log record with the current time and the given level and message.
func (l *Logger) Log(ctx context.Context, s Severity, msg string, args ...any) {
	l.log(ctx, s, msg, args...)
}

// With returns a Logger that includes the given attributes in each output operation.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.inner.With(args...)}
}

// WithHTTPRequest returns a Logger that includes the given httpRequest in each output operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func (l *Logger) WithHTTPRequest(req *HTTPRequest) *Logger {
	return l.withAttrs(slog.Any(keys.HTTPRequest, req))
}

// WithInsertID returns a Logger that includes the given insertId in each output operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func (l *Logger) WithInsertID(id string) *Logger {
	return l.withAttrs(slog.String(keys.InsertID, id))
}

func (l *Logger) log(ctx context.Context, s Severity, msg string, args ...any) {
	// skip [runtime.Callers, source, this function, clog exported function]
	src := getSourceLocation(4)
	l.logWithSource(ctx, s, src, msg, args...)
}

func (l *Logger) logAttrs(ctx context.Context, s Severity, msg string, attrs ...slog.Attr) {
	// skip [runtime.Callers, source, this function, clog exported function]
	src := getSourceLocation(4)
	l.logAttrsWithSource(ctx, s, src, msg, attrs...)
}

func (l *Logger) withAttrs(attrs ...slog.Attr) *Logger {
	return &Logger{slog.New(l.inner.Handler().WithAttrs(attrs))}
}

func (l *Logger) err(ctx context.Context, s Severity, err error, args ...any) {
	if err == nil {
		return
	}

	attrs := argsToAttrs(args)

	var ews errors.ErrorWithStack
	if errors.As(err, &ews) {
		attrs = append(attrs, slog.String(keys.StackTrace, formatStack(ews)))
	}

	// skip [runtime.Callers, source, this function, clog exported function]
	src := getSourceLocation(4)
	l.logAttrsWithSource(ctx, s, src, err.Error(), attrs...)
}

func (l *Logger) logWithSource(ctx context.Context, s Severity, src *sourceLocation, msg string, args ...any) {
	args = append(args, keys.SourceLocation, src)
	l.inner.Log(ctx, s, msg, args...)
}

func (l *Logger) logAttrsWithSource(ctx context.Context, s Severity, src *sourceLocation, msg string, attrs ...slog.Attr) {
	attrs = append(attrs, slog.Any(keys.SourceLocation, src))
	l.inner.LogAttrs(ctx, s, msg, attrs...)
}

func argsToAttrs(args []any) []slog.Attr {
	const badKey = "!BADKEY"

	var attrs []slog.Attr

	for len(args) > 0 {
		switch x := args[0].(type) {
		case string:
			if len(args) == 1 {
				attrs = append(attrs, slog.String(badKey, x))
				break
			} else {
				attrs = append(attrs, slog.Any(x, args[1]))
				args = args[2:]
			}
		case slog.Attr:
			attrs = append(attrs, x)
			args = args[1:]
		default:
			attrs = append(attrs, slog.Any(badKey, x))
			args = args[1:]
		}
	}

	return attrs
}

func formatStack(e errors.ErrorWithStack) string {
	return e.Error() + "\n\n" + string(e.Stack())
}
