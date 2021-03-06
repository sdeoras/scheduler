package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Func implements Scheduler interface to provide a functional literal for use in errgroup.
// Such func will execute the input func deferred in time.
func (d *deferredErrGroupScheduler) Go(ctx context.Context,
	f func() error) string {

	// build new key
	key := uuid.New().String()
	// derive new context from input context
	ctx, cancelFunc := context.WithCancel(ctx)

	// map key to the cancel func
	d.mu.Lock()
	d.m[key] = cancelFunc
	d.mu.Unlock()

	var trig context.Context
	var trigCancelFunc context.CancelFunc

	switch v := d.trig.value.(type) {
	case time.Duration:
		// create a context based on timeout
		trig, trigCancelFunc = context.WithTimeout(context.Background(), v)
		// spawn a go-routine that will call cancel func once triggered
		go func() {
			select {
			case <-trig.Done():
				trigCancelFunc()
			}
		}()
	case context.Context:
		trig = v
	}

	// spawn func within an error group.
	d.eg.Go(func() error {
		select {
		case <-trig.Done():
			err := f()
			d.Cancel(key)
			return err
		case <-ctx.Done():
			return FuncExecCancelled
		}
	})

	return key
}

// Cancel cancels execution of func associated with the key if it is pending.
func (d *deferredErrGroupScheduler) Cancel(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if f, ok := d.m[key]; !ok {
		return KeyNotFound
	} else {
		if f != nil {
			f()
		}
	}

	delete(d.m, key)
	return nil
}

// Wait waits for all function executions to be over and performs cleanup.
func (d *deferredErrGroupScheduler) Wait() error {
	err := d.eg.Wait()

	d.mu.Lock()
	defer d.mu.Unlock()
	for key, f := range d.m {
		if f != nil {
			f()
		}
		delete(d.m, key)
	}

	return err
}
