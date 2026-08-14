package charm

import (
	"context"
	"errors"
	"fmt"
	"io"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// Prompter adapts Huh forms to Aruo prompt requests.
type Prompter struct {
	in           io.Reader
	out          io.Writer
	capabilities tux.Capabilities
	policy       tux.Policy
}

// NewPrompter constructs a prompt adapter over explicit terminal streams.
func NewPrompter(in io.Reader, out io.Writer, capabilities tux.Capabilities, policy tux.Policy) *Prompter {
	return &Prompter{in: in, out: out, capabilities: capabilities, policy: policy}
}

// Input resolves one validated line of text. A default never occupies the
// editable buffer as pre-filled text the user has to delete; it only shows
// as placeholder hint text, and is substituted after the fact if the user
// submits the field genuinely empty, matching the plain adapter's contract.
func (p *Prompter) Input(ctx context.Context, request tux.InputRequest) (string, error) {
	if err := p.available(); err != nil {
		return "", err
	}
	value := ""
	placeholder := request.Placeholder
	if request.Default != nil && *request.Default != "" {
		placeholder = *request.Default
	}
	field := huh.NewInput().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Prompt(": ").
		Placeholder(placeholder).
		Suggestions(request.Suggestions).
		Value(&value).
		Validate(func(value string) error {
			if value == "" && !request.Optional {
				return fmt.Errorf("%s is required", request.Label)
			}
			if request.Validate != nil {
				return request.Validate(value)
			}
			return nil
		})
	if err := p.run(ctx, field); err != nil {
		return "", err
	}
	if value == "" && request.Default != nil {
		value = *request.Default
	}
	return value, nil
}

// Secret resolves one validated value without echo.
func (p *Prompter) Secret(ctx context.Context, request tux.SecretRequest) (string, error) {
	if err := p.available(); err != nil {
		return "", err
	}
	value := ""
	field := huh.NewInput().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Prompt(": ").
		EchoMode(huh.EchoModeNone).
		Value(&value).
		Validate(func(value string) error {
			if value == "" && !request.Optional {
				return fmt.Errorf("%s is required", request.Label)
			}
			if request.Validate != nil {
				return request.Validate(value)
			}
			return nil
		})
	if err := p.run(ctx, field); err != nil {
		return "", err
	}
	return value, nil
}

// Confirm resolves one yes-or-no decision.
func (p *Prompter) Confirm(ctx context.Context, request tux.ConfirmRequest) (bool, error) {
	if err := p.available(); err != nil {
		return false, err
	}
	value := request.Default
	field := huh.NewConfirm().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Affirmative("Yes").
		Negative("No").
		Value(&value)
	if err := p.run(ctx, field); err != nil {
		return false, err
	}
	return value, nil
}

// Select resolves one stable option identifier.
func (p *Prompter) Select(ctx context.Context, request tux.SelectRequest) (tux.OptionID, error) {
	if err := p.available(); err != nil {
		return "", err
	}
	options := p.huhOptions(request.Options, nil)
	if len(options) == 0 {
		return "", errors.New("selection has no enabled options")
	}
	value := options[0].Value
	if request.Default != nil {
		value = *request.Default
	}
	field := huh.NewSelect[tux.OptionID]().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Options(options...).
		Value(&value)
	if request.HighlightActive {
		field.WithTheme(newHighlightedHuhTheme(p.capabilities, p.policy))
	}
	if request.Validate != nil {
		field.Validate(request.Validate)
	}
	if err := p.run(ctx, field); err != nil {
		return "", err
	}
	return value, nil
}

// MultiSelect resolves stable option identifiers.
func (p *Prompter) MultiSelect(ctx context.Context, request tux.MultiSelectRequest) ([]tux.OptionID, error) {
	if err := p.available(); err != nil {
		return nil, err
	}
	defaults := make(map[tux.OptionID]struct{}, len(request.Defaults))
	for _, id := range request.Defaults {
		defaults[id] = struct{}{}
	}
	options := p.huhOptions(request.Options, defaults)
	if len(options) == 0 && request.Minimum > 0 {
		return nil, errors.New("selection has no enabled options")
	}
	value := append([]tux.OptionID(nil), request.Defaults...)
	field := huh.NewMultiSelect[tux.OptionID]().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Options(options...).
		Value(&value).
		Filterable(len(options) > 8).
		Validate(func(selected []tux.OptionID) error {
			if len(selected) < request.Minimum {
				return fmt.Errorf("select at least %d options", request.Minimum)
			}
			if request.Maximum > 0 && len(selected) > request.Maximum {
				return fmt.Errorf("select at most %d options", request.Maximum)
			}
			if request.Validate != nil {
				return request.Validate(selected)
			}
			return nil
		})
	if request.Maximum > 0 {
		field.Limit(request.Maximum)
	}
	if err := p.run(ctx, field); err != nil {
		return nil, err
	}
	return value, nil
}

func (p *Prompter) available() error {
	if p.policy.Mode != tux.ModeInteractive || p.policy.Input == tux.FeatureNever {
		return tux.ErrUnavailable
	}
	return nil
}

func (p *Prompter) run(ctx context.Context, field huh.Field) error {
	return p.runForm(ctx, huh.NewForm(huh.NewGroup(field)))
}

// runForm applies the shared terminal wiring and error mapping to an
// already-built form, whether it's a single-field form from run or the
// multi-group form Guide assembles.
func (p *Prompter) runForm(ctx context.Context, form *huh.Form) error {
	if err := ctx.Err(); err != nil {
		return errors.Join(tux.ErrCancelled, err)
	}
	width := p.capabilities.Width
	if width <= 0 {
		width = 80
	}
	form = form.
		WithInput(p.in).
		WithOutput(p.out).
		WithWidth(width).
		WithTheme(newHuhTheme(p.capabilities, p.policy)).
		WithAccessible(p.policy.Accessible)
	if err := form.RunWithContext(ctx); err != nil {
		switch {
		case errors.Is(err, huh.ErrUserAborted), ctx.Err() != nil:
			return errors.Join(tux.ErrCancelled, err, ctx.Err())
		default:
			return fmt.Errorf("run prompt: %w", err)
		}
	}
	return nil
}

func (p *Prompter) huhOptions(options []tux.Option, defaults map[tux.OptionID]struct{}) []huh.Option[tux.OptionID] {
	colorEligible := p.policy.Color != tux.FeatureNever && !p.policy.Accessible && p.capabilities.Color == tux.ColorTrueColor
	result := make([]huh.Option[tux.OptionID], 0, len(options))
	for _, option := range options {
		if option.Disabled {
			continue
		}
		label := option.Label
		if option.Description != "" {
			label += "  " + option.Description
		}
		if colorEligible && option.Color != "" {
			label = lipgloss.NewStyle().Foreground(lipgloss.Color(option.Color)).Render(label)
		}
		converted := huh.NewOption(label, option.ID)
		if _, selected := defaults[option.ID]; selected {
			converted = converted.Selected(true)
		}
		result = append(result, converted)
	}
	return result
}
