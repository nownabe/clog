package clog

import (
	"context"
	"os"
	"sync/atomic"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(New(os.Stdout, SeverityInfo, false))
}

// SetDefault makes l the default Logger.
func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

// Default returns the default Logger.
func Default() *Logger {
	l, ok := defaultLogger.Load().(*Logger)
	if !ok {
		panic("clog: default logger is not *Logger")
	}
	return l
}

// StartOperation returns a new context and a function to end the opration, starting the operation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
func StartOperation(ctx context.Context, s Severity, msg, id, producer string) (context.Context, func(msg string)) {
	return Default().startOperation(ctx, s, msg, id, producer)
}
