# go.nownabe.dev/clog

[![PkgGoDev](https://pkg.go.dev/badge/go.nownabe.dev/clog)](https://pkg.go.dev/go.nownabe.dev/clog)
[![License](https://img.shields.io/github/license/nownabe/clog.svg?style=popout)](https://github.com/nownabe/clog/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nownabe/clog)](https://goreportcard.com/report/github.com/nownabe/clog)
[![codecov](https://codecov.io/gh/nownabe/clog/graph/badge.svg?token=hbX8ZaXRRC)](https://codecov.io/gh/nownabe/clog)

clog is a structured logger optimized for [Cloud Logging] based on [log/slog](https://pkg.go.dev/log/slog).

clog supports following special fields in JSON log entries:

- `severity`
- `message`
- `httpRequest`
- `time`
- `logging.googleapis.com/insertId`
- `logging.googleapis.com/labels`
- `logging.googleapis.com/operation`
- `logging.googleapis.com/sourceLocation`
- `logging.googleapis.com/spanId`
- `logging.googleapis.com/trace`
- `logging.googleapis.com/trace_sampled`
- `logging.googleapis.com/stack_trace`

See [Cloud Logging documentation] and [Cloud Error Reporting documentation] for more details.

## Severity

clog uses Severity in the `severity` field instead of log levels. 8 severities are supported:

- DEBUG
- INFO
- NOTICE
- WARNING
- ERROR
- CRITICAL
- ALERT
- EMERGENCY

## Usage

Each severity has three methods like Info, Infof, and InfoErr.

```go
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
```

See [Examples](https://pkg.go.dev/go.nownabe.dev/clog#pkg-examples) for more details.

[Cloud Logging]: https://cloud.google.com/logging
[Cloud Logging documentation]: https://cloud.google.com/logging/docs/structured-logging
[Cloud Error Reporting documentation]: https://cloud.google.com/error-reporting/docs/formatting-error-messages
