// Package policy resolves terminal behavior from capabilities and explicit choices.
package policy

import (
	"strings"

	"github.com/aruodore/aruo-cli/internal/tux"
)

// Overrides contains command-line choices that take precedence over environment hints.
type Overrides struct {
	Format     tux.OutputFormat
	Accessible *bool
	Input      *tux.FeaturePolicy
	Color      *tux.FeaturePolicy
	Unicode    *tux.FeaturePolicy
	Motion     *tux.FeaturePolicy
	Pager      *tux.FeaturePolicy
}

// Resolve produces a deterministic policy for one invocation.
func Resolve(capabilities tux.Capabilities, environment map[string]string, overrides Overrides) tux.Policy {
	format := overrides.Format
	if format == "" {
		format = tux.OutputHuman
	}
	resolved := tux.Policy{
		Format:  format,
		Input:   tux.FeatureAuto,
		Color:   tux.FeatureAuto,
		Unicode: tux.FeatureAuto,
		Motion:  tux.FeatureAuto,
		Pager:   tux.FeatureAuto,
	}
	resolved.Accessible = truthy(environment["ARUO_ACCESSIBLE"])
	if truthy(environment["NO_COLOR"]) || hasKey(environment, "NO_COLOR") {
		resolved.Color = tux.FeatureNever
	}
	if strings.EqualFold(environment["ARUO_MOTION"], "never") {
		resolved.Motion = tux.FeatureNever
	}
	if truthy(environment["ARUO_NO_INPUT"]) {
		resolved.Input = tux.FeatureNever
	}
	if capabilities.CI {
		resolved.Input = tux.FeatureNever
		resolved.Motion = tux.FeatureNever
		resolved.Pager = tux.FeatureNever
	}
	if resolved.Accessible {
		resolved.Motion = tux.FeatureNever
		resolved.Pager = tux.FeatureNever
	}
	apply(&resolved.Accessible, overrides.Accessible)
	apply(&resolved.Input, overrides.Input)
	apply(&resolved.Color, overrides.Color)
	apply(&resolved.Unicode, overrides.Unicode)
	apply(&resolved.Motion, overrides.Motion)
	apply(&resolved.Pager, overrides.Pager)

	switch {
	case format != tux.OutputHuman && format != tux.OutputPlain:
		resolved.Mode = tux.ModeMachine
		resolved.Input = tux.FeatureNever
		resolved.Color = tux.FeatureNever
		resolved.Unicode = tux.FeatureNever
		resolved.Motion = tux.FeatureNever
		resolved.Pager = tux.FeatureNever
	case format == tux.OutputPlain:
		resolved.Mode = tux.ModeNonInteractive
		resolved.Input = tux.FeatureNever
		resolved.Color = tux.FeatureNever
		resolved.Motion = tux.FeatureNever
		resolved.Pager = tux.FeatureNever
	case resolved.Input != tux.FeatureNever && capabilities.InputTTY && capabilities.ErrorTTY:
		resolved.Mode = tux.ModeInteractive
	default:
		resolved.Mode = tux.ModeNonInteractive
	}
	return resolved
}

func apply[T any](target *T, override *T) {
	if override != nil {
		*target = *override
	}
}

func truthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func hasKey(environment map[string]string, key string) bool {
	_, found := environment[key]
	return found
}
