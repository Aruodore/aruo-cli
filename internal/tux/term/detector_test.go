package term_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/aruodore/aruo-cli/internal/tux"
	tuxterm "github.com/aruodore/aruo-cli/internal/tux/term"
)

type stream struct{ fd uintptr }

func (s stream) Read([]byte) (int, error)    { return 0, errors.New("not implemented") }
func (s stream) Write(p []byte) (int, error) { return len(p), nil }
func (s stream) Fd() uintptr                 { return s.fd }

type probe struct {
	terminals map[int]bool
	width     int
	height    int
}

func (p probe) IsTerminal(fd int) bool     { return p.terminals[fd] }
func (p probe) Size(int) (int, int, error) { return p.width, p.height, nil }

func TestDetectCapabilities(t *testing.T) {
	t.Parallel()

	detector := tuxterm.NewDetectorWithProbe(probe{
		terminals: map[int]bool{1: true, 2: true, 3: true},
		width:     132,
		height:    40,
	})
	capabilities := detector.Detect(stream{1}, stream{2}, stream{3}, map[string]string{
		"TERM":           "xterm-256color",
		"COLORTERM":      "truecolor",
		"LANG":           "en_US.UTF-8",
		"TERM_PROGRAM":   "WezTerm",
		"SSH_CONNECTION": "client server",
		"TMUX":           "/tmp/tmux",
	})
	if !capabilities.InputTTY || !capabilities.OutputTTY || !capabilities.ErrorTTY {
		t.Fatalf("terminal detection = %+v", capabilities)
	}
	if capabilities.Width != 132 || capabilities.Height != 40 {
		t.Fatalf("size = %dx%d", capabilities.Width, capabilities.Height)
	}
	if capabilities.Color != tux.ColorTrueColor || !capabilities.Unicode || !capabilities.Hyperlinks {
		t.Fatalf("presentation capabilities = %+v", capabilities)
	}
	if !capabilities.SSH || capabilities.Multiplexer != "tmux" {
		t.Fatalf("remote capabilities = %+v", capabilities)
	}
}

func TestDetectConservativeFallback(t *testing.T) {
	t.Parallel()

	detector := tuxterm.NewDetectorWithProbe(probe{terminals: map[int]bool{}})
	capabilities := detector.Detect(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, map[string]string{
		"TERM":     "dumb",
		"NO_COLOR": "",
		"CI":       "true",
	})
	if capabilities.Width != 80 || capabilities.Height != 24 {
		t.Fatalf("size = %dx%d", capabilities.Width, capabilities.Height)
	}
	if capabilities.Color != tux.ColorNone || capabilities.Unicode || capabilities.CursorAddressing || capabilities.Hyperlinks {
		t.Fatalf("fallback capabilities = %+v", capabilities)
	}
	if !capabilities.CI {
		t.Fatal("CI = false")
	}
}

func TestEnvironmentUsesLastValue(t *testing.T) {
	t.Parallel()

	got := tuxterm.Environment([]string{"TERM=dumb", "VALUE=with=equals", "TERM=xterm"})
	if got["TERM"] != "xterm" || got["VALUE"] != "with=equals" {
		t.Fatalf("Environment() = %#v", got)
	}
}
