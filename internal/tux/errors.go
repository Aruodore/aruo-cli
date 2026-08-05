package tux

import "errors"

var (
	// ErrCancelled reports an explicit user cancellation.
	ErrCancelled = errors.New("interaction cancelled")
	// ErrEndOfInput reports that the input stream closed before a decision.
	ErrEndOfInput = errors.New("end of input")
	// ErrUnavailable reports that the selected interaction cannot run in this session.
	ErrUnavailable = errors.New("interaction unavailable")
)
