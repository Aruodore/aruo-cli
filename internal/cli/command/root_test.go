package command

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// fakeProbe reports a fixed, non-terminal environment so session.New never
// blocks on real terminal state during tests.
type fakeProbe struct{}

func (fakeProbe) IsTerminal(int) bool        { return false }
func (fakeProbe) Size(int) (int, int, error) { return 80, 24, nil }

func TestNewRootOmitsCommandsForNilDependencies(t *testing.T) {
	t.Parallel()

	root := NewRoot(iostreams.IOStreams{In: strings.NewReader(""), Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
		buildinfo.Info{Version: "dev"}, nil, nil, nil, nil, nil, fakeProbe{})
	root.InitDefaultCompletionCmd()

	names := map[string]bool{}
	for _, sub := range root.Commands() {
		names[sub.Name()] = true
	}
	if names["create"] {
		t.Error("root has a create command, want it omitted when the catalog/creator are nil")
	}
	if names["doctor"] {
		t.Error("root has a doctor command, want it omitted when the doctor service is nil")
	}
	if !names["version"] {
		t.Error("root is missing the version command, which has no dependencies")
	}
	if !names["completion"] {
		t.Error("root is missing the generated completion command")
	}
}

func TestSessionFactoryBuildForceNoInputOverridesGlobalPolicy(t *testing.T) {
	t.Parallel()

	factory := sessionFactory{
		streams:     iostreams.IOStreams{In: strings.NewReader(""), Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
		environment: map[string]string{},
		probe:       fakeProbe{},
		global:      newGlobalOptions(),
	}
	terminal, err := factory.build(context.Background(), true)
	if err != nil {
		t.Fatalf("build() error = %v", err)
	}
	defer func() { _ = terminal.Close() }()
	if terminal.Policy().Mode == tux.ModeInteractive {
		t.Error("Policy().Mode = interactive, want forceNoInput to disable prompting regardless of global flags")
	}
}

func TestSessionFactoryBuildPropagatesInvalidGlobalFlags(t *testing.T) {
	t.Parallel()

	global := newGlobalOptions()
	global.color = "not-a-policy"
	factory := sessionFactory{
		streams:     iostreams.IOStreams{In: strings.NewReader(""), Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}},
		environment: map[string]string{},
		probe:       fakeProbe{},
		global:      global,
	}
	if _, err := factory.build(context.Background(), false); err == nil {
		t.Fatal("build() error = nil, want the invalid --color value to surface")
	}
}
