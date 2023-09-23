package clog_test

import (
	"context"
	"testing"
	"time"

	"go.nownabe.dev/clog"
)

func TestLogger_HTTPReq(t *testing.T) {
	t.Parallel()

	msg := "msg"

	tests := map[string]struct {
		r    *clog.HTTPRequest
		want map[string]any
	}{
		"empty": {
			r:    &clog.HTTPRequest{},
			want: buildWantLog("INFO", msg),
		},
		"full": {
			r: &clog.HTTPRequest{
				RequestMethod:                  "GET",
				RequestURL:                     "https://example.com/foo",
				RequestSize:                    123,
				Status:                         200,
				ResponseSize:                   456,
				UserAgent:                      "clog",
				RemoteIP:                       "203.0.113.1",
				ServerIP:                       "203.0.113.2",
				Referer:                        "https://example.com/referer",
				Latency:                        123*time.Second + 456*time.Nanosecond,
				CacheLookup:                    true,
				CacheHit:                       true,
				CacheValidatedWithOriginServer: true,
				CacheFillBytes:                 789,
				Protocol:                       "HTTP/1.1",
			},
			want: buildWantLog("INFO", msg, "httpRequest", map[string]any{
				"requestMethod":                  "GET",
				"requestUrl":                     "https://example.com/foo",
				"requestSize":                    "123",
				"status":                         200,
				"responseSize":                   "456",
				"userAgent":                      "clog",
				"remoteIp":                       "203.0.113.1",
				"serverIp":                       "203.0.113.2",
				"referer":                        "https://example.com/referer",
				"latency":                        "123.000000456s",
				"cacheLookup":                    true,
				"cacheHit":                       true,
				"cacheValidatedWithOriginServer": true,
				"cacheFillBytes":                 "789",
				"protocol":                       "HTTP/1.1",
			}),
		},
		"only requestMethod": {
			r:    &clog.HTTPRequest{RequestMethod: "GET"},
			want: buildWantLog("INFO", msg, "httpRequest", map[string]any{"requestMethod": "GET"}),
		},
		"internal error": {
			r:    &clog.HTTPRequest{RequestMethod: "GET", Status: 500},
			want: buildWantLog("ERROR", msg, "httpRequest", map[string]any{"requestMethod": "GET", "status": 500}),
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l, w := newLogger(clog.SeverityInfo)
			l.HTTPReq(ctx, tt.r, msg)
			w.assertLog(t, tt.want)
		})
	}
}

func ExampleHTTPReq() {
	req := &clog.HTTPRequest{
		RequestMethod:                  "GET",
		RequestURL:                     "https://example.com/foo",
		RequestSize:                    123,
		Status:                         200,
		ResponseSize:                   456,
		UserAgent:                      "clog",
		RemoteIP:                       "203.0.113.1",
		ServerIP:                       "203.0.113.2",
		Referer:                        "https://example.com/referer",
		Latency:                        123*time.Second + 456*time.Nanosecond,
		CacheLookup:                    true,
		CacheHit:                       true,
		CacheValidatedWithOriginServer: true,
		CacheFillBytes:                 789,
		Protocol:                       "HTTP/1.1",
	}
	clog.HTTPReq(context.Background(), req, "GET /foo")
}
