package charm

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aruodore/aruo/internal/tux"
)

func TestAccessibleInputUsesAruoRequest(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	prompter := NewPrompter(strings.NewReader("Aruo\n"), &output, tux.Capabilities{Width: 80}, tux.Policy{
		Mode:       tux.ModeInteractive,
		Accessible: true,
	})
	value, err := prompter.Input(context.Background(), tux.InputRequest{
		ID:          "name",
		Label:       "Project name",
		Description: "Name the new project.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != "Aruo" {
		t.Fatalf("Input() = %q", value)
	}
	if !strings.Contains(output.String(), "Project name") {
		t.Fatalf("output = %q", output.String())
	}
}

func TestAccessibleInputSubmittingEmptyUsesDefaultWithoutPreFilling(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	defaultValue := "Aruodore"
	prompter := NewPrompter(strings.NewReader("\n"), &output, tux.Capabilities{Width: 80}, tux.Policy{
		Mode:       tux.ModeInteractive,
		Accessible: true,
	})
	value, err := prompter.Input(context.Background(), tux.InputRequest{
		ID:       "author",
		Label:    "Author or organization",
		Optional: true,
		Default:  &defaultValue,
	})
	if err != nil {
		t.Fatalf("Input() error = %v, want a bare Enter on an optional field with a default to succeed", err)
	}
	if value != defaultValue {
		t.Fatalf("Input() = %q, want the default %q substituted after an empty submission", value, defaultValue)
	}
}

func TestPrompterRefusesNonInteractiveSession(t *testing.T) {
	t.Parallel()

	prompter := NewPrompter(strings.NewReader("ignored\n"), &bytes.Buffer{}, tux.Capabilities{}, tux.Policy{Mode: tux.ModeNonInteractive})
	_, err := prompter.Input(context.Background(), tux.InputRequest{Label: "Name"})
	if !errors.Is(err, tux.ErrUnavailable) {
		t.Fatalf("Input() error = %v", err)
	}
}

func TestPrompterMapsCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	prompter := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{}, tux.Policy{Mode: tux.ModeInteractive})
	_, err := prompter.Confirm(ctx, tux.ConfirmRequest{Label: "Continue?"})
	if !errors.Is(err, tux.ErrCancelled) || !errors.Is(err, context.Canceled) {
		t.Fatalf("Confirm() error = %v", err)
	}
}

func TestHuhOptionsExcludeDisabledAndPreserveDefaults(t *testing.T) {
	t.Parallel()

	prompter := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{}, tux.Policy{})
	defaults := map[tux.OptionID]struct{}{"go-cli": {}}
	options := prompter.huhOptions([]tux.Option{
		{ID: "go-library", Label: "Go library"},
		{ID: "go-cli", Label: "Go CLI", Description: "Command-line application"},
		{ID: "legacy", Label: "Legacy", Disabled: true},
	}, defaults)
	if len(options) != 2 {
		t.Fatalf("len(options) = %d", len(options))
	}
	if options[1].Value != "go-cli" || options[1].Key != "Go CLI — Command-line application" {
		t.Fatalf("option = %+v", options[1])
	}
}

func TestHuhOptionsAppliesColorOnlyWhenTrueColorAndPolicyAllow(t *testing.T) {
	t.Parallel()

	colored := []tux.Option{{ID: "go-library", Label: "Go library", Color: "#00ADD8"}}

	trueColor := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{})
	options := trueColor.huhOptions(colored, nil)
	if options[0].Key == "Go library" {
		t.Error("huhOptions() did not style the label under true-color capability")
	}

	noColorPolicy := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{Color: tux.FeatureNever})
	options = noColorPolicy.huhOptions(colored, nil)
	if options[0].Key != "Go library" {
		t.Errorf("huhOptions() = %q, want unstyled when --color=never overrides true-color capability", options[0].Key)
	}

	ansi256 := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorANSI256}, tux.Policy{})
	options = ansi256.huhOptions(colored, nil)
	if options[0].Key != "Go library" {
		t.Errorf("huhOptions() = %q, want unstyled below true-color capability rather than an approximated color", options[0].Key)
	}
}
