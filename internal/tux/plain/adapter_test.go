package plain_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/aruodore/aruo-cli/internal/tux/plain"
)

func TestInputRepeatsAfterValidationFailure(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader("bad\ngood\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	value, err := adapter.Input(context.Background(), tux.InputRequest{
		Label: "Name",
		Validate: func(value string) error {
			if value != "good" {
				return errors.New("use good")
			}
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != "good" {
		t.Fatalf("Input() = %q, want good", value)
	}
	if got := diagnostic.String(); !strings.Contains(got, "Invalid value: use good\nName: ") {
		t.Errorf("diagnostic = %q", got)
	}
}

func TestInputEmptyDefaultOmitsBracketHint(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	empty := ""
	adapter := plain.New(strings.NewReader("\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	value, err := adapter.Input(context.Background(), tux.InputRequest{
		Label:    "Author or organization",
		Optional: true,
		Default:  &empty,
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != "" {
		t.Fatalf("Input() = %q, want empty when the default itself is empty", value)
	}
	if got := diagnostic.String(); strings.Contains(got, "[]") {
		t.Errorf("diagnostic = %q, want no empty-bracket hint for a blank default", got)
	}
}

func TestSelectUsesStableID(t *testing.T) {
	t.Parallel()

	adapter := plain.New(strings.NewReader("2\n"), &bytes.Buffer{}, &bytes.Buffer{}, true, nil)
	selected, err := adapter.Select(context.Background(), tux.SelectRequest{
		Label: "Template",
		Options: []tux.Option{
			{ID: "go-library", Label: "Go library"},
			{ID: "go-cli", Label: "Go command-line application"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if selected != "go-cli" {
		t.Fatalf("Select() = %q, want go-cli", selected)
	}
}

func TestConfirmRejectsAmbiguousAnswer(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader("perhaps\nyes\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	confirmed, err := adapter.Confirm(context.Background(), tux.ConfirmRequest{Label: "Continue?"})
	if err != nil {
		t.Fatal(err)
	}
	if !confirmed {
		t.Fatal("Confirm() = false, want true")
	}
	if !strings.Contains(diagnostic.String(), "Invalid value: enter yes or no") {
		t.Errorf("diagnostic = %q", diagnostic.String())
	}
}

func TestInputDisabled(t *testing.T) {
	t.Parallel()

	adapter := plain.New(strings.NewReader("ignored\n"), &bytes.Buffer{}, &bytes.Buffer{}, false, nil)
	_, err := adapter.Input(context.Background(), tux.InputRequest{Label: "Name"})
	if !errors.Is(err, tux.ErrUnavailable) {
		t.Fatalf("Input() error = %v, want ErrUnavailable", err)
	}
}

func TestPresenterKeepsDiagnosticsOffResultStream(t *testing.T) {
	t.Parallel()

	var result bytes.Buffer
	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader(""), &result, &diagnostic, false, nil)
	if err := adapter.Message(context.Background(), tux.Message{Kind: tux.MessageSuccess, Text: "created project"}); err != nil {
		t.Fatal(err)
	}
	if err := adapter.Diagnostic(context.Background(), tux.Diagnostic{Summary: "could not create project", Suggestion: "choose another path"}); err != nil {
		t.Fatal(err)
	}
	if result.String() != "SUCCESS: created project\n" {
		t.Errorf("result = %q", result.String())
	}
	if !strings.Contains(diagnostic.String(), "Error: could not create project") {
		t.Errorf("diagnostic = %q", diagnostic.String())
	}
}

func nameStep() tux.Step {
	return tux.Step{
		ID:   "name",
		Kind: tux.StepInput,
		Input: func(tux.Answers) tux.InputRequest {
			return tux.InputRequest{Label: "Project name"}
		},
	}
}

func kindStep(skip func() bool) tux.Step {
	return tux.Step{
		ID:   "kind",
		Kind: tux.StepSelect,
		Select: func(tux.Answers) tux.SelectRequest {
			return tux.SelectRequest{
				Label: "What are you building?",
				Options: []tux.Option{
					{ID: "app", Label: "Application"},
					{ID: "library", Label: "Library"},
				},
			}
		},
		Skip: skip,
	}
}

func confirmStep() tux.Step {
	return tux.Step{
		ID:   "confirm",
		Kind: tux.StepConfirm,
		Confirm: func(tux.Answers) tux.ConfirmRequest {
			return tux.ConfirmRequest{Label: "Create this project?"}
		},
	}
}

func TestGuideLinearHappyPath(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader("my-app\n2\nyes\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	answers, err := adapter.Guide(context.Background(), []tux.Step{nameStep(), kindStep(nil), confirmStep()})
	if err != nil {
		t.Fatalf("Guide() error = %v", err)
	}
	if answers["name"] != "my-app" {
		t.Errorf("answers[name] = %v, want my-app", answers["name"])
	}
	if answers["kind"] != tux.OptionID("library") {
		t.Errorf("answers[kind] = %v, want library", answers["kind"])
	}
	if answers["confirm"] != true {
		t.Errorf("answers[confirm] = %v, want true", answers["confirm"])
	}
	if !strings.Contains(diagnostic.String(), "Type back at any prompt to return to the previous question.") {
		t.Errorf("diagnostic = %q, want the back-navigation tip", diagnostic.String())
	}
}

func TestGuideBackReturnsToPreviousStepWithPriorAnswerAsDefault(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	// name=my-app, then on the kind step type "back", re-see name with
	// "my-app" as the default (bare Enter keeps it), then answer kind=1
	// (app) and confirm.
	adapter := plain.New(strings.NewReader("my-app\nback\n\n1\nyes\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	answers, err := adapter.Guide(context.Background(), []tux.Step{nameStep(), kindStep(nil), confirmStep()})
	if err != nil {
		t.Fatalf("Guide() error = %v", err)
	}
	if answers["name"] != "my-app" {
		t.Fatalf("answers[name] = %v, want my-app to survive a bare Enter on revisit", answers["name"])
	}
	if answers["kind"] != tux.OptionID("app") {
		t.Fatalf("answers[kind] = %v, want app", answers["kind"])
	}
	if !strings.Contains(diagnostic.String(), "Project name [my-app]:") {
		t.Errorf("diagnostic = %q, want the revisit prompt to show the prior answer as its default", diagnostic.String())
	}
}

func TestGuideBackAtFirstStepPrintsNoticeAndReprompts(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader("back\nmy-app\n2\nyes\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	answers, err := adapter.Guide(context.Background(), []tux.Step{nameStep(), kindStep(nil), confirmStep()})
	if err != nil {
		t.Fatalf("Guide() error = %v", err)
	}
	if answers["name"] != "my-app" {
		t.Fatalf("answers[name] = %v, want my-app", answers["name"])
	}
	if !strings.Contains(diagnostic.String(), "Already at the first question.") {
		t.Errorf("diagnostic = %q, want a notice instead of erroring", diagnostic.String())
	}
}

func TestGuideSkipsStepInBothDirections(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	// kind is always skipped; going "back" from confirm must land on name,
	// not the hidden kind step.
	adapter := plain.New(strings.NewReader("my-app\nback\nrenamed\nyes\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	answers, err := adapter.Guide(context.Background(), []tux.Step{
		nameStep(),
		kindStep(func() bool { return true }),
		confirmStep(),
	})
	if err != nil {
		t.Fatalf("Guide() error = %v", err)
	}
	if _, answered := answers["kind"]; answered {
		t.Errorf("answers[kind] = %v, want the skipped step to have no entry", answers["kind"])
	}
	if answers["name"] != "renamed" {
		t.Fatalf("answers[name] = %v, want renamed after going back and re-answering", answers["name"])
	}
	if strings.Contains(diagnostic.String(), "What are you building?") {
		t.Errorf("diagnostic = %q, want the skipped kind step never rendered", diagnostic.String())
	}
}

func TestGuideOmitsTipWithOnlyOneVisibleStep(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader("my-app\n"), &bytes.Buffer{}, &diagnostic, true, nil)
	_, err := adapter.Guide(context.Background(), []tux.Step{
		nameStep(),
		kindStep(func() bool { return true }),
	})
	if err != nil {
		t.Fatalf("Guide() error = %v", err)
	}
	if strings.Contains(diagnostic.String(), "Type back") {
		t.Errorf("diagnostic = %q, want no back-navigation tip when only one step is visible", diagnostic.String())
	}
}

func TestGuideDisabledReturnsUnavailable(t *testing.T) {
	t.Parallel()

	adapter := plain.New(strings.NewReader("ignored\n"), &bytes.Buffer{}, &bytes.Buffer{}, false, nil)
	_, err := adapter.Guide(context.Background(), []tux.Step{nameStep()})
	if !errors.Is(err, tux.ErrUnavailable) {
		t.Fatalf("Guide() error = %v, want ErrUnavailable", err)
	}
}

func TestProgressWritesOnlyDurableTransitions(t *testing.T) {
	t.Parallel()

	var diagnostic bytes.Buffer
	adapter := plain.New(strings.NewReader(""), &bytes.Buffer{}, &diagnostic, false, nil)
	for _, event := range []tux.TaskEvent{
		{Kind: tux.TaskStarted, Label: "audit repository"},
		{Kind: tux.TaskAdvanced, Label: "audit repository", Current: 1, Total: 2},
		{Kind: tux.TaskCompleted, Label: "audit repository"},
	} {
		if err := adapter.Emit(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	}
	if diagnostic.String() != "started: audit repository\ncompleted: audit repository\n" {
		t.Errorf("diagnostic = %q", diagnostic.String())
	}
}
