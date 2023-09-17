package clog

import (
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
	return defaultLogger.Load().(*Logger)
}
