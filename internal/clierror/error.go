// Package clierror carries stable process exit semantics across command adapters.
package clierror

// Error wraps a user-safe error with an exit code.
type Error struct {
	Code   int
	Err    error
	Silent bool
}

func (e *Error) Error() string { return e.Err.Error() }
func (e *Error) Unwrap() error { return e.Err }
func (e *Error) ExitCode() int { return e.Code }

// SuppressMessage reports whether the command already rendered the outcome.
func (e *Error) SuppressMessage() bool { return e.Silent }
