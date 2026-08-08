// Package cli composes and runs the Aruo command-line application.
package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/cli/command"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/doctor"
	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/aruodore/aruo-cli/internal/tux/term"
)

// Dependencies contains process-level resources. Command constructors receive
// narrower values from this composition root instead of reading globals.
type Dependencies struct {
	Build   buildinfo.Info
	Catalog catalog.Catalog
	Creator *create.Service
	Doctor  *doctor.Service
	Logger  *slog.Logger
	Streams iostreams.IOStreams
	// Environment supplies presentation-policy hints (NO_COLOR, CI, TERM, and
	// similar). A nil map resolves policy with no environment hints, which
	// keeps tests hermetic; cmd/aruo supplies the real process environment.
	Environment map[string]string
	// Probe overrides terminal detection for tests. A nil probe uses the
	// system terminal.
	Probe term.Probe
}

// Run executes one CLI invocation and returns a process exit code.
func Run(ctx context.Context, args []string, dependencies Dependencies) int {
	environment := dependencies.Environment
	if environment == nil {
		environment = map[string]string{}
	}
	root := command.NewRoot(dependencies.Streams, dependencies.Build, dependencies.Catalog, dependencies.Creator, dependencies.Doctor, environment, dependencies.Probe)
	root.SetArgs(args)

	if err := root.ExecuteContext(ctx); err != nil {
		if dependencies.Logger != nil {
			dependencies.Logger.DebugContext(ctx, "command failed", "error", err)
		}
		if isCancellation(err) {
			_, _ = fmt.Fprintln(dependencies.Streams.ErrOut, "Cancelled.")
			return 130
		}
		var silent interface{ SuppressMessage() bool }
		if !errors.As(err, &silent) || !silent.SuppressMessage() {
			_, _ = fmt.Fprintf(dependencies.Streams.ErrOut, "Error: %v\n", err)
		}
		var exitCoder interface{ ExitCode() int }
		if errors.As(err, &exitCoder) {
			return exitCoder.ExitCode()
		}
		return 1
	}

	return 0
}

// isCancellation reports whether err stems from Ctrl+C or SIGTERM rather than
// an operational failure. Interactive rich adapters own raw terminal mode
// while running, which keeps the terminal driver from generating SIGINT, so
// cancellation can surface as context.Canceled (the process-level signal
// path) or tux.ErrCancelled (an adapter's own Ctrl+C key handling) without
// the other.
func isCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, tux.ErrCancelled)
}
