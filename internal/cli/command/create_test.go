package command

import (
	"context"
	"errors"
	"testing"

	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/tux"
)

func TestParseVariables(t *testing.T) {
	t.Parallel()

	got, err := parseVariables([]string{"name=demo", "private=true", "count=false"})
	if err != nil {
		t.Fatalf("parseVariables() error = %v", err)
	}
	want := map[string]any{"name": "demo", "private": true, "count": false}
	if len(got) != len(want) || got["name"] != want["name"] || got["private"] != want["private"] || got["count"] != want["count"] {
		t.Fatalf("parseVariables() = %#v, want %#v", got, want)
	}
}

func TestParseVariablesRejectsMissingSeparator(t *testing.T) {
	t.Parallel()

	if _, err := parseVariables([]string{"nokeyvalue"}); err == nil {
		t.Fatal("parseVariables() error = nil, want error for missing '='")
	}
}

func TestParseVariablesRejectsEmptyKey(t *testing.T) {
	t.Parallel()

	if _, err := parseVariables([]string{"=value"}); err == nil {
		t.Fatal("parseVariables() error = nil, want error for empty key")
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	if !contains([]string{"MIT", "Apache-2.0"}, "Apache-2.0") {
		t.Error("contains() = false, want true for present value")
	}
	if contains([]string{"MIT"}, "Apache-2.0") {
		t.Error("contains() = true, want false for absent value")
	}
	if contains(nil, "MIT") {
		t.Error("contains() = true, want false for nil slice")
	}
}

func TestFilterEntries(t *testing.T) {
	t.Parallel()

	entries := []catalog.Entry{
		{ID: "go-library", Language: "go", Kind: "library"},
		{ID: "go-service", Language: "go", Kind: "service"},
		{ID: "py-library", Language: "python", Kind: "library"},
	}

	byLanguage := filterEntries(entries, "go", "")
	if len(byLanguage) != 2 {
		t.Fatalf("filterEntries(language=go) len = %d, want 2", len(byLanguage))
	}

	byBoth := filterEntries(entries, "go", "library")
	if len(byBoth) != 1 || byBoth[0].ID != "go-library" {
		t.Fatalf("filterEntries(go, library) = %#v, want [go-library]", byBoth)
	}

	none := filterEntries(entries, "rust", "")
	if len(none) != 0 {
		t.Fatalf("filterEntries(rust) len = %d, want 0", len(none))
	}
}

func TestUniqueValues(t *testing.T) {
	t.Parallel()

	entries := []catalog.Entry{{Language: "go"}, {Language: "go"}, {Language: "python"}}
	got := uniqueValues(entries, func(e catalog.Entry) string { return e.Language })
	if len(got) != 2 || got[0] != "go" || got[1] != "python" {
		t.Fatalf("uniqueValues() = %#v, want [go python] in first-seen order", got)
	}
}

func TestModuleLabel(t *testing.T) {
	t.Parallel()

	if got := moduleLabel(catalog.Entry{}); got != "Package or module path" {
		t.Errorf("moduleLabel(zero value) = %q, want default label", got)
	}
	custom := catalog.Entry{Prompts: catalog.ProjectPrompts{ModuleLabel: "Crate path"}}
	if got := moduleLabel(custom); got != "Crate path" {
		t.Errorf("moduleLabel(custom) = %q, want %q", got, "Crate path")
	}
}

// stubPrompter answers Input with a fixed value/error and fails every other
// interaction, since resolveInput only exercises Input.
type stubPrompter struct {
	inputValue string
	inputErr   error
	tux.Prompter
}

func (s stubPrompter) Input(context.Context, tux.InputRequest) (string, error) {
	return s.inputValue, s.inputErr
}

func TestResolveInputReturnsExistingValueWithoutPrompting(t *testing.T) {
	t.Parallel()

	prompter := stubPrompter{inputErr: errors.New("must not be called")}
	got, err := resolveInput(context.Background(), prompter, "already-set", tux.InputRequest{Label: "Name"}, "--name")
	if err != nil {
		t.Fatalf("resolveInput() error = %v", err)
	}
	if got != "already-set" {
		t.Fatalf("resolveInput() = %q, want %q", got, "already-set")
	}
}

func TestResolveInputPromptsWhenEmpty(t *testing.T) {
	t.Parallel()

	prompter := stubPrompter{inputValue: "from-prompt"}
	got, err := resolveInput(context.Background(), prompter, "", tux.InputRequest{Label: "Name"}, "--name")
	if err != nil {
		t.Fatalf("resolveInput() error = %v", err)
	}
	if got != "from-prompt" {
		t.Fatalf("resolveInput() = %q, want %q", got, "from-prompt")
	}
}

func TestResolveInputMapsUnavailableToActionableError(t *testing.T) {
	t.Parallel()

	prompter := stubPrompter{inputErr: tux.ErrUnavailable}
	_, err := resolveInput(context.Background(), prompter, "", tux.InputRequest{Label: "Project name"}, "--name")
	if err == nil {
		t.Fatal("resolveInput() error = nil, want an actionable error naming --name")
	}
	if errors.Is(err, tux.ErrUnavailable) {
		t.Fatal("resolveInput() should translate ErrUnavailable into a user-facing message, not pass it through")
	}
	const want = "Project name is required; provide --name"
	if err.Error() != want {
		t.Fatalf("resolveInput() error = %q, want %q", err.Error(), want)
	}
}

func TestResolveInputPropagatesOtherErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("terminal disconnected")
	prompter := stubPrompter{inputErr: sentinel}
	_, err := resolveInput(context.Background(), prompter, "", tux.InputRequest{Label: "Name"}, "--name")
	if !errors.Is(err, sentinel) {
		t.Fatalf("resolveInput() error = %v, want it to wrap %v", err, sentinel)
	}
}
