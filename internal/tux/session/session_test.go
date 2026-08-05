package session_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/aruodore/aruo/internal/tux"
	"github.com/aruodore/aruo/internal/tux/charm"
	"github.com/aruodore/aruo/internal/tux/plain"
	"github.com/aruodore/aruo/internal/tux/policy"
	"github.com/aruodore/aruo/internal/tux/session"
)

type stream struct{ fd uintptr }

func (s stream) Read([]byte) (int, error)    { return 0, errors.New("not implemented") }
func (s stream) Write(p []byte) (int, error) { return len(p), nil }
func (s stream) Fd() uintptr                 { return s.fd }

type probe struct{ terminals map[int]bool }

func (p probe) IsTerminal(fd int) bool     { return p.terminals[fd] }
func (p probe) Size(int) (int, int, error) { return 80, 24, nil }

func TestNewSelectsRichAdaptersForInteractiveHuman(t *testing.T) {
	t.Parallel()

	terminalProbe := probe{terminals: map[int]bool{1: true, 2: true}}
	built := session.New(context.Background(), stream{1}, stream{1}, stream{2}, map[string]string{}, policy.Overrides{}, terminalProbe)

	if built.Policy().Mode != tux.ModeInteractive {
		t.Fatalf("Mode = %v, want interactive", built.Policy().Mode)
	}
	if _, ok := built.Prompter().(*charm.Prompter); !ok {
		t.Fatalf("Prompter() = %T, want *charm.Prompter", built.Prompter())
	}
	if _, ok := built.Presenter().(*charm.Presenter); !ok {
		t.Fatalf("Presenter() = %T, want *charm.Presenter", built.Presenter())
	}
}

func TestNewSelectsPlainAdapterWhenNotATerminal(t *testing.T) {
	t.Parallel()

	var in, out, errOut bytes.Buffer
	built := session.New(context.Background(), &in, &out, &errOut, map[string]string{}, policy.Overrides{}, probe{})

	if built.Policy().Mode == tux.ModeInteractive {
		t.Fatalf("Mode = %v, want non-interactive", built.Policy().Mode)
	}
	if _, ok := built.Prompter().(*plain.Adapter); !ok {
		t.Fatalf("Prompter() = %T, want *plain.Adapter", built.Prompter())
	}
}

func TestNewSelectsPlainAdapterWhenAccessible(t *testing.T) {
	t.Parallel()

	terminalProbe := probe{terminals: map[int]bool{1: true, 2: true}}
	accessible := true
	built := session.New(context.Background(), stream{1}, stream{1}, stream{2}, map[string]string{}, policy.Overrides{Accessible: &accessible}, terminalProbe)

	if built.Policy().Mode != tux.ModeInteractive {
		t.Fatalf("Mode = %v, want interactive", built.Policy().Mode)
	}
	if _, ok := built.Prompter().(*plain.Adapter); !ok {
		t.Fatalf("Prompter() = %T, want *plain.Adapter", built.Prompter())
	}
	if _, ok := built.Presenter().(*plain.Adapter); !ok {
		t.Fatalf("Presenter() = %T, want *plain.Adapter", built.Presenter())
	}
}

func TestNewSelectsPlainAdapterForMachineFormat(t *testing.T) {
	t.Parallel()

	terminalProbe := probe{terminals: map[int]bool{1: true, 2: true}}
	built := session.New(context.Background(), stream{1}, stream{1}, stream{2}, map[string]string{}, policy.Overrides{Format: tux.OutputJSON}, terminalProbe)

	if built.Policy().Mode != tux.ModeMachine {
		t.Fatalf("Mode = %v, want machine", built.Policy().Mode)
	}
	if _, ok := built.Presenter().(*plain.Adapter); !ok {
		t.Fatalf("Presenter() = %T, want *plain.Adapter", built.Presenter())
	}
}

func TestProgressUsesReferenceAdapterWithoutMotion(t *testing.T) {
	t.Parallel()

	terminalProbe := probe{terminals: map[int]bool{1: true, 2: true}}
	never := tux.FeatureNever
	built := session.New(context.Background(), stream{1}, stream{1}, stream{2}, map[string]string{}, policy.Overrides{Motion: &never}, terminalProbe)
	t.Cleanup(func() { _ = built.Close() })

	if _, ok := built.Progress().(*plain.Adapter); !ok {
		t.Fatalf("Progress() = %T, want *plain.Adapter", built.Progress())
	}
}

func TestCloseWithoutProgressIsSafe(t *testing.T) {
	t.Parallel()

	built := session.New(context.Background(), &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, map[string]string{}, policy.Overrides{}, probe{})
	if err := built.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
}
