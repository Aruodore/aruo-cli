package charm

import (
	"context"
	"errors"
	"fmt"
	"io"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/aruodore/aruo/internal/tux"
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

// Input resolves one validated line of text.
func (p *Prompter) Input(ctx context.Context, request tux.InputRequest) (string, error) {
	if err := p.available(); err != nil {
		return "", err
	}
	value := ""
	if request.Default != nil {
		value = *request.Default
	}
	field := huh.NewInput().
		Key(request.ID).
		Title(request.Label).
		Description(request.Description).
		Placeholder(request.Placeholder).
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
	if err := ctx.Err(); err != nil {
		return errors.Join(tux.ErrCancelled, err)
	}
	width := p.capabilities.Width
	if width <= 0 {
		width = 80
	}
	form := huh.NewForm(huh.NewGroup(field)).
		WithInput(p.in).
		WithOutput(p.out).
		WithWidth(width).
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
	colorEligible := p.policy.Color != tux.FeatureNever && p.capabilities.Color == tux.ColorTrueColor
	result := make([]huh.Option[tux.OptionID], 0, len(options))
	for _, option := range options {
		if option.Disabled {
			continue
		}
		label := option.Label
		if option.Description != "" {
			label += " — " + option.Description
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
