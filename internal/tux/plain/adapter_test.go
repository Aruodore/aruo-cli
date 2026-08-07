package plain_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aruodore/aruo/internal/tux"
	"github.com/aruodore/aruo/internal/tux/plain"
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
