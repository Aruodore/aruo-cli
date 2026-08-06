package command

import (
	"fmt"
	"strings"

	"github.com/aruodore/aruo/internal/tux"
	"github.com/aruodore/aruo/internal/tux/policy"
	"github.com/spf13/cobra"
)

// globalOptions carries the terminal-behavior flags shared by every command.
type globalOptions struct {
	noInput     bool
	interactive bool
	accessible  bool
	color       string
	motion      string
}

func newGlobalOptions() *globalOptions {
	return &globalOptions{color: "auto", motion: "auto"}
}

var featurePolicyCompletions = func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveNoFileComp
}

func (g *globalOptions) register(root *cobra.Command) {
	flags := root.PersistentFlags()
	flags.BoolVar(&g.noInput, "no-input", false, "disable prompts and fail on missing input")
	flags.BoolVar(&g.interactive, "interactive", false, "allow prompts even when CI or redirection would otherwise disable them")
	flags.BoolVar(&g.accessible, "accessible", false, "use line-oriented accessible prompts and output")
	flags.StringVar(&g.color, "color", "auto", "color policy: auto, always, or never")
	flags.StringVar(&g.motion, "motion", "auto", "animation policy: auto, always, or never")
	_ = root.RegisterFlagCompletionFunc("color", featurePolicyCompletions)
	_ = root.RegisterFlagCompletionFunc("motion", featurePolicyCompletions)
}

// overrides converts the resolved flags into policy overrides for one command's format.
func (g *globalOptions) overrides(format tux.OutputFormat) (policy.Overrides, error) {
	result := policy.Overrides{Format: format}
	if g.accessible {
		accessible := true
		result.Accessible = &accessible
	}
	switch {
	case g.noInput:
		never := tux.FeatureNever
		result.Input = &never
	case g.interactive:
		always := tux.FeatureAlways
		result.Input = &always
	}
	color, err := parseFeaturePolicy("--color", g.color)
	if err != nil {
		return policy.Overrides{}, err
	}
	result.Color = color
	motion, err := parseFeaturePolicy("--motion", g.motion)
	if err != nil {
		return policy.Overrides{}, err
	}
	result.Motion = motion
	return result, nil
}

func parseFeaturePolicy(flag, value string) (*tux.FeaturePolicy, error) {
	switch strings.ToLower(value) {
	case "", "auto":
		return nil, nil
	case "always":
		resolved := tux.FeatureAlways
		return &resolved, nil
	case "never":
		resolved := tux.FeatureNever
		return &resolved, nil
	default:
		return nil, fmt.Errorf("invalid %s %q; use auto, always, or never", flag, value)
	}
}
