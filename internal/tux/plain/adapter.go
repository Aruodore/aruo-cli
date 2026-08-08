// Package plain provides the line-oriented terminal reference adapter.
package plain

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/aruodore/aruo-cli/internal/tux"
)

// SecretReader reads a value without exposing it through the ordinary line reader.
type SecretReader func(context.Context, string) (string, error)

// Adapter implements accessible prompts and durable presentation.
type Adapter struct {
	scanner      *bufio.Scanner
	result       io.Writer
	diagnostic   io.Writer
	inputEnabled bool
	readSecret   SecretReader
	mu           sync.Mutex
}

// New constructs a plain adapter over explicit streams.
func New(in io.Reader, result, diagnostic io.Writer, inputEnabled bool, readSecret SecretReader) *Adapter {
	return &Adapter{
		scanner:      bufio.NewScanner(in),
		result:       result,
		diagnostic:   diagnostic,
		inputEnabled: inputEnabled,
		readSecret:   readSecret,
	}
}

// Input resolves one line of validated text.
func (a *Adapter) Input(ctx context.Context, request tux.InputRequest) (string, error) {
	if !a.inputEnabled {
		return "", tux.ErrUnavailable
	}
	for {
		if err := ctx.Err(); err != nil {
			return "", errors.Join(tux.ErrCancelled, err)
		}
		if err := a.writePrompt(request.Description, request.Example, inputLabel(request)); err != nil {
			return "", err
		}
		answer, err := a.scan()
		if err != nil {
			return "", err
		}
		answer = strings.TrimSpace(answer)
		if answer == "" && request.Default != nil {
			answer = *request.Default
		}
		if answer == "" && !request.Optional {
			if err := a.validationError(fmt.Errorf("%s is required", request.Label)); err != nil {
				return "", err
			}
			continue
		}
		if request.Validate != nil {
			if err := request.Validate(answer); err != nil {
				if writeErr := a.validationError(err); writeErr != nil {
					return "", writeErr
				}
				continue
			}
		}
		return answer, nil
	}
}

// Secret resolves one no-echo value through an injected terminal primitive.
func (a *Adapter) Secret(ctx context.Context, request tux.SecretRequest) (string, error) {
	if !a.inputEnabled || a.readSecret == nil {
		return "", tux.ErrUnavailable
	}
	if err := ctx.Err(); err != nil {
		return "", errors.Join(tux.ErrCancelled, err)
	}
	if request.Description != "" {
		if _, err := fmt.Fprintln(a.diagnostic, request.Description); err != nil {
			return "", fmt.Errorf("write secret description: %w", err)
		}
	}
	value, err := a.readSecret(ctx, request.Label+": ")
	if err != nil {
		return "", err
	}
	if value == "" && !request.Optional {
		return "", fmt.Errorf("%s is required", request.Label)
	}
	if request.Validate != nil {
		if err := request.Validate(value); err != nil {
			return "", err
		}
	}
	return value, nil
}

// Confirm resolves a yes-or-no decision.
func (a *Adapter) Confirm(ctx context.Context, request tux.ConfirmRequest) (bool, error) {
	defaultLabel := "y/N"
	if request.Default {
		defaultLabel = "Y/n"
	}
	for {
		answer, err := a.Input(ctx, tux.InputRequest{
			ID:          request.ID,
			Label:       request.Label + " [" + defaultLabel + "]",
			Description: request.Description,
			Optional:    true,
		})
		if err != nil {
			return false, err
		}
		switch strings.ToLower(answer) {
		case "":
			return request.Default, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			if err := a.validationError(errors.New("enter yes or no")); err != nil {
				return false, err
			}
		}
	}
}

// Select resolves one numbered option.
func (a *Adapter) Select(ctx context.Context, request tux.SelectRequest) (tux.OptionID, error) {
	options := enabledOptions(request.Options)
	if len(options) == 0 {
		return "", errors.New("selection has no enabled options")
	}
	for {
		if err := a.writeOptions(request.Description, request.Label, options, request.Default); err != nil {
			return "", err
		}
		answer, err := a.readChoice(ctx)
		if err != nil {
			return "", err
		}
		if answer == "" && request.Default != nil {
			answer = string(*request.Default)
		}
		selected, ok := resolveOption(answer, options)
		if !ok {
			if err := a.validationError(errors.New("enter an option number or identifier")); err != nil {
				return "", err
			}
			continue
		}
		if request.Validate != nil {
			if err := request.Validate(selected); err != nil {
				if writeErr := a.validationError(err); writeErr != nil {
					return "", writeErr
				}
				continue
			}
		}
		return selected, nil
	}
}

