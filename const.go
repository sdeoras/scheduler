package scheduler

// Error is used to build error constants.
type Error string

// Error implements Error interface for Error type.
func (e Error) Error() string {
	return string(e)
}

const (
	FuncExecCancelled Error = "function execution was cancelled by context"
	KeyNotFound       Error = "key not found"
)
