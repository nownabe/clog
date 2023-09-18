/*
Package clog is a structured logger optimized for [Cloud Logging] based on [log/slog].

clog supports following special fields in JSON log entries:

  - severity
  - message
  - httpRequest
  - time
  - logging.googleapis.com/insertId
  - logging.googleapis.com/labels
  - logging.googleapis.com/operation
  - logging.googleapis.com/sourceLocation
  - logging.googleapis.com/spanId
  - logging.googleapis.com/trace
  - logging.googleapis.com/trace_sampled
  - logging.googleapis.com/stack_trace

See [Cloud Logging documentation] and [Cloud Error Reporting documentation] for more details.

# Severity

clod uses [Severity] in the "severity" field instead of log levels. 8 severities are supported:

  - DEBUG
  - INFO
  - NOTICE
  - WARNING
  - ERROR
  - CRITICAL
  - ALERT
  - EMERGENCY

# Usage

Each severity has three methods like [Info], [Infof], and [InfoErr].

	clog.Info(ctx, "simple logging with args", "key", "value")
	// {"severity":"INFO", "message":"simple logging with args", "key":"value",
	//  "time": "...", "logging.googleapis.com/sourceLocation": {...}}

	clog.Noticef(ctx, "formatted message %s %s", "like", "fmt.Printf")
	// {"severity":"NOTICE", "message":"formatted message like fmt.Printf",
	//  "time": "...", "logging.googleapis.com/sourceLocation": {...}}

	clog.CriticalErr(ctx, errors.New("critical error!!"), "key", "value")
	// {"severity":"CRITICAL", "message":"critical error!!", "key":"value",
	//  "time": "...", "logging.googleapis.com/sourceLocation": {...}}

	// clog.ErrorErr has a shorthand clog.Err.
	clog.Err(ctx, errors.New("error!"))
	// {"severity":"ERROR", "message":"error!",
	//  "time": "...", "logging.googleapis.com/sourceLocation": {...}}

See [Examples] for more details.

[Cloud Logging]: https://cloud.google.com/logging
[Cloud Logging documentation]: https://cloud.google.com/logging/docs/structured-logging
[Cloud Error Reporting documentation]: https://cloud.google.com/error-reporting/docs/formatting-error-messages
*/
package clog // import "go.nownabe.dev/clog"