// MultiSelect resolves comma-separated numbered options.
func (a *Adapter) MultiSelect(ctx context.Context, request tux.MultiSelectRequest) ([]tux.OptionID, error) {
	options := enabledOptions(request.Options)
	if len(options) == 0 && request.Minimum > 0 {
		return nil, errors.New("selection has no enabled options")
	}
	for {
		if err := a.writeOptions(request.Description, request.Label+" (comma-separated)", options, nil); err != nil {
			return nil, err
		}
		answer, err := a.readChoice(ctx)
		if err != nil {
			return nil, err
		}
		var selected []tux.OptionID
		if strings.TrimSpace(answer) == "" {
			selected = append(selected, request.Defaults...)
		} else {
			seen := make(map[tux.OptionID]struct{})
			valid := true
			for _, value := range strings.Split(answer, ",") {
				option, ok := resolveOption(strings.TrimSpace(value), options)
				if !ok {
					valid = false
					break
				}
				if _, exists := seen[option]; !exists {
					seen[option] = struct{}{}
					selected = append(selected, option)
				}
			}
			if !valid {
				if err := a.validationError(errors.New("enter comma-separated option numbers or identifiers")); err != nil {
					return nil, err
				}
				continue
			}
		}
		if len(selected) < request.Minimum || request.Maximum > 0 && len(selected) > request.Maximum {
			if err := a.validationError(fmt.Errorf("select between %d and %d options", request.Minimum, request.Maximum)); err != nil {
				return nil, err
			}
			continue
		}
		if request.Validate != nil {
			if err := request.Validate(selected); err != nil {
				if writeErr := a.validationError(err); writeErr != nil {
					return nil, writeErr
				}
				continue
			}
		}
		return selected, nil
	}
}

// Guide runs a multi-screen flow as one session. Typing "back" (case
// insensitive) at any prompt returns to the previous visible step instead
// of being taken as that step's value. A step that genuinely needs the
// literal text "back" can still be supplied via its flag, bypassing the
// prompt entirely.
func (a *Adapter) Guide(ctx context.Context, steps []tux.Step) (tux.Answers, error) {
	if !a.inputEnabled {
		return nil, tux.ErrUnavailable
	}
	if err := ctx.Err(); err != nil {
		return nil, errors.Join(tux.ErrCancelled, err)
	}
	if countVisibleSteps(steps) > 1 {
		if _, err := fmt.Fprintln(a.diagnostic, "Type back at any prompt to return to the previous question."); err != nil {
			return nil, err
		}
	}
	answers := make(tux.Answers, len(steps))
	index := 0
	for index < len(steps) {
		step := steps[index]
		if step.Skip != nil && step.Skip() {
			index++
			continue
		}
		var (
			value any
			back  bool
			err   error
		)
		switch step.Kind {
		case tux.StepSelect:
			value, back, err = a.guideSelectStep(ctx, step, answers)
		case tux.StepConfirm:
			value, back, err = a.guideConfirmStep(ctx, step, answers)
		default:
			value, back, err = a.guideInputStep(ctx, step, answers)
		}
		if err != nil {
			return nil, err
		}
		if back {
			prev := prevVisibleStep(steps, index)
			if prev < 0 {
				if _, err := fmt.Fprintln(a.diagnostic, "Already at the first question."); err != nil {
					return nil, err
				}
				continue
			}
			index = prev
			continue
		}
		answers[step.ID] = value
		index++
	}
	return answers, nil
}

