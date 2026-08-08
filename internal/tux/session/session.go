// Package session assembles one invocation's terminal capabilities, policy,
// and adapters behind the tux.Session port.
package session

import (
	"context"
	"io"
	"sync"

	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/aruodore/aruo-cli/internal/tux/charm"
	"github.com/aruodore/aruo-cli/internal/tux/plain"
	"github.com/aruodore/aruo-cli/internal/tux/policy"
	"github.com/aruodore/aruo-cli/internal/tux/term"
)

// Session resolves capabilities and policy once and exposes the adapters
// chosen for one invocation. It is the only place that decides between the
// rich Charm adapters and the plain reference adapter.
type Session struct {
	capabilities tux.Capabilities
	resolved     tux.Policy
	prompter     tux.Prompter
	presenter    tux.Presenter
	reference    *plain.Adapter

	ctx          context.Context
	errOut       io.Writer
	progressOnce sync.Once
	progress     tux.ProgressSink
}

var _ tux.Session = (*Session)(nil)

// New detects terminal capabilities, resolves policy, and selects the
// prompter, presenter, and progress adapters for one command invocation.
// probe is optional; a nil probe uses the system terminal.
func New(ctx context.Context, in io.Reader, out, errOut io.Writer, environment map[string]string, overrides policy.Overrides, probe term.Probe) *Session {
	detector := term.NewDetector()
	if probe != nil {
		detector = term.NewDetectorWithProbe(probe)
	}
	capabilities := detector.Detect(in, out, errOut, environment)
	resolved := policy.Resolve(capabilities, environment, overrides)
	reference := plain.New(in, out, errOut, resolved.Mode == tux.ModeInteractive, nil)

	session := &Session{
		capabilities: capabilities,
		resolved:     resolved,
		reference:    reference,
		ctx:          ctx,
		errOut:       errOut,
	}

	// CursorAddressing is false under TERM=dumb even on a real TTY; Huh's
	// forms redraw in place and have no graceful degradation for a terminal
	// that cannot address the cursor, so fall back to the line-oriented
	// reference adapter rather than the rich one.
	if resolved.Mode == tux.ModeInteractive && !resolved.Accessible && capabilities.CursorAddressing {
		session.prompter = charm.NewPrompter(in, errOut, capabilities, resolved)
	} else {
		session.prompter = reference
	}

	if resolved.Format == tux.OutputHuman && !resolved.Accessible {
		session.presenter = charm.NewPresenter(out, errOut, capabilities, resolved)
	} else {
		session.presenter = reference
	}

	return session
}

// Capabilities returns the immutable capability snapshot for this invocation.
func (s *Session) Capabilities() tux.Capabilities { return s.capabilities }

// Policy returns the resolved presentation and input policy for this invocation.
func (s *Session) Policy() tux.Policy { return s.resolved }

// Prompter returns the selected prompt adapter.
func (s *Session) Prompter() tux.Prompter { return s.prompter }

// Presenter returns the selected output adapter.
func (s *Session) Presenter() tux.Presenter { return s.presenter }

// Progress returns the selected progress adapter, starting a live renderer
// on first use only when motion is enabled and the terminal supports it.
func (s *Session) Progress() tux.ProgressSink {
	s.progressOnce.Do(func() {
		rich := s.resolved.Mode == tux.ModeInteractive &&
			!s.resolved.Accessible &&
			s.resolved.Motion != tux.FeatureNever &&
			s.capabilities.CursorAddressing
		if rich {
			s.progress = charm.NewProgress(s.ctx, s.errOut, s.capabilities)
		} else {
			s.progress = s.reference
		}
	})
	return s.progress
}

// Close releases any live renderer started by Progress.
func (s *Session) Close() error {
	if closer, ok := s.progress.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}
