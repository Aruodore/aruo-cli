package policy_test

import (
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/aruodore/aruo-cli/internal/tux/policy"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	enabled := tux.FeatureAlways
	tests := []struct {
		name         string
		capabilities tux.Capabilities
		environment  map[string]string
		overrides    policy.Overrides
		wantMode     tux.Mode
		wantInput    tux.FeaturePolicy
		wantColor    tux.FeaturePolicy
		wantMotion   tux.FeaturePolicy
	}{
		{
			name:         "terminal selects interactive human mode",
			capabilities: tux.Capabilities{InputTTY: true, ErrorTTY: true},
			wantMode:     tux.ModeInteractive,
			wantInput:    tux.FeatureAuto,
			wantColor:    tux.FeatureAuto,
			wantMotion:   tux.FeatureAuto,
		},
		{
			name:         "CI disables transient interaction",
			capabilities: tux.Capabilities{InputTTY: true, ErrorTTY: true, CI: true},
			wantMode:     tux.ModeNonInteractive,
			wantInput:    tux.FeatureNever,
			wantColor:    tux.FeatureAuto,
			wantMotion:   tux.FeatureNever,
		},
		{
			name:         "explicit interaction overrides CI input",
			capabilities: tux.Capabilities{InputTTY: true, ErrorTTY: true, CI: true},
			overrides:    policy.Overrides{Input: &enabled},
			wantMode:     tux.ModeInteractive,
			wantInput:    tux.FeatureAlways,
			wantColor:    tux.FeatureAuto,
			wantMotion:   tux.FeatureNever,
		},
		{
			name:         "machine format disables presentation features",
			capabilities: tux.Capabilities{InputTTY: true, ErrorTTY: true},
			overrides:    policy.Overrides{Format: tux.OutputJSON},
			wantMode:     tux.ModeMachine,
			wantInput:    tux.FeatureNever,
			wantColor:    tux.FeatureNever,
			wantMotion:   tux.FeatureNever,
		},
		{
			name:         "plain format remains human and non-interactive",
			capabilities: tux.Capabilities{InputTTY: true, ErrorTTY: true},
			overrides:    policy.Overrides{Format: tux.OutputPlain},
			wantMode:     tux.ModeNonInteractive,
			wantInput:    tux.FeatureNever,
			wantColor:    tux.FeatureNever,
			wantMotion:   tux.FeatureNever,
		},
		{
			name:        "NO_COLOR presence disables color",
			environment: map[string]string{"NO_COLOR": ""},
			wantMode:    tux.ModeNonInteractive,
			wantInput:   tux.FeatureAuto,
			wantColor:   tux.FeatureNever,
			wantMotion:  tux.FeatureAuto,
		},
		{
			name:        "accessible mode disables motion",
			environment: map[string]string{"ARUO_ACCESSIBLE": "1"},
			wantMode:    tux.ModeNonInteractive,
			wantInput:   tux.FeatureAuto,
			wantColor:   tux.FeatureAuto,
			wantMotion:  tux.FeatureNever,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := policy.Resolve(test.capabilities, test.environment, test.overrides)
			if got.Mode != test.wantMode || got.Input != test.wantInput || got.Color != test.wantColor || got.Motion != test.wantMotion {
				t.Fatalf("Resolve() = %+v", got)
			}
		})
	}
}
