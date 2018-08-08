# scheduler
this package provides an easy way to schedule functions using an event trigger.

# scheduling a function
A timeout trigger can be defined as follows:
```go
trigger := scheduler.NewTimeoutTrigger(time.Second)
```
which can then be used to schedule functions as follows:
```go
group, ctx := scheduler.New(ctx, trigger)
```
where, `ctx` in input is the parent context and one in output is to be used to pass during scheduling functions.
```go
group.Go(ctx, func() error { 
	// your logic
	return nil
})
```
Finally wait for all functions to return using `Wait`:
```go
// wait for func executions to be over.
err := group.Wait()
```

# triggers
triggers can be of following types:
```go
// a timeout trigger allows execution on a timeout
timeoutTrigger := scheduler.NewTimeoutTrigger(timeDuration)
// trigger now causes functions to execute immediately
triggerNow := scheduler.TriggerNow()
// similarly, an event can used to trigger functions, event being ctx.Done()
triggerOnEvent := scheduler.NewContextTrigger(ctx)
```