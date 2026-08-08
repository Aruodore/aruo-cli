package charm

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// guideField tracks one step's bound Go value alongside its ID, so a later
// step's reactive Title/Description/Options can read an earlier answer,
// and so Guide can assemble the final Answers once the form completes.
type guideField struct {
	pointer any // *string, *tux.OptionID, or *bool, matching the step's Kind
}

func (f guideField) value() any {
	switch pointer := f.pointer.(type) {
	case *string:
		return *pointer
	case *tux.OptionID:
		return *pointer
	case *bool:
		return *pointer
	default:
		return nil
	}
}

// Guide runs a multi-screen flow as one continuous Huh form, giving the
// user Huh's native shift+tab back-navigation across every screen instead
// of only within one (today's Input/Select/Confirm each build their own
// isolated single-field form, which is structurally always "the first
// field" and can never go back). A step whose Skip reports true becomes a
// hidden group, skipped in both directions by Huh itself.
func (p *Prompter) Guide(ctx context.Context, steps []tux.Step) (tux.Answers, error) {
	if err := p.available(); err != nil {
		return nil, err
	}
	if countVisibleGuideSteps(steps) == 0 {
		return tux.Answers{}, nil
	}
	form, fields := p.buildGuideForm(steps)
	if err := p.runForm(ctx, form); err != nil {
		return nil, err
	}
	return collectGuideAnswers(steps, fields), nil
}

// buildGuideForm assembles the multi-group form and the bound pointers
// backing it. Split out from Guide so tests can drive the form directly
// via Form.Update, without going through RunWithContext's real terminal
// program.
func (p *Prompter) buildGuideForm(steps []tux.Step) (*huh.Form, []guideField) {
	fields := make([]guideField, len(steps))
	groups := make([]*huh.Group, len(steps))
	for index, step := range steps {
		field := p.buildGuideField(step, steps, fields, index)
		group := huh.NewGroup(field)
		if step.Skip != nil {
			group = group.WithHide(step.Skip())
		}
		groups[index] = group
	}
	return huh.NewForm(groups...), fields
}

// collectGuideAnswers reads the final bound values into an Answers map
// once the form has completed (or, in tests, once the caller is satisfied
// the relevant fields have settled).
func collectGuideAnswers(steps []tux.Step, fields []guideField) tux.Answers {
	answers := make(tux.Answers, len(steps))
	for index, step := range steps {
		if step.Skip != nil && step.Skip() {
			continue
		}
		value := fields[index].value()
		if step.Kind == tux.StepInput {
			// A default never occupies the editable buffer as pre-filled
			// text (matches Input's own contract); apply it here, once,
			// only if the field was left genuinely empty.
			if text, ok := value.(string); ok && text == "" {
				if request := step.Input(answers); request.Default != nil && *request.Default != "" {
					value = *request.Default
				}
			}
		}
		answers[step.ID] = value
	}
	return answers
}

// buildGuideField constructs the Huh field for one step and records its
// bound pointer into fields[index]. Every *Func is wired to a snapshot of
// the preceding, non-hidden steps' current answers, so a step can react to
// an earlier one (for example the template list narrowing to the chosen
// kind) without Guide needing to know which steps actually depend on which.
func (p *Prompter) buildGuideField(step tux.Step, steps []tux.Step, fields []guideField, index int) huh.Field {
	snapshot := func() tux.Answers {
		answers := make(tux.Answers, index)
		for i := 0; i < index; i++ {
			if steps[i].Skip != nil && steps[i].Skip() {
				continue
			}
			answers[steps[i].ID] = fields[i].value()
		}
		return answers
	}
	// huh hashes bindings to decide whether a *Func needs recomputing, and
	// an Eval's bindingsHash starts at its zero value (0). Hashing a truly
	// empty []any{} -- which step 0 would otherwise pass, having no
	// preceding steps -- also happens to produce exactly 0, so its *Func
	// would look "already up to date" and never run even once. Seed with
	// the step's own ID so bindings is never empty.
	bindings := make([]any, 0, index+1)
	bindings = append(bindings, step.ID)
	for i := 0; i < index; i++ {
		bindings = append(bindings, fields[i].pointer)
	}

	switch step.Kind {
	case tux.StepSelect:
		return p.buildGuideSelect(step, snapshot, bindings, &fields[index])
	case tux.StepConfirm:
		return buildGuideConfirm(step, snapshot, bindings, &fields[index])
	default:
		return buildGuideInput(step, snapshot, bindings, &fields[index])
	}
}

func buildGuideInput(step tux.Step, snapshot func() tux.Answers, bindings []any, out *guideField) huh.Field {
	value := ""
	out.pointer = &value
	request := func() tux.InputRequest { return step.Input(snapshot()) }
	return huh.NewInput().
		Key(step.ID).
		TitleFunc(func() string { return request().Label }, bindings).
		DescriptionFunc(func() string { return request().Description }, bindings).
		PlaceholderFunc(func() string {
			req := request()
			if req.Default != nil && *req.Default != "" {
				return *req.Default
			}
			return req.Placeholder
		}, bindings).
		SuggestionsFunc(func() []string { return request().Suggestions }, bindings).
		Value(&value).
		Validate(func(v string) error {
			req := request()
			if v == "" && !req.Optional {
				return fmt.Errorf("%s is required", req.Label)
			}
			if req.Validate != nil {
				return req.Validate(v)
			}
			return nil
		})
}

func (p *Prompter) buildGuideSelect(step tux.Step, snapshot func() tux.Answers, bindings []any, out *guideField) huh.Field {
	initial := step.Select(snapshot())
	value := firstOptionID(initial.Options)
	if initial.Default != nil {
		value = *initial.Default
	}
	out.pointer = &value
	request := func() tux.SelectRequest { return step.Select(snapshot()) }
	return huh.NewSelect[tux.OptionID]().
		Key(step.ID).
		TitleFunc(func() string { return request().Label }, bindings).
		DescriptionFunc(func() string { return request().Description }, bindings).
		OptionsFunc(func() []huh.Option[tux.OptionID] { return p.huhOptions(request().Options, nil) }, bindings).
		Value(&value).
		Validate(func(v tux.OptionID) error {
			if req := request(); req.Validate != nil {
				return req.Validate(v)
			}
			return nil
		})
}

func buildGuideConfirm(step tux.Step, snapshot func() tux.Answers, bindings []any, out *guideField) huh.Field {
	value := step.Confirm(snapshot()).Default
	out.pointer = &value
	request := func() tux.ConfirmRequest { return step.Confirm(snapshot()) }
	return huh.NewConfirm().
		Key(step.ID).
		TitleFunc(func() string { return request().Label }, bindings).
		DescriptionFunc(func() string { return request().Description }, bindings).
		Affirmative("Yes").
		Negative("No").
		Value(&value)
}

func firstOptionID(options []tux.Option) tux.OptionID {
	for _, option := range options {
		if !option.Disabled {
			return option.ID
		}
	}
	return ""
}

func countVisibleGuideSteps(steps []tux.Step) int {
	count := 0
	for _, step := range steps {
		if step.Skip == nil || !step.Skip() {
			count++
		}
	}
	return count
}
