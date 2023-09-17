package keys

/*
These keys in JSON log entries are automatically extracted into the LogEntry structure.
See https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
*/
const (
	apiPrefix = "logging.googleapis.com/"

	HTTPRequest    = "httpRequest"
	InsertID       = apiPrefix + "insertId"
	Labels         = apiPrefix + "labels"
	Severity       = "severity"
	SourceLocation = apiPrefix + "sourceLocation"
	StackTrace     = "stack_trace"
)
