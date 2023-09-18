package clog_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/nownabe/clog"
	"github.com/nownabe/clog/errors"
)

type ctxUserIDKey struct{}

func InitializeClog() {
	// You can add labels to classify logs.
	appLabels := map[string]string{
		"app":     "myapp",
		"version": "1.0.0",
	}
	// Your Google Cloud Project ID.
	projectID := "my-gcp-project"

	// Custome handler to add user ID to logs.
	customHandler := func(next clog.HandleFunc) clog.HandleFunc {
		return func(ctx context.Context, r slog.Record) error {
			if ctx == nil {
				return next(ctx, r)
			}

			if userID, ok := ctx.Value(ctxUserIDKey{}).(string); ok {
				r.AddAttrs(slog.String("user_id", userID))
			}

			return next(ctx, r)
		}
	}

	// Create a custom logger.
	logger := clog.New(os.Stdout, clog.SeverityInfo, true,
		clog.WithLabels(appLabels),         // Add logger-level labels to classify logs.
		clog.WithTrace(projectID),          // Add trace and span ID to logs.
		clog.WithHandleFunc(customHandler), // Add custom log handler.
	)

	// You can add custome fields to logs.
	logger = logger.With("key1", "value1", "key2", "value2")

	clog.SetDefault(logger)
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tracer := otel.Tracer("tracer")
	ctx, span := tracer.Start(ctx, "RequestHandler")
	defer span.End()

	ctx, removeLabel := clog.ContextWithLabel(ctx, "handler", "RequestHandler")
	defer removeLabel()

	clog.Info(ctx, "request received")
	/*
	   {
	     "time": "2023-09-18T20:46:03.448753224+09:00",
	     "severity": "INFO",
	     "message": "request received",
	     "logging.googleapis.com/sourceLocation": {
	       "file": "github.com/nownabe/clog/example_test.go",
	       "line": "66",
	       "function": "main.RequestHandler"
	     },
	     "user_id": "user1",
	     "logging.googleapis.com/trace": "projects/my-gcp-project/traces/abcdabcdabcdabcdabcdabcdabcdabcd",
	     "logging.googleapis.com/spanId": "1234123412341234",
	     "logging.googleapis.com/trace_sampled": false,
	     "logging.googleapis.com/labels": {
	       "handler": "RequestHandler",
	       "app": "myapp",
	       "version": "1.0.0"
	     }
	   }
	*/
	defer clog.Info(ctx, "request processed")

	err := errors.New("something went wrong")
	if err != nil {
		clog.Err(ctx, err)
		/*
		   {
		     "time": "2023-09-18T21:14:02.980412394+09:00",
		     "severity": "ERROR",
		     "message": "something went wrong",
		     "stack_trace": "something went wrong\n\ngoroutine 0 [running]:\nmain.RequestHandler(...)\n\tgithub.com/nownabe/clog/example_test.go:90\nnet/http.HandlerFunc.ServeHTTP(...)\n\tnet/http/server.go:2136\nmain.Example(...)\n\tgithub.com/nownabe/clog/example_test.go:135\nmain.main(...)\n\tgithub.com/nownabe/clog/main.go:20\nruntime.main(...)\n\truntime/proc.go:267\n",
		     "logging.googleapis.com/sourceLocation": {
		       "file": "github.com/nownabe/clog/example_test.go",
		       "line": "92",
		       "function": "main.RequestHandler"
		     },
		     "user_id": "user1",
		     "logging.googleapis.com/trace": "projects/my-gcp-project/traces/abcdabcdabcdabcdabcdabcdabcdabcd",
		     "logging.googleapis.com/spanId": "1234123412341234",
		     "logging.googleapis.com/trace_sampled": false,
		     "logging.googleapis.com/labels": {
		       "handler": "RequestHandler",
		       "app": "myapp",
		       "version": "1.0.0"
		     }
		   }
		*/
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
		/*
		   {
		     "time": "2023-09-18T21:14:02.980418868+09:00",
		     "severity": "INFO",
		     "message": "request processed",
		     "logging.googleapis.com/sourceLocation": {
		       "file": "github.com/nownabe/clog/example_test.go",
		       "line": "116",
		       "function": "main.RequestHandler"
		     },
		     "user_id": "user1",
		     "logging.googleapis.com/trace": "projects/my-gcp-project/traces/abcdabcdabcdabcdabcdabcdabcdabcd",
		     "logging.googleapis.com/spanId": "1234123412341234",
		     "logging.googleapis.com/trace_sampled": false,
		     "logging.googleapis.com/labels": {
		       "handler": "RequestHandler",
		       "app": "myapp",
		       "version": "1.0.0"
		     }
		   }
		*/
	}

	w.Write([]byte("Hello, World!"))
}

func Example() {
	ctx := context.Background()
	ctx = withTraceSpan(ctx)
	ctx = context.WithValue(ctx, ctxUserIDKey{}, "user1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		panic(err)
	}
	InitializeClog()

	clog.Notice(ctx, "app started")
	/*
	   {
	     "time": "2023-09-18T20:46:03.448702881+09:00",
	     "severity": "NOTICE",
	     "message": "app started",
	     "logging.googleapis.com/sourceLocation": {
	       "file": "github.com/nownabe/clog/example_test.go",
	       "line": "92",
	       "function": "main.Example"
	     },
	     "user_id": "user1",
	     "logging.googleapis.com/trace": "projects/my-gcp-project/traces/abcdabcdabcdabcdabcdabcdabcdabcd",
	     "logging.googleapis.com/spanId": "1234123412341234",
	     "logging.googleapis.com/trace_sampled": false,
	     "logging.googleapis.com/labels": {
	       "app": "myapp",
	       "version": "1.0.0"
	     }
	   }
	*/

	(http.HandlerFunc(RequestHandler)).ServeHTTP(httptest.NewRecorder(), req)

	clog.Notice(ctx, "app finished")
	/*
	   {
	     "time": "2023-09-18T21:14:02.980422814+09:00",
	     "severity": "NOTICE",
	     "message": "app finished",
	     "logging.googleapis.com/sourceLocation": {
	       "file": "github.com/nownabe/clog/tmp/example.go",
	       "line": "178",
	       "function": "main.Example"
	     },
	     "user_id": "user1",
	     "logging.googleapis.com/trace": "projects/my-gcp-project/traces/abcdabcdabcdabcdabcdabcdabcdabcd",
	     "logging.googleapis.com/spanId": "1234123412341234",
	     "logging.googleapis.com/trace_sampled": false,
	     "logging.googleapis.com/labels": {
	       "app": "myapp",
	       "version": "1.0.0"
	     }
	   }
	*/
}

func withTraceSpan(ctx context.Context) context.Context {
	traceID, err := trace.TraceIDFromHex("abcdabcdabcdabcdabcdabcdabcdabcd")
	if err != nil {
		panic(err)
	}
	spanID, err := trace.SpanIDFromHex("1234123412341234")
	if err != nil {
		panic(err)
	}
	cfg := trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	}
	sc := trace.NewSpanContext(cfg)

	return trace.ContextWithSpanContext(ctx, sc)
}
