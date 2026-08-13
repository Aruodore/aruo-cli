package command

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/tux"
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

func TestResolveInputOptionalUnavailableUsesDefaultInsteadOfErroring(t *testing.T) {
	t.Parallel()

	placeholder := "TODO: describe this project."
	prompter := stubPrompter{inputErr: tux.ErrUnavailable}
	got, err := resolveInput(context.Background(), prompter, "", tux.InputRequest{
		Label: "Short description", Optional: true, Default: &placeholder,
	}, "--description")
	if err != nil {
		t.Fatalf("resolveInput() error = %v, want nil for an optional field with a default", err)
	}
	if got != placeholder {
		t.Fatalf("resolveInput() = %q, want the default %q", got, placeholder)
	}
}

func TestResolveInputOptionalUnavailableWithoutDefaultReturnsEmpty(t *testing.T) {
	t.Parallel()

	prompter := stubPrompter{inputErr: tux.ErrUnavailable}
	got, err := resolveInput(context.Background(), prompter, "", tux.InputRequest{
		Label: "Author or organization", Optional: true,
	}, "--author")
	if err != nil {
		t.Fatalf("resolveInput() error = %v, want nil for an optional field with no default", err)
	}
	if got != "" {
		t.Fatalf("resolveInput() = %q, want empty string", got)
	}
}

// fakeCatalog is a minimal catalog.Catalog over an in-memory entry list,
// so runGuide's step-building logic can be tested without depending on
// the real builtin catalog's specific templates.
type fakeCatalog struct {
	entries []catalog.Entry
}

func (f fakeCatalog) List(context.Context) ([]catalog.Entry, error) { return f.entries, nil }

func (f fakeCatalog) Resolve(_ context.Context, id string) (catalog.Entry, error) {
	for _, entry := range f.entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return catalog.Entry{}, fmt.Errorf("unknown template %q", id)
}

// stubGuidePrompter answers Guide with a fixed Answers/error while
// recording the steps it received, so a test can drive individual step
// builders directly without a real terminal. When set, inspect runs
// before Guide returns -- runGuide mutates its createOptions immediately
// after Guide returns, and a Step's Skip/request builder reads that same
// struct live, so assertions about pre-mutation state (like whether a
// step was offered because no flag answered it yet) must run inside
// inspect rather than after runGuide has fully returned.
type stubGuidePrompter struct {
	answers  tux.Answers
	err      error
	gotSteps []tux.Step
	inspect  func([]tux.Step)
	tux.Prompter
}

func (s *stubGuidePrompter) Guide(_ context.Context, steps []tux.Step) (tux.Answers, error) {
	s.gotSteps = steps
	if s.inspect != nil {
		s.inspect(steps)
	}
	return s.answers, s.err
}

func findGuideStep(t *testing.T, steps []tux.Step, id string) tux.Step {
	t.Helper()
	for _, step := range steps {
		if step.ID == id {
			return step
		}
	}
	t.Fatalf("no step with ID %q among %d steps", id, len(steps))
	return tux.Step{}
}

func TestRunGuideSkipsKindStepWhenCatalogHasOneKind(t *testing.T) {
	t.Parallel()

	templateCatalog := fakeCatalog{entries: []catalog.Entry{{ID: "go-library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"}}}
	prompter := &stubGuidePrompter{answers: tux.Answers{"template": tux.OptionID("go-library"), "confirm": true}}
	prompter.inspect = func(steps []tux.Step) {
		kind := findGuideStep(t, steps, "kind")
		if kind.Skip == nil || !kind.Skip() {
			t.Error("kind step Skip() = false, want true when the catalog only has one kind")
		}
	}
	options := &createOptions{module: "example.com/x"}
	if _, err := runGuide(context.Background(), prompter, templateCatalog, options); err != nil {
		t.Fatalf("runGuide() error = %v", err)
	}
}

