package charm

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
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

func TestAccessibleInputNeverEmitsThemeEscapes(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	prompter := NewPrompter(strings.NewReader("Aruo\n"), &output, tux.Capabilities{Width: 80, Color: tux.ColorTrueColor}, tux.Policy{
		Mode:       tux.ModeInteractive,
		Accessible: true,
	})
	if _, err := prompter.Input(context.Background(), tux.InputRequest{ID: "name", Label: "Project name"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("accessible output contains terminal styling: %q", output.String())
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
	if options[1].Value != "go-cli" || options[1].Key != "Go CLI  Command-line application" {
		t.Fatalf("option = %+v", options[1])
	}
}

func TestHuhOptionsApplyCatalogColorsOnlyToRichTrueColorMenus(t *testing.T) {
	t.Parallel()

	colored := []tux.Option{{ID: "go-library", Label: "Go library", Color: "#00ADD8"}}

	trueColor := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{})
	options := trueColor.huhOptions(colored, nil)
	if options[0].Key == "Go library" || !strings.Contains(options[0].Key, "38;2;0;173;216") {
		t.Errorf("huhOptions() = %q, want the catalog ecosystem color", options[0].Key)
	}

	accessible := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{Accessible: true})
	if got := accessible.huhOptions(colored, nil)[0].Key; got != "Go library" {
		t.Errorf("accessible huhOptions() = %q, want no ANSI styling", got)
	}

	ansi256 := NewPrompter(strings.NewReader(""), &bytes.Buffer{}, tux.Capabilities{Color: tux.ColorANSI256}, tux.Policy{})
	if got := ansi256.huhOptions(colored, nil)[0].Key; got != "Go library" {
		t.Errorf("ANSI-256 huhOptions() = %q, want no inaccurate color approximation", got)
	}
}

func TestAruoPromptThemeUnderlinesFullSelectionWithoutReplacingCatalogColor(t *testing.T) {
	t.Parallel()

	styles := newHuhTheme(tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{}).Theme(true)
	if got := styles.Focused.SelectSelector.Render(""); strings.TrimSpace(got) != "" {
		t.Fatalf("focused selector = %q, want spacing only", got)
	}
	if got := styles.Blurred.SelectSelector.Render(""); strings.TrimSpace(got) != "" {
		t.Fatalf("blurred selector = %q", got)
	}
	if styles.Focused.Base.GetBorderLeft() {
		t.Fatal("focused field has a decorative border")
	}
	colored := "\x1b[38;2;97;218;251mReact\x1b[m"
	if got := styles.Focused.SelectedOption.Render(colored); got != "\x1b[1;4m"+colored+"\x1b[22;24m" {
		t.Fatalf("selected option = %q, want one full-label underline preserving its ANSI color", got)
	}
	if rendered := styles.Focused.Title.Render("Project name"); strings.Contains(rendered, "38;2;") {
		t.Fatalf("field title uses a decorative brand color: %q", rendered)
	}
}

func TestHighlightedPromptThemeFillsFullActiveOption(t *testing.T) {
	t.Parallel()

	styles := newHighlightedHuhTheme(tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{}).Theme(true)
	if got := styles.Focused.SelectedOption.Render("Application"); got != "\x1b[7mApplication\x1b[27m" {
		t.Fatalf("selected option = %q, want full-label highlight", got)
	}
}

func TestPromptSecondaryTextAdaptsToTerminalBackground(t *testing.T) {
	t.Parallel()

	theme := newHuhTheme(tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{})
	dark := theme.Theme(true).Focused.Description.Render("Secondary")
	if !strings.Contains(dark, "38;2;167;171;184") {
		t.Fatalf("dark-terminal secondary text = %q, want readable light neutral", dark)
	}
	light := theme.Theme(false).Focused.Description.Render("Secondary")
	if !strings.Contains(light, "38;2;85;90;102") {
		t.Fatalf("light-terminal secondary text = %q, want readable dark neutral", light)
	}
}