// guideInputStep prompts for step's Input request, recognizing "back" as a
// navigation command rather than a literal value. A previous answer for
// this step (from before the user went back) becomes the default, so a
// bare Enter on revisit keeps it.
func (a *Adapter) guideInputStep(ctx context.Context, step tux.Step, answers tux.Answers) (string, bool, error) {
	request := step.Input(answers)
	if previous, ok := answers[step.ID].(string); ok {
		request.Default = &previous
	}
	for {
		if err := ctx.Err(); err != nil {
			return "", false, errors.Join(tux.ErrCancelled, err)
		}
		if err := a.writePrompt(request.Description, request.Example, inputLabel(request)); err != nil {
			return "", false, err
		}
		raw, err := a.scan()
		if err != nil {
			return "", false, err
		}
		trimmed := strings.TrimSpace(raw)
		if isBackCommand(trimmed) {
			return "", true, nil
		}
		answer := trimmed
		if answer == "" && request.Default != nil {
			answer = *request.Default
		}
		if answer == "" && !request.Optional {
			if err := a.validationError(fmt.Errorf("%s is required", request.Label)); err != nil {
				return "", false, err
			}
			continue
		}
		if request.Validate != nil {
			if err := request.Validate(answer); err != nil {
				if writeErr := a.validationError(err); writeErr != nil {
					return "", false, writeErr
				}
				continue
			}
		}
		return answer, false, nil
	}
}

// guideSelectStep mirrors guideInputStep for a Select step.
func (a *Adapter) guideSelectStep(ctx context.Context, step tux.Step, answers tux.Answers) (tux.OptionID, bool, error) {
	request := step.Select(answers)
	options := enabledOptions(request.Options)
	if len(options) == 0 {
		return "", false, errors.New("selection has no enabled options")
	}
	defaultID := request.Default
	if previous, ok := answers[step.ID].(tux.OptionID); ok {
		defaultID = &previous
	}
	for {
		if err := a.writeOptions(request.Description, request.Label, options, defaultID); err != nil {
			return "", false, err
		}
		raw, err := a.readChoice(ctx)
		if err != nil {
			return "", false, err
		}
		trimmed := strings.TrimSpace(raw)
		if isBackCommand(trimmed) {
			return "", true, nil
		}
		answer := trimmed
		if answer == "" && defaultID != nil {
			answer = string(*defaultID)
		}
		selected, ok := resolveOption(answer, options)
		if !ok {
			if err := a.validationError(errors.New("enter an option number or identifier")); err != nil {
				return "", false, err
			}
			continue
		}
		if request.Validate != nil {
			if err := request.Validate(selected); err != nil {
				if writeErr := a.validationError(err); writeErr != nil {
					return "", false, writeErr
				}
				continue
			}
		}
		return selected, false, nil
	}
}

// guideConfirmStep mirrors Confirm's own y/n handling (including its
// delegation to the Input primitives), adding back-detection and
// per-visit default injection the same way the other two guide steps do.
func (a *Adapter) guideConfirmStep(ctx context.Context, step tux.Step, answers tux.Answers) (bool, bool, error) {
	request := step.Confirm(answers)
	defaultValue := request.Default
	if previous, ok := answers[step.ID].(bool); ok {
		defaultValue = previous
	}
	defaultLabel := "y/N"
	if defaultValue {
		defaultLabel = "Y/n"
	}
	inputStep := tux.Step{
		ID: step.ID,
		Input: func(tux.Answers) tux.InputRequest {
			return tux.InputRequest{
				Label:       request.Label + " [" + defaultLabel + "]",
				Description: request.Description,
				Optional:    true,
			}
		},
	}
	for {
		answer, back, err := a.guideInputStep(ctx, inputStep, tux.Answers{})
		if err != nil {
			return false, false, err
		}
		if back {
			return false, true, nil
		}
		switch strings.ToLower(strings.TrimSpace(answer)) {
		case "":
			return defaultValue, false, nil
		case "y", "yes":
			return true, false, nil
		case "n", "no":
			return false, false, nil
		default:
			if err := a.validationError(errors.New("enter yes or no")); err != nil {
				return false, false, err
			}
		}
	}
}

// isBackCommand reports whether a scanned line is the reserved navigation
// keyword rather than a literal value.
func isBackCommand(value string) bool {
	return strings.EqualFold(value, "back")
}

func prevVisibleStep(steps []tux.Step, from int) int {
	for i := from - 1; i >= 0; i-- {
		if steps[i].Skip == nil || !steps[i].Skip() {
			return i
		}
	}
	return -1
}

func countVisibleSteps(steps []tux.Step) int {
	count := 0
	for _, step := range steps {
		if step.Skip == nil || !step.Skip() {
			count++
		}
	}
	return count
}

// Message writes one durable semantic message.
func (a *Adapter) Message(ctx context.Context, message tux.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	label := strings.ToUpper(string(message.Kind))
	if message.Kind == tux.MessageInfo || message.Kind == "" {
		_, err := fmt.Fprintln(a.result, message.Text)
		return err
	}
	_, err := fmt.Fprintf(a.result, "%s: %s\n", label, message.Text)
	return err
}

