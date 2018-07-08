# scheduler
a func scheduler interface

# deferred execution
an implementation of scheduler providing deferred execution of
a func. Such scheduling happens within the framework of an error
group such that the fist func to error out with non-nil error
causes all other pending executions to get cancelled.