func TestRunGuideKindStepOffersEcosystemNamesPerKind(t *testing.T) {
	t.Parallel()

	templateCatalog := fakeCatalog{entries: []catalog.Entry{
		{ID: "next", Name: "Next.js", Kind: "app", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
		{ID: "go-library", Name: "Go library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
		{ID: "js-library", Name: "JavaScript library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
	}}
	prompter := &stubGuidePrompter{answers: tux.Answers{"kind": tux.OptionID("library"), "template": tux.OptionID("go-library"), "confirm": true}}
	prompter.inspect = func(steps []tux.Step) {
		kind := findGuideStep(t, steps, "kind")
		if kind.Skip != nil && kind.Skip() {
			t.Fatal("kind step Skip() = true, want false with two kinds present")
		}
		request := kind.Select(tux.Answers{})
		if len(request.Options) != 2 {
			t.Fatalf("kind options = %#v, want 2 kinds", request.Options)
		}
		app, library := request.Options[0], request.Options[1]
		if app.ID != "app" || app.Label != "Application" || app.Description != "Next.js" {
			t.Errorf("app option = %+v", app)
		}
		if library.ID != "library" || library.Label != "Library" || library.Description != "Go library, JavaScript library" {
			t.Errorf("library option = %+v", library)
		}
	}
	options := &createOptions{module: "example.com/x"}
	if _, err := runGuide(context.Background(), prompter, templateCatalog, options); err != nil {
		t.Fatalf("runGuide() error = %v", err)
	}
}

func TestRunGuideTemplateStepNarrowsToChosenKind(t *testing.T) {
	t.Parallel()

	templateCatalog := fakeCatalog{entries: []catalog.Entry{
		{ID: "next", Kind: "app", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
		{ID: "go-library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
		{ID: "js-library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
	}}
	prompter := &stubGuidePrompter{answers: tux.Answers{"kind": tux.OptionID("library"), "template": tux.OptionID("go-library")}}
	options := &createOptions{yes: true, module: "example.com/x"}
	if _, err := runGuide(context.Background(), prompter, templateCatalog, options); err != nil {
		t.Fatalf("runGuide() error = %v", err)
	}
	template := findGuideStep(t, prompter.gotSteps, "template")
	request := template.Select(tux.Answers{"kind": tux.OptionID("library")})
	if len(request.Options) != 2 {
		t.Fatalf("template options = %#v, want the 2 library templates only", request.Options)
	}
	for _, option := range request.Options {
		if option.ID == "next" {
			t.Errorf("template options include %q, want it narrowed away once kind=library", option.ID)
		}
	}
}

func TestRunGuideCopiesAnswersBackIntoOptionsAndResolvesEntry(t *testing.T) {
	t.Parallel()

	templateCatalog := fakeCatalog{entries: []catalog.Entry{
		{ID: "go-library", Name: "Go library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
	}}
	prompter := &stubGuidePrompter{answers: tux.Answers{
		"name":        "my-app",
		"template":    tux.OptionID("go-library"),
		"description": "A test project.",
		"author":      "Jane Doe",
		"confirm":     true,
		// No "module" answer at all: there's no module step to produce
		// one anymore. Even if a stub Guide somehow returned one, it
		// must be ignored -- module always defaults to the name.
		"module": "should-be-ignored",
	}}
	options := &createOptions{}
	entry, err := runGuide(context.Background(), prompter, templateCatalog, options)
	if err != nil {
		t.Fatalf("runGuide() error = %v", err)
	}
	if entry.ID != "go-library" {
		t.Errorf("entry.ID = %q, want go-library", entry.ID)
	}
	if options.name != "my-app" || options.destination != "my-app" {
		t.Errorf("options.name/destination = %q/%q, want my-app/my-app", options.name, options.destination)
	}
	if options.module != "my-app" {
		t.Errorf("options.module = %q, want it to default to the project name, ignoring any stray Guide answer", options.module)
	}
	if options.templateID != "go-library" ||
		options.description != "A test project." || options.author != "Jane Doe" {
		t.Errorf("options = %+v, want the Guide answers copied in", options)
	}
}

func TestRunGuideReturnsCancelledWhenConfirmDeclined(t *testing.T) {
	t.Parallel()

	templateCatalog := fakeCatalog{entries: []catalog.Entry{
		{ID: "go-library", Kind: "library", Licenses: []string{"MIT"}, DefaultLicense: "MIT"},
	}}
	prompter := &stubGuidePrompter{answers: tux.Answers{
		"name": "my-app", "template": tux.OptionID("go-library"), "confirm": false,
	}}
	options := &createOptions{}
	_, err := runGuide(context.Background(), prompter, templateCatalog, options)
	if err == nil || err.Error() != "creation cancelled" {
		t.Fatalf("runGuide() error = %v, want \"creation cancelled\"", err)
	}
}

func TestRunGuidePropagatesGuideError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("terminal disconnected")
	templateCatalog := fakeCatalog{entries: []catalog.Entry{{ID: "go-library", Kind: "library"}}}
	prompter := &stubGuidePrompter{err: sentinel}
	_, err := runGuide(context.Background(), prompter, templateCatalog, &createOptions{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("runGuide() error = %v, want it to propagate %v", err, sentinel)
	}
}
