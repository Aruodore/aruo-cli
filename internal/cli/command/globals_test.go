package command

import (
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
)

func TestParseFeaturePolicy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		value string
		want  *tux.FeaturePolicy
	}{
		{"", nil},
		{"auto", nil},
		{"always", ptr(tux.FeatureAlways)},
		{"never", ptr(tux.FeatureNever)},
		{"ALWAYS", ptr(tux.FeatureAlways)},
	}
	for _, tc := range cases {
		got, err := parseFeaturePolicy("--color", tc.value)
		if err != nil {
			t.Fatalf("parseFeaturePolicy(%q) error = %v", tc.value, err)
		}
		switch {
		case tc.want == nil && got != nil:
			t.Errorf("parseFeaturePolicy(%q) = %v, want nil", tc.value, *got)
		case tc.want != nil && (got == nil || *got != *tc.want):
			t.Errorf("parseFeaturePolicy(%q) = %v, want %v", tc.value, got, *tc.want)
		}
	}
}

func TestParseFeaturePolicyRejectsUnknownValue(t *testing.T) {
	t.Parallel()

	if _, err := parseFeaturePolicy("--motion", "sometimes"); err == nil {
		t.Fatal("parseFeaturePolicy() error = nil, want error for invalid value")
	}
}

func TestGlobalOptionsOverridesDefaults(t *testing.T) {
	t.Parallel()

	overrides, err := newGlobalOptions().overrides()
	if err != nil {
		t.Fatalf("overrides() error = %v", err)
	}
	if overrides.Input != nil || overrides.Color != nil || overrides.Motion != nil || overrides.Accessible != nil {
		t.Fatalf("overrides() = %+v, want every override unset by default", overrides)
	}
	if overrides.Format != tux.OutputHuman {
		t.Fatalf("overrides().Format = %v, want %v", overrides.Format, tux.OutputHuman)
	}
}

func TestGlobalOptionsOverridesNoInputWinsWhenBothSet(t *testing.T) {
	t.Parallel()

	options := newGlobalOptions()
	options.noInput = true
	options.interactive = true
	overrides, err := options.overrides()
	if err != nil {
		t.Fatalf("overrides() error = %v", err)
	}
	if overrides.Input == nil || *overrides.Input != tux.FeatureNever {
		t.Fatalf("overrides().Input = %v, want FeatureNever when --no-input and --interactive are both set", overrides.Input)
	}
}

func TestGlobalOptionsOverridesInteractiveForcesInput(t *testing.T) {
	t.Parallel()

	options := newGlobalOptions()
	options.interactive = true
	overrides, err := options.overrides()
	if err != nil {
		t.Fatalf("overrides() error = %v", err)
	}
	if overrides.Input == nil || *overrides.Input != tux.FeatureAlways {
		t.Fatalf("overrides().Input = %v, want FeatureAlways", overrides.Input)
	}
}

func TestGlobalOptionsOverridesRejectsInvalidColor(t *testing.T) {
	t.Parallel()

	options := newGlobalOptions()
	options.color = "purple"
	if _, err := options.overrides(); err == nil {
		t.Fatal("overrides() error = nil, want error for invalid --color value")
	}
}

func TestGlobalOptionsOverridesRejectsInvalidMotion(t *testing.T) {
	t.Parallel()

	options := newGlobalOptions()
	options.motion = "sometimes"
	if _, err := options.overrides(); err == nil {
		t.Fatal("overrides() error = nil, want error for invalid --motion value")
	}
}

func ptr(v tux.FeaturePolicy) *tux.FeaturePolicy { return &v }
