package clog

import (
	"log/slog"
)

/*
Severity is the severity of the log event.
These severities are defined in the Cloud Logging API v2 as an enum.

See following links.
https://github.com/googleapis/googleapis/blob/master/google/logging/type/log_severity.proto
https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity

Though these packages provide predefined severity constants, we don't use them not to depend on external package just for it.
https://pkg.go.dev/google.golang.org/genproto/googleapis/logging/type#LogSeverity
https://pkg.go.dev/cloud.google.com/go/logging#Severity
*/
type Severity = slog.Level

const (
	// SeverityDefault indicates that the log entry has no assigned severity level.
	SeverityDefault = Severity(0)
	// SeverityDebug indicates debug or trace information.
	SeverityDebug = Severity(100)
	// SeverityInfo indicates routine information, such as ongoing status or performance.
	SeverityInfo = Severity(200)
	// SeverityNotice indicates normal but significant events, such as start up, shut down, or a configuration change.
	SeverityNotice = Severity(300)
	// SeverityWarning indicates warning events that might cause problems.
	SeverityWarning = Severity(400)
	// SeverityError indicates error events that are likely to cause problems.
	SeverityError = Severity(500)
	// SeverityCritical indicates critical events that cause more severe problems or outages.
	SeverityCritical = Severity(600)
	// SeverityAlert indicates that a person must take an action immediately.
	SeverityAlert = Severity(700)
	// SeverityEmergency indicates that one or more systems are unusable.
	SeverityEmergency = Severity(800)
)

func severityString(s Severity) string {
	switch s {
	case SeverityDefault:
		return "DEFAULT"
	case SeverityDebug:
		return "DEBUG"
	case SeverityInfo:
		return "INFO"
	case SeverityNotice:
		return "NOTICE"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	case SeverityAlert:
		return "ALERT"
	case SeverityEmergency:
		return "EMERGENCY"
	}

	return "DEFAULT"
}