// Table writes an adaptive plain-text table.
func (a *Adapter) Table(ctx context.Context, table tux.Table) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, row := range table.Rows {
		for index, column := range table.Columns {
			if index >= len(row) || column.Sensitive {
				continue
			}
			if _, err := fmt.Fprintf(a.result, "%s: %s\n", column.Heading, row[index]); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(a.result); err != nil {
			return err
		}
	}
	return nil
}

// Diagnostic writes one actionable problem to the diagnostic stream.
func (a *Adapter) Diagnostic(ctx context.Context, diagnostic tux.Diagnostic) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(a.diagnostic, "Error: %s\n", diagnostic.Summary); err != nil {
		return err
	}
	for _, value := range []string{diagnostic.Detail, diagnostic.Effect, diagnostic.Suggestion} {
		if value == "" {
			continue
		}
		if _, err := fmt.Fprintf(a.diagnostic, "\n%s\n", value); err != nil {
			return err
		}
	}
	if diagnostic.Code != "" {
		_, err := fmt.Fprintf(a.diagnostic, "\nDebug reference: %s\n", diagnostic.Code)
		return err
	}
	return nil
}

// Emit writes sparse durable progress transitions.
func (a *Adapter) Emit(ctx context.Context, event tux.TaskEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	switch event.Kind {
	case tux.TaskStarted, tux.TaskRetrying, tux.TaskCompleted, tux.TaskFailed, tux.TaskCancelled:
		_, err := fmt.Fprintf(a.diagnostic, "%s: %s\n", event.Kind, event.Label)
		return err
	default:
		return nil
	}
}

func (a *Adapter) scan() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.scanner.Scan() {
		if err := a.scanner.Err(); err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}
		return "", tux.ErrEndOfInput
	}
	return a.scanner.Text(), nil
}

func (a *Adapter) readChoice(ctx context.Context) (string, error) {
	if !a.inputEnabled {
		return "", tux.ErrUnavailable
	}
	if err := ctx.Err(); err != nil {
		return "", errors.Join(tux.ErrCancelled, err)
	}
	if _, err := fmt.Fprint(a.diagnostic, "Choice: "); err != nil {
		return "", err
	}
	return a.scan()
}

func (a *Adapter) writePrompt(description, example, label string) error {
	if description != "" {
		if _, err := fmt.Fprintln(a.diagnostic, description); err != nil {
			return err
		}
	}
	if example != "" {
		if _, err := fmt.Fprintf(a.diagnostic, "Example: %s\n", example); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(a.diagnostic, label+": ")
	return err
}

func (a *Adapter) writeOptions(description, label string, options []tux.Option, defaultID *tux.OptionID) error {
	if description != "" {
		if _, err := fmt.Fprintln(a.diagnostic, description); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(a.diagnostic, label); err != nil {
		return err
	}
	for index, option := range options {
		defaultMarker := ""
		if defaultID != nil && option.ID == *defaultID {
			defaultMarker = " (default)"
		}
		if _, err := fmt.Fprintf(a.diagnostic, "  %d. %s%s\n", index+1, option.Label, defaultMarker); err != nil {
			return err
		}
		if option.Description != "" {
			if _, err := fmt.Fprintf(a.diagnostic, "     %s\n", option.Description); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Adapter) validationError(err error) error {
	_, writeErr := fmt.Fprintf(a.diagnostic, "Invalid value: %s\n", err)
	return writeErr
}

func inputLabel(request tux.InputRequest) string {
	if request.Default == nil || *request.Default == "" {
		return request.Label
	}
	return fmt.Sprintf("%s [%s]", request.Label, *request.Default)
}

func enabledOptions(options []tux.Option) []tux.Option {
	result := make([]tux.Option, 0, len(options))
	for _, option := range options {
		if !option.Disabled {
			result = append(result, option)
		}
	}
	return result
}

func resolveOption(value string, options []tux.Option) (tux.OptionID, bool) {
	if index, err := strconv.Atoi(value); err == nil && index > 0 && index <= len(options) {
		return options[index-1].ID, true
	}
	for _, option := range options {
		if string(option.ID) == value {
			return option.ID, true
		}
	}
	return "", false
}
