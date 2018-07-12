package scheduler

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Trigger struct {
	value interface{}
}

// deferredErrGroupScheduler implements Scheduler interface providing deferred execution for a func.
type deferredErrGroupScheduler struct {
	// trig is trigger indicating how user wants function to be triggered.
	// trigger is either a timeout or an event defined by cancellation of a context.
	// Use NewTimeoutTrigger or NewEvent to create options in New method.
	trig Trigger
	// m is map of cancel func that can be accessed using keys.
	m map[string]context.CancelFunc
	// mu is a lock for modifying above map.
	mu sync.Mutex
	// eg is error group
	eg *errgroup.Group
}
