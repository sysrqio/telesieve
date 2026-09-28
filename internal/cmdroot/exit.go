package cmdroot

import "fmt"

// ExitError signals a successful command run that should terminate with a non-zero exit code.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func newExitError(code int) error {
	return &ExitError{Code: code}
}
