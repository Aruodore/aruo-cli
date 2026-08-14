package charm_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/aruodore/aruo-cli/internal/tux/charm"
)

func TestMessageDegradesWithoutColorOrUnicode(t *testing.T) {
	t.Parallel()

	var result bytes.Buffer
	presenter := charm.NewPresenter(&result, &bytes.Buffer{}, tux.Capabilities{}, tux.Policy{Color: tux.FeatureNever, Unicode: tux.FeatureNever})
	if err := presenter.Message(context.Background(), tux.Message{Kind: tux.MessageSuccess, Text: "created project"}); err != nil {
		t.Fatal(err)
	}
	if got := result.String(); got != "Success: created project\n" || strings.Contains(got, "\x1b[") {
		t.Fatalf("result = %q", got)
	}
}

func TestMessageUsesSemanticRichPresentation(t *testing.T) {
	t.Parallel()

	var result bytes.Buffer
	presenter := charm.NewPresenter(&result, &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorTrueColor, Unicode: true}, tux.Policy{})
	if err := presenter.Message(context.Background(), tux.Message{Kind: tux.MessageSuccess, Text: "created project"}); err != nil {
		t.Fatal(err)
	}
	if got := result.String(); !strings.Contains(got, "✓") || !strings.Contains(got, "created project") || !strings.Contains(got, "\x1b[") {
		t.Fatalf("result = %q", got)
	}
}

func TestTableDropsLowerPriorityColumns(t *testing.T) {
	t.Parallel()

	var result bytes.Buffer
	presenter := charm.NewPresenter(&result, &bytes.Buffer{}, tux.Capabilities{Width: 18}, tux.Policy{Color: tux.FeatureNever})
	table := tux.Table{
		Columns: []tux.Column{
			{ID: "name", Heading: "Name", Priority: 0},
			{ID: "score", Heading: "Score", Alignment: tux.AlignRight, Priority: 1},
			{ID: "detail", Heading: "Detail", Priority: 10},
			{ID: "secret", Heading: "Secret", Priority: 0, Sensitive: true},
		},
		Rows: [][]string{{"Aruo", "98", "production ready", "token"}},
	}
	if err := presenter.Table(context.Background(), table); err != nil {
		t.Fatal(err)
	}
	got := result.String()
	if !strings.Contains(got, "NAME") || !strings.Contains(got, "SCORE") || strings.Contains(got, "DETAIL") || strings.Contains(got, "token") {
		t.Fatalf("result = %q", got)
	}
}

func TestDiagnosticUsesDiagnosticStream(t *testing.T) {
	t.Parallel()

	var result bytes.Buffer
	var diagnostic bytes.Buffer
	presenter := charm.NewPresenter(&result, &diagnostic, tux.Capabilities{}, tux.Policy{Color: tux.FeatureNever})
	if err := presenter.Diagnostic(context.Background(), tux.Diagnostic{
		Summary:    "target exists",
		Effect:     "Nothing was changed.",
		Suggestion: "Choose another path.",
	}); err != nil {
		t.Fatal(err)
	}
	if result.Len() != 0 || !strings.Contains(diagnostic.String(), "Error: target exists") {
		t.Fatalf("result=%q diagnostic=%q", result.String(), diagnostic.String())
	}
}
