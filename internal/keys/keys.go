package keys

/*
These keys in JSON log entries are automatically extracted into the LogEntry structure.
See https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
*/
const (
	HTTPRequest    = "httpRequest"
	Severity       = "severity"
	SourceLocation = "logging.googleapis.com/sourceLocation"
	StackTrace     = "stack_trace"
)
