package clog_test

import (
	"context"
	"testing"
	"time"

	"github.com/nownabe/clog"
)

func Test_HTTPRequest(t *testing.T) {
	t.Parallel()

	sev := "INFO"
	msg := "msg"

	tests := map[string]struct {
		r    *clog.HTTPRequest
		want map[string]any
	}{
		"empty": {
			r:    &clog.HTTPRequest{},
			want: buildWantLog(sev, msg),
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
			want: buildWantLog(sev, msg, "httpRequest", map[string]any{
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
			want: buildWantLog(sev, msg, "httpRequest", map[string]any{"requestMethod": "GET"}),
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l, w := newLogger(clog.SeverityInfo)
			l.WithHTTPRequest(tt.r).Info(ctx, msg)
			w.assertLog(t, tt.want)
		})
	}
}
