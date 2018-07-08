package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// test if func gets executed after timeout.
func Test_Exec_NormalBehavior(t *testing.T) {
	startTime := time.Now()
	execTime := time.Now()
	// new context with a 2 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	// get a new deferredErrGroupScheduler with a 1 second timeout.
	// functions will be executed after one second delay.
	group, ctx := NewDeferredErrGroup(ctx, time.Second)

	// defer execute func setting value of execTime to be 1.
	group.Go(ctx, func() error { execTime = time.Now(); return nil })

	// wait for func executions to be over.
	err := group.Wait()

	// there should be no error.
	if err != nil {
		t.Fatal(err)
	}

	// and func should have executed as expected.
	if execTime.Sub(startTime) < time.Second {
		t.Fatal("func did not run as expected")
	}
}

// test if global timeout cancells func execution.
func Test_Exec_GloabalTimeout(t *testing.T) {
	var c int

	// new context with a 1 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	// get a new deferredErrGroupScheduler with a 2 second timeout.
	// functions will be executed after one second delay.
	group, ctx := NewDeferredErrGroup(ctx, 2*time.Second)

	// defer execute func.
	group.Go(ctx, func() error { c = 1; return nil })

	// wait for executions of all funcs to be over.
	err := group.Wait()

	// the error should be not be nil since global context is over before deferred timeout.
	if err == nil {
		t.Fatal(err)
	}

	// the error should be of known value.
	if err != FuncExecCancelled {
		t.Fatal("did not received expected error")
	}

	// func should not have executed
	if c != 0 {
		t.Fatal("func did not execute as expected")
	}
}

// test if cancelling global context cancells the func execution
func Test_Exec_GlobalCancel(t *testing.T) {
	var c int
	// new context with a cancel func
	ctx, cancel := context.WithCancel(context.Background())

	// get a new deferredErrGroupScheduler with a 1 second timeout.
	// functions will be executed after one second delay.
	group, ctx := NewDeferredErrGroup(ctx, time.Second)

	// defer execute func that sets value of c to 1.
	group.Go(ctx, func() error { c = 1; return nil })

	// cancel global context.
	cancel()

	// wait for deferred execution func calls to be over.
	err := group.Wait()

	// error should not be nil since func was cancelled using global context canceller.
	if err == nil {
		t.Fatal(err)
	}

	// error should be of expected value.
	if err != FuncExecCancelled {
		t.Fatal("did not received expected error")
	}

	// func should not have executed.
	if c != 0 {
		t.Fatal("func did not run as expected")
	}
}

// cancel one func but let others run
func Test_Exec_CancelOneFunc(t *testing.T) {
	var c [2]int

	// new context with a 2 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	// get a new deferredErrGroupScheduler with a 1 second timeout.
	// functions will be executed after one second delay.
	group, ctx := NewDeferredErrGroup(ctx, time.Second)

	// defer execute func setting value of c to be 1.
	key := group.Go(ctx, func() error { c[0] = 1; return nil })
	group.Go(ctx, func() error { c[1] = 1; return nil })

	// cancel execution of first func.
	if err := group.Cancel(key); err != nil {
		t.Fatal(err)
	}

	// wait for func executions to be over.
	err := group.Wait()

	// error should not be nil.
	if err == nil {
		t.Fatal(err)
	}

	if err != FuncExecCancelled {
		t.Fatal("error value not as expected")
	}

	// and func should have executed such that only second func ran.
	if c[0] != 0 && c[1] != 1 {
		t.Fatal("func did not run as expected")
	}
}

// one of the funcs errors out. see if everything gets cancelled.
func Test_Exec_OneFuncErrorsOut(t *testing.T) {
	var c [2]int

	// new context with a 2 second timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	// get a new deferredErrGroupScheduler with a 1 second timeout.
	// functions will be executed after one second delay.
	group, ctx := NewDeferredErrGroup(ctx, time.Second)

	// defer execute func setting value of c to be 1.
	group.Go(ctx, func() error { c[0] = 1; return fmt.Errorf("to err is human") })
	group.Go(ctx, func() error { time.Sleep(time.Millisecond * 100); c[1] = 1; return nil })

	// wait for func executions to be over.
	err := group.Wait()

	// error should not be nil.
	if err == nil {
		t.Fatal(err)
	}

	if err == FuncExecCancelled {
		t.Fatal("error value not as expected")
	}

	// and func should have executed such that only second func ran.
	if c[0] != 1 && c[1] != 0 {
		t.Fatal("func did not run as expected")
	}
}
