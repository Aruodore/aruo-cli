// Package term adapts operating-system terminal primitives to Aruo capabilities.
package term

import (
	"io"
	"os"
	"strings"

	"github.com/aruodore/aruo/internal/tux"
	xterm "golang.org/x/term"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

type descriptor interface {
	Fd() uintptr
}

// Probe isolates low-level terminal calls for qualification tests.
type Probe interface {
	IsTerminal(int) bool
	Size(int) (int, int, error)
}

type systemProbe struct{}

func (systemProbe) IsTerminal(fd int) bool        { return xterm.IsTerminal(fd) }
func (systemProbe) Size(fd int) (int, int, error) { return xterm.GetSize(fd) }

// Detector calculates one immutable capability snapshot.
type Detector struct {
	probe Probe
}

// NewDetector constructs a detector with the system terminal probe.
func NewDetector() Detector {
	return Detector{probe: systemProbe{}}
}

// NewDetectorWithProbe constructs a detector with an injected qualification probe.
func NewDetectorWithProbe(probe Probe) Detector {
	return Detector{probe: probe}
}

// Detect inspects explicit streams and environment without taking terminal ownership.
func (d Detector) Detect(in io.Reader, out, errOut io.Writer, environment map[string]string) tux.Capabilities {
	capabilities := tux.Capabilities{
		Width:  defaultWidth,
		Height: defaultHeight,
		CI:     detectCI(environment),
		SSH:    environment["SSH_CONNECTION"] != "" || environment["SSH_CLIENT"] != "" || environment["SSH_TTY"] != "",
	}
	_, capabilities.InputTTY = d.isTerminal(in)
	_, capabilities.OutputTTY = d.isTerminal(out)
	_, capabilities.ErrorTTY = d.isTerminal(errOut)
	if environment["TMUX"] != "" {
		capabilities.Multiplexer = "tmux"
	} else if environment["STY"] != "" {
		capabilities.Multiplexer = "screen"
	}
	termName := strings.ToLower(environment["TERM"])
	dumb := termName == "dumb"
	capabilities.Unicode = !dumb && unicodeLocale(environment)
	capabilities.CursorAddressing = !dumb && (capabilities.OutputTTY || capabilities.ErrorTTY)
	capabilities.AlternateScreen = capabilities.CursorAddressing && capabilities.Multiplexer != "screen"
	capabilities.Color = detectColor(capabilities, environment, dumb)
	capabilities.Hyperlinks = detectHyperlinks(capabilities, environment, dumb)
	for _, stream := range []any{errOut, out} {
		fd, terminal := d.isTerminal(stream)
		if !terminal {
			continue
		}
		if width, height, err := d.probe.Size(fd); err == nil && width > 0 && height > 0 {
			capabilities.Width = width
			capabilities.Height = height
			break
		}
	}
	return capabilities
}

// Environment converts an environment vector into last-value-wins lookup data.
func Environment(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, value := range values {
		key, content, found := strings.Cut(value, "=")
		if found {
			result[key] = content
		}
	}
	return result
}

// SystemEnvironment returns the current process environment as lookup data.
func SystemEnvironment() map[string]string {
	return Environment(os.Environ())
}

func (d Detector) isTerminal(stream any) (int, bool) {
	file, ok := stream.(descriptor)
	if !ok {
		return 0, false
	}
	value := file.Fd()
	if value > uintptr(^uint(0)>>1) {
		return 0, false
	}
	fd := int(value)
	return fd, d.probe.IsTerminal(fd)
}

func detectColor(capabilities tux.Capabilities, environment map[string]string, dumb bool) tux.ColorLevel {
	if dumb || !capabilities.OutputTTY && !capabilities.ErrorTTY {
		return tux.ColorNone
	}
	if _, disabled := environment["NO_COLOR"]; disabled {
		return tux.ColorNone
	}
	colorTerm := strings.ToLower(environment["COLORTERM"])
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return tux.ColorTrueColor
	}
	if strings.Contains(strings.ToLower(environment["TERM"]), "256color") {
		return tux.ColorANSI256
	}
	return tux.ColorANSI16
}

func detectHyperlinks(capabilities tux.Capabilities, environment map[string]string, dumb bool) bool {
	if dumb || !capabilities.OutputTTY && !capabilities.ErrorTTY {
		return false
	}
	if environment["WT_SESSION"] != "" {
		return true
	}
	switch strings.ToLower(environment["TERM_PROGRAM"]) {
	case "iterm.app", "wezterm", "vscode":
		return true
	default:
		return false
	}
}

func detectCI(environment map[string]string) bool {
	for _, key := range []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "BUILDKITE", "TF_BUILD", "TEAMCITY_VERSION"} {
		if value, found := environment[key]; found && (value == "" || !strings.EqualFold(value, "false")) {
			return true
		}
	}
	return false
}

func unicodeLocale(environment map[string]string) bool {
	locale := environment["LC_ALL"]
	if locale == "" {
		locale = environment["LC_CTYPE"]
	}
	if locale == "" {
		locale = environment["LANG"]
	}
	locale = strings.ToLower(locale)
	return locale == "" || locale != "c" && locale != "posix"
}
