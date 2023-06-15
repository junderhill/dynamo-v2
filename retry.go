package helixddb

import (
	"context"
	"errors"
	"time"

	smithytime "github.com/aws/smithy-go/time"
	"github.com/cenkalti/backoff/v4"
)

// RetryTimeout defines the maximum amount of time that requests will
// attempt to automatically retry for. In other words, this is the maximum
// amount of time that dynamo operations will block.
// RetryTimeout is only considered by methods that do not take a context.
// Higher values are better when using tables with lower throughput.
var RetryTimeout = 1 * time.Minute

func defaultContext() (context.Context, context.CancelFunc) {
	if RetryTimeout == 0 {
		return context.Background(), (func() {})
	}
	return context.WithDeadline(context.Background(), time.Now().Add(RetryTimeout))
}

func retry(ctx context.Context, f func() error) error {
	var err error
	var next time.Duration
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	for {
		if err = f(); err == nil {
			return nil
		}

		if !canRetry(err) {
			return err
		}

		if next = b.NextBackOff(); next == backoff.Stop {
			return err
		}

		if err = smithytime.SleepWithContext(ctx, next); err != nil {
			return err
		}
	}
}

// errRetry is a sentinel error to retry, should never be returned to user
var errRetry = errors.New("dynamo: retry")

func canRetry(err error) bool {
	return errors.Is(err, errRetry)
}
