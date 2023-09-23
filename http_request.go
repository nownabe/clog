package clog

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"
)

/*
HTTPRequest represents HttpRequest.
See these links: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest,
https://github.com/googleapis/googleapis/blob/master/google/logging/type/http_request.proto.

If HTTPRequest is empty, it will be omitted.
cf. https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/log/slog/record.go;l=118
*/
type HTTPRequest struct {
	RequestMethod                  string
	RequestURL                     string
	RequestSize                    int64
	Status                         int
	ResponseSize                   int64
	UserAgent                      string
	RemoteIP                       string
	ServerIP                       string
	Referer                        string
	Latency                        time.Duration
	CacheLookup                    bool
	CacheHit                       bool
	CacheValidatedWithOriginServer bool
	CacheFillBytes                 int64
	Protocol                       string
}

// LogValue returns slog.Value.
func (r *HTTPRequest) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 15)

	if r.RequestMethod != "" {
		attrs = append(attrs, slog.String("requestMethod", r.RequestMethod))
	}
	if r.RequestURL != "" {
		attrs = append(attrs, slog.String("requestUrl", r.RequestURL))
	}
	if r.RequestSize != 0 {
		attrs = append(attrs, slog.String("requestSize", strconv.FormatInt(r.RequestSize, 10)))
	}
	if r.Status != 0 {
		attrs = append(attrs, slog.Int("status", r.Status))
	}
	if r.ResponseSize != 0 {
		attrs = append(attrs, slog.String("responseSize", strconv.FormatInt(r.ResponseSize, 10)))
	}
	if r.UserAgent != "" {
		attrs = append(attrs, slog.String("userAgent", r.UserAgent))
	}
	if r.RemoteIP != "" {
		attrs = append(attrs, slog.String("remoteIp", r.RemoteIP))
	}
	if r.ServerIP != "" {
		attrs = append(attrs, slog.String("serverIp", r.ServerIP))
	}
	if r.Referer != "" {
		attrs = append(attrs, slog.String("referer", r.Referer))
	}
	if r.Latency != 0 {
		// https://protobuf.dev/reference/protobuf/google.protobuf/#duration
		attrs = append(attrs, slog.String("latency", fmt.Sprintf("%.9fs", r.Latency.Seconds())))
	}
	if r.CacheLookup {
		attrs = append(attrs, slog.Bool("cacheLookup", r.CacheLookup))
		attrs = append(attrs, slog.Bool("cacheHit", r.CacheHit))
		if r.CacheHit {
			attrs = append(attrs, slog.Bool("cacheValidatedWithOriginServer", r.CacheValidatedWithOriginServer))
		}
	}
	if r.CacheFillBytes != 0 {
		attrs = append(attrs, slog.String("cacheFillBytes", strconv.FormatInt(r.CacheFillBytes, 10)))
	}
	if r.Protocol != "" {
		attrs = append(attrs, slog.String("protocol", r.Protocol))
	}

	return slog.GroupValue(attrs...)
}

func (r *HTTPRequest) msg() string {
	msg := strings.Join(
		slices.DeleteFunc(
			[]string{r.RequestMethod, r.RequestURL, r.Protocol},
			func(s string) bool { return s == "" }),
		" ")
	if msg == "" {
		msg = "HTTP request"
	}
	return msg
}
