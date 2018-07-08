// package scheduler defines func scheduling interface and provides
// implementation for deferred execution.
package scheduler

import (
	"context"
)

// Scheduler defines interface to schedule execution of a func.
type Scheduler interface {
	// Go executes function and returns a key which can be used to cancel func execution.
	Go(ctx context.Context, f func() error) string
	// Cancel cancels execution of func if pending execution.
	Cancel(key string) error
	// Wait waits for all functions dispatched using Go to finish.
	Wait() error
}
