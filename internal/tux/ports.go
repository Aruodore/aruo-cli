package tux

import "context"

// Prompter resolves human decisions without exposing a terminal library.
type Prompter interface {
	Input(context.Context, InputRequest) (string, error)
	Secret(context.Context, SecretRequest) (string, error)
	Confirm(context.Context, ConfirmRequest) (bool, error)
	Select(context.Context, SelectRequest) (OptionID, error)
	MultiSelect(context.Context, MultiSelectRequest) ([]OptionID, error)
	// Guide runs a multi-screen flow as one session, letting the user move
	// backward across screens rather than only within one prompt.
	Guide(context.Context, []Step) (Answers, error)
}

// Presenter renders semantic results and diagnostics.
type Presenter interface {
	Message(context.Context, Message) error
	Table(context.Context, Table) error
	Diagnostic(context.Context, Diagnostic) error
}

// ProgressSink accepts ordered task lifecycle events.
type ProgressSink interface {
	Emit(context.Context, TaskEvent) error
}

// Session exposes the resolved interaction services for one invocation.
type Session interface {
	Capabilities() Capabilities
	Policy() Policy
	Prompter() Prompter
	Presenter() Presenter
	Progress() ProgressSink
	Close() error
}
