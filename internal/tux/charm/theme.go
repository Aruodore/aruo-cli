package charm

import (
	"image/color"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// theme is intentionally small. Accent communicates focus and structure;
// outcome colors are reserved for states with those meanings.
type theme struct {
	accent   lipgloss.Style
	emphasis lipgloss.Style
	success  lipgloss.Style
	warning  lipgloss.Style
	danger   lipgloss.Style
	muted    lipgloss.Style
}

func newTheme(capabilities tux.Capabilities, policy tux.Policy) theme {
	if policy.Color == tux.FeatureNever || capabilities.Color == tux.ColorNone {
		return theme{
			accent:   lipgloss.NewStyle(),
			emphasis: lipgloss.NewStyle(),
			success:  lipgloss.NewStyle(),
			warning:  lipgloss.NewStyle(),
			danger:   lipgloss.NewStyle().Bold(true),
			muted:    lipgloss.NewStyle(),
		}
	}
	return theme{
		accent:   lipgloss.NewStyle().Bold(true),
		emphasis: lipgloss.NewStyle().Bold(true),
		success:  lipgloss.NewStyle().Foreground(semanticColor(capabilities.Color, "2", "71", "#4F8A63")),
		warning:  lipgloss.NewStyle().Foreground(semanticColor(capabilities.Color, "3", "136", "#9A6B18")),
		danger:   lipgloss.NewStyle().Bold(true).Foreground(semanticColor(capabilities.Color, "1", "167", "#B54F4F")),
		muted:    lipgloss.NewStyle().Foreground(semanticColor(capabilities.Color, "8", "243", "#77736C")),
	}
}

func semanticColor(level tux.ColorLevel, ansi16, ansi256, trueColor string) color.Color {
	switch level {
	case tux.ColorTrueColor:
		return lipgloss.Color(trueColor)
	case tux.ColorANSI256:
		return lipgloss.Color(ansi256)
	default:
		return lipgloss.Color(ansi16)
	}
}

func newHuhTheme(capabilities tux.Capabilities, policy tux.Policy) huh.Theme {
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		styles := huh.ThemeBase(isDark)
		visual := newPromptTheme(capabilities, policy, isDark)

		styles.Form.Base = lipgloss.NewStyle()
		styles.FieldSeparator = lipgloss.NewStyle().SetString("\n")
		styles.Focused.Base = lipgloss.NewStyle().PaddingLeft(1)
		styles.Focused.Card = styles.Focused.Base
		styles.Focused.Title = visual.emphasis
		styles.Focused.Description = visual.muted.MarginBottom(1)
		styles.Focused.ErrorIndicator = visual.danger.SetString(" !")
		styles.Focused.ErrorMessage = visual.danger
		styles.Focused.SelectSelector = lipgloss.NewStyle().SetString("  ")
		styles.Focused.MultiSelectSelector = lipgloss.NewStyle()
		styles.Focused.SelectedPrefix = visual.emphasis.SetString("[x] ")
		styles.Focused.UnselectedPrefix = visual.muted.SetString("[ ] ")
		// Transform wraps the complete rendered label instead of styling each
		// rune. This preserves nested ecosystem colors while underlining the full
		// active sentence; Style.Underline corrupts nested ANSI sequences.
		styles.Focused.SelectedOption = lipgloss.NewStyle().Transform(func(label string) string {
			return "\x1b[1;4m" + label + "\x1b[22;24m"
		})
		styles.Focused.UnselectedOption = lipgloss.NewStyle()
		styles.Focused.Option = lipgloss.NewStyle()
		styles.Focused.NextIndicator = visual.emphasis.SetString("→")
		styles.Focused.PrevIndicator = visual.muted.SetString("←")
		styles.Focused.FocusedButton = visual.emphasis.Reverse(true).Padding(0, 2).MarginRight(1)
		styles.Focused.BlurredButton = visual.muted.Padding(0, 2).MarginRight(1)
		styles.Focused.TextInput.Cursor = lipgloss.NewStyle()
		styles.Focused.TextInput.Prompt = visual.emphasis
		styles.Focused.TextInput.Placeholder = visual.muted

		styles.Blurred = styles.Focused
		styles.Blurred.Base = lipgloss.NewStyle().PaddingLeft(1)
		styles.Blurred.Card = styles.Blurred.Base
		styles.Blurred.Title = visual.muted
		styles.Blurred.SelectSelector = lipgloss.NewStyle().SetString("  ")
		styles.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
		styles.Blurred.NextIndicator = lipgloss.NewStyle()
		styles.Blurred.PrevIndicator = lipgloss.NewStyle()

		styles.Group.Title = visual.emphasis.MarginBottom(1)
		styles.Group.Description = visual.muted.MarginBottom(1)
		styles.Help.ShortKey = visual.emphasis
		styles.Help.ShortDesc = visual.muted
		styles.Help.ShortSeparator = visual.muted
		styles.Help.FullKey = visual.emphasis
		styles.Help.FullDesc = visual.muted
		styles.Help.FullSeparator = visual.muted
		styles.Help.Ellipsis = visual.muted
		return styles
	})
}

func newHighlightedHuhTheme(capabilities tux.Capabilities, policy tux.Policy) huh.Theme {
	base := newHuhTheme(capabilities, policy)
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		styles := base.Theme(isDark)
		styles.Focused.SelectedOption = lipgloss.NewStyle().Transform(func(label string) string {
			return "\x1b[7m" + label + "\x1b[27m"
		})
		return styles
	})
}

func newPromptTheme(capabilities tux.Capabilities, policy tux.Policy, isDark bool) theme {
	if policy.Accessible || policy.Color == tux.FeatureNever || capabilities.Color == tux.ColorNone {
		return newTheme(tux.Capabilities{}, tux.Policy{Color: tux.FeatureNever})
	}
	visual := newTheme(capabilities, policy)
	if isDark {
		visual.muted = lipgloss.NewStyle().Foreground(semanticColor(capabilities.Color, "7", "249", "#A7ABB8"))
	} else {
		visual.muted = lipgloss.NewStyle().Foreground(semanticColor(capabilities.Color, "8", "240", "#555A66"))
	}
	return visual
}
