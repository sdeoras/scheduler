package scheduler

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

// NewDeferredErrGroup provides a new instance of Scheduler and a contexct derived from the input context.
func NewDeferredErrGroup(ctx context.Context, dur time.Duration) (Scheduler, context.Context) {
	s := new(deferredErrGroupScheduler)
	eg, ctx := errgroup.WithContext(ctx)
	s.eg, s.dur, s.m = eg, dur, make(map[string]context.CancelFunc)
	return s, ctx
}
