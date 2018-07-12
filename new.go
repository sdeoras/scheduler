package scheduler

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

// New creates a new scheduler based on an input context and trigger mode.
// Scheduler will cancel all pending jobs if the input context is cancelled.
// Scheduler will trigger jobs based on trigger.
// It returns a context that should be used in each function.
// This output context is cancelled on first error in any of the scheduled jobs.
func New(ctx context.Context, trig Trigger) (Scheduler, context.Context) {
	s := new(deferredErrGroupScheduler)
	eg, ctx := errgroup.WithContext(ctx)
	s.eg, s.trig, s.m = eg, trig, make(map[string]context.CancelFunc)
	return s, ctx
}

// NewTimeoutTrigger provides a new trigger based on a timeout.
func NewTimeoutTrigger(dur time.Duration) Trigger {
	return Trigger{value: dur}
}

// NewContextTrigger provides a new trigger based on a context.
// Trigger fires when input context context is done.
func NewContextTrigger(ctx context.Context) Trigger {
	return Trigger{value: ctx}
}

// NewChannelTrigger provides a new trigger based on a channel.
// Trigger fires when input chan is read.
func NewChannelTrigger(c chan struct{}) Trigger {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-c:
			cancel()
		}
	}()
	return Trigger{value: ctx}
}
