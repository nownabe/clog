package clog

import (
	"log/slog"
	"runtime"
	"strconv"
)

// sourceLocation represents LogEntrySourceLocation.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
type sourceLocation struct {
	file     string
	line     string
	function string
}

func (s *sourceLocation) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("file", s.file),
		slog.String("line", s.line),
		slog.String("function", s.function),
	)
}

// source returns slog.Source.
// slog has built-in source mechanism, but it will be wrong when slog is wrapped.
// cf. https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/log/slog/logger.go;l=209
func getSourceLocation(skip int) *sourceLocation {
	pcs := make([]uintptr, 1)

	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return nil
	}

	fs := runtime.CallersFrames(pcs)
	f, _ := fs.Next()

	return &sourceLocation{
		file:     f.File,
		line:     strconv.Itoa(f.Line),
		function: f.Function,
	}
}
