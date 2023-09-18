package clog_test

import (
	"context"

	"github.com/nownabe/clog"
)

func DoLongRunningOperation(ctx context.Context) {
	ctx, end := clog.StartOperation(ctx, clog.SeverityInfo, "long-running operation started", "long-running", "my-app")
	defer end("long-running operation ended")

	clog.Info(ctx, "long-running operation is running")
}

func ExampleStartOperation() {
	DoLongRunningOperation(context.Background())
	// Logs are like:
	// {"severity":"INFO", "message":"long-running operation started",
	//  "logging.googleapis.com/operation":{"id":"long-running","producer":"my-app","first":true}, ...}
	// {"severity":"INFO", "message":"long-running operation is running",
	//  "logging.googleapis.com/operation":{"id":"long-running","producer":"my-app"}, ...}
	// {"severity":"INFO", "message":"long-running operation ended",
	//  "logging.googleapis.com/operation":{"id":"long-running","producer":"my-app","last":true}, ...}
}
