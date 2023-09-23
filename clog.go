package clog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"

	"go.nownabe.dev/clog/internal/keys"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(New(os.Stdout, SeverityInfo, true))
}

// SetDefault makes l the default Logger.
func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

// SetOptions sets options to the default Logger.
func SetOptions(opts ...Option) {
	l := Default()

	h := l.inner.Handler()
	for _, o := range opts {
		h = o.apply(h)
	}

	SetDefault(&Logger{slog.New(h)})
}

// Default returns the default Logger.
func Default() *Logger {
	l, ok := defaultLogger.Load().(*Logger)
	if !ok {
		panic("clog: default logger is not *Logger")
	}
	return l
}

// Debug logs at SeverityDebug.
func Debug(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityDebug, msg, args...)
}

// Debugf logs formatted in the manner of fmt.Printf at SeverityDebug.
func Debugf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityDebug, fmt.Sprintf(format, a...))
}

// DebugErr logs an error at SeverityDebug.
func DebugErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityDebug, err, args...)
}

// Info logs at SeverityInfo.
func Info(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityInfo, msg, args...)
}

// Infof logs formatted in the manner of fmt.Printf at SeverityInfo.
func Infof(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityInfo, fmt.Sprintf(format, a...))
}

// InfoErr logs an error at SeverityInfo.
func InfoErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityInfo, err, args...)
}

// Notice logs at SeverityNotice.
func Notice(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityNotice, msg, args...)
}

// Noticef logs formatted in the manner of fmt.Printf at SeverityNotice.
func Noticef(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityNotice, fmt.Sprintf(format, a...))
}

// NoticeErr logs an error at SeverityNotice.
func NoticeErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityNotice, err, args...)
}

// Warning logs at SeverityWarning.
func Warning(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityWarning, msg, args...)
}

// Warningf logs formatted in the manner of fmt.Printf at SeverityWarning.
func Warningf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityWarning, fmt.Sprintf(format, a...))
}

// WarningErr logs an error at SeverityWarning.
func WarningErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityWarning, err, args...)
}

// Error logs at SeverityError.
func Error(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityError, msg, args...)
}

// Errorf logs formatted in the manner of fmt.Printf at SeverityError.
func Errorf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityError, fmt.Sprintf(format, a...))
}

// ErrorErr logs an error at SeverityError.
func ErrorErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityError, err, args...)
}

// Critical logs at SeverityCriticaDefault().
func Critical(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityCritical, msg, args...)
}

// Criticalf logs formatted in the manner of fmt.Printf at SeverityCriticaDefault().
func Criticalf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityCritical, fmt.Sprintf(format, a...))
}

// CriticalErr logs an error at SeverityCriticaDefault().
func CriticalErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityCritical, err, args...)
}

// Alert logs at SeverityAlert.
func Alert(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityAlert, msg, args...)
}

// Alertf logs formatted in the manner of fmt.Printf at SeverityAlert.
func Alertf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityAlert, fmt.Sprintf(format, a...))
}

// AlertErr logs an error at SeverityAlert.
func AlertErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityAlert, err, args...)
}

// Emergency logs at SeverityEmergency.
func Emergency(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, SeverityEmergency, msg, args...)
}

// Emergencyf logs formatted in the manner of fmt.Printf at SeverityEmergency.
func Emergencyf(ctx context.Context, format string, a ...any) {
	Default().log(ctx, SeverityEmergency, fmt.Sprintf(format, a...))
}

// EmergencyErr logs an error at SeverityEmergency.
func EmergencyErr(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityEmergency, err, args...)
}

// Log emits a log record with the current time and the given level and message.
func Log(ctx context.Context, s Severity, msg string, args ...any) {
	Default().log(ctx, s, msg, args...)
}

// Err is a shorthand for ErrorErr.
func Err(ctx context.Context, err error, args ...any) {
	Default().err(ctx, SeverityError, err, args...)
}

// Enabled reports whether the Logger emits log records at the given context and leveDefault().
func Enabled(ctx context.Context, s Severity) bool {
	return Default().Enabled(ctx, s)
}

// With returns a Logger that includes the given attributes in each output operation.
func With(args ...any) *Logger {
	return Default().With(args...)
}

// HTTPReq emits a log with the given [HTTPRequest].
// If status >= 500, the log is at SeverityError.
// Otherwise, the log is at SeverityInfo.
func HTTPReq(ctx context.Context, req *HTTPRequest, msg string, args ...any) {
	s := SeverityInfo
	if req.Status >= 500 {
		s = SeverityError
	}
	args = append(args, keys.HTTPRequest, req)
	Default().log(ctx, s, msg, args...)
}

// WithInsertID returns a Logger that includes the given insertId in each output operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func WithInsertID(id string) *Logger {
	return Default().WithInsertID(id)
}

// StartOperation returns a new context and a function to end the opration, starting the operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
func StartOperation(ctx context.Context, s Severity, msg, id, producer string) (context.Context, func(msg string)) {
	return Default().startOperation(ctx, s, msg, id, producer)
}
