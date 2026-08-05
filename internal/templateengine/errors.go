package templateengine

import "fmt"

// Stage identifies the operation that failed without exposing template data.
type Stage string

const (
	StageValidate Stage = "validate"
	StageRead     Stage = "read"
	StageParse    Stage = "parse"
	StageExecute  Stage = "execute"
)

// Error adds safe blueprint and file context to a rendering failure.
type Error struct {
	Stage       Stage
	BlueprintID string
	Source      string
	Destination string
	Err         error
}

func (e *Error) Error() string {
	location := e.Source
	if e.Destination != "" {
		location += " -> " + e.Destination
	}
	if location == "" {
		return fmt.Sprintf("template %s for blueprint %q: %v", e.Stage, e.BlueprintID, e.Err)
	}
	return fmt.Sprintf("template %s for blueprint %q (%s): %v", e.Stage, e.BlueprintID, location, e.Err)
}

// Unwrap exposes the underlying error for errors.Is and errors.As.
func (e *Error) Unwrap() error { return e.Err }
