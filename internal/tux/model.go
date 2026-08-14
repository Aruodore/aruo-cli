// Package tux defines Aruo's terminal interaction contracts.
package tux

import "time"

// Mode selects an interaction contract independently of terminal styling.
type Mode uint8

const (
	ModeInteractive Mode = iota
	ModeNonInteractive
	ModeMachine
)

// ColorLevel describes the terminal color profile available to a renderer.
type ColorLevel uint8

const (
	ColorNone ColorLevel = iota
	ColorANSI16
	ColorANSI256
	ColorTrueColor
)

// FeaturePolicy controls capability-dependent presentation features.
type FeaturePolicy uint8

const (
	FeatureAuto FeaturePolicy = iota
	FeatureAlways
	FeatureNever
)

// OutputFormat selects the stable result representation.
type OutputFormat string

const (
	OutputHuman    OutputFormat = "human"
	OutputPlain    OutputFormat = "plain"
	OutputJSON     OutputFormat = "json"
	OutputJSONL    OutputFormat = "jsonl"
	OutputYAML     OutputFormat = "yaml"
	OutputSARIF    OutputFormat = "sarif"
	OutputMarkdown OutputFormat = "markdown"
)

// Capabilities records independently detected terminal features.
type Capabilities struct {
	InputTTY         bool
	OutputTTY        bool
	ErrorTTY         bool
	Width            int
	Height           int
	Color            ColorLevel
	Unicode          bool
	Hyperlinks       bool
	CursorAddressing bool
	AlternateScreen  bool
	CI               bool
	SSH              bool
	Multiplexer      string
}

// Policy contains resolved user and environment presentation choices.
type Policy struct {
	Mode       Mode
	Format     OutputFormat
	Accessible bool
	Input      FeaturePolicy
	Color      FeaturePolicy
	Unicode    FeaturePolicy
	Motion     FeaturePolicy
	Pager      FeaturePolicy
}

// MessageKind gives prose a semantic role without prescribing styling.
type MessageKind string

const (
	MessageInfo    MessageKind = "info"
	MessageSuccess MessageKind = "success"
	MessageWarning MessageKind = "warning"
	MessageNote    MessageKind = "note"
)

// Message is one durable human-facing statement.
type Message struct {
	Kind MessageKind
	Text string
}

// Alignment describes cell alignment independently of a table renderer.
type Alignment uint8

const (
	AlignLeft Alignment = iota
	AlignRight
)

// Column describes one semantic table field. Lower priorities are retained first.
type Column struct {
	ID        string
	Heading   string
	Alignment Alignment
	Priority  int
	Sensitive bool
}

// Table contains renderer-neutral tabular data.
type Table struct {
	Columns []Column
	Rows    [][]string
}

// Diagnostic is one actionable user-facing problem.
type Diagnostic struct {
	Code       string
	Summary    string
	Detail     string
	Effect     string
	Suggestion string
}

// OptionID is a stable choice identifier independent of its displayed label.
type OptionID string

// Option is one selectable value.
type Option struct {
	ID          OptionID
	Label       string
	Description string
	// Color is an optional terminal-readable truecolor hex hint.
	// Adapters with no color rendering ignore it. Color-capable adapters
	// apply it only when the active color policy and capability permit
	// color at all; it never substitutes for Label as the only way to
	// distinguish options.
	Color    string
	Disabled bool
}

// InputRequest describes one text decision.
type InputRequest struct {
	ID          string
	Label       string
	Description string
	Example     string
	Placeholder string
	Suggestions []string
	Default     *string
	Optional    bool
	Validate    func(string) error
}

// SecretRequest describes one value that must not be echoed or logged.
type SecretRequest struct {
	ID          string
	Label       string
	Description string
	Optional    bool
	Validate    func(string) error
}

// ConfirmRequest describes a yes-or-no decision.
type ConfirmRequest struct {
	ID          string
	Label       string
	Description string
	Default     bool
}

// SelectRequest describes one selection from stable options.
type SelectRequest struct {
	ID          string
	Label       string
	Description string
	Options     []Option
	Default     *OptionID
	// HighlightActive uses a filled, terminal-adaptive active state instead
	// of the default underline. Reserve it for short categorical choices.
	HighlightActive bool
	Validate        func(OptionID) error
}

// MultiSelectRequest describes zero or more selections from stable options.
type MultiSelectRequest struct {
	ID          string
	Label       string
	Description string
	Options     []Option
	Defaults    []OptionID
	Minimum     int
	Maximum     int
	Validate    func([]OptionID) error
}

// Answers accumulates resolved values from a guided flow (Prompter.Guide),
// keyed by each Step's ID, so a later step's request can react to an
// earlier answer -- for example narrowing a template list to the chosen
// project kind. A step that was skipped has no entry.
type Answers map[string]any

// StepKind selects which prompt primitive a Step renders as.
type StepKind uint8

const (
	StepInput StepKind = iota
	StepSelect
	StepConfirm
)

// Step describes one question in a multi-screen guided flow that supports
// moving backward across the whole flow, not only within one prompt. The
// Input/Select/Confirm builder matching Kind is called fresh each time the
// step is (re)shown, so its request can read earlier Answers.
type Step struct {
	ID      string
	Kind    StepKind
	Input   func(Answers) InputRequest
	Select  func(Answers) SelectRequest
	Confirm func(Answers) ConfirmRequest
	// Skip reports whether to bypass this step entirely, in both
	// directions. Evaluated once per Guide call: every skip condition
	// Aruo has today is a static flag decided before the guide starts,
	// not something that changes based on answers gathered mid-flow.
	Skip func() bool
}

// TaskEventKind identifies a transition in a long-running task.
type TaskEventKind string

const (
	TaskStarted   TaskEventKind = "started"
	TaskAdvanced  TaskEventKind = "advanced"
	TaskMessage   TaskEventKind = "message"
	TaskRetrying  TaskEventKind = "retrying"
	TaskCompleted TaskEventKind = "completed"
	TaskFailed    TaskEventKind = "failed"
	TaskCancelled TaskEventKind = "cancelled"
)

// TaskEvent is a renderer-neutral progress transition.
type TaskEvent struct {
	TaskID   string
	ParentID string
	Kind     TaskEventKind
	Label    string
	Current  int64
	Total    int64
	Unit     string
	At       time.Time
	Detail   map[string]string
}
