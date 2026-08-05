// Package cli composes and runs the Aruo command-line application.
package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aruodore/aruo/internal/buildinfo"
	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/cli/command"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
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
}

// Run executes one CLI invocation and returns a process exit code.
func Run(ctx context.Context, args []string, dependencies Dependencies) int {
	root := command.NewRoot(dependencies.Streams, dependencies.Build, dependencies.Catalog, dependencies.Creator, dependencies.Doctor)
	root.SetArgs(args)

	if err := root.ExecuteContext(ctx); err != nil {
		if dependencies.Logger != nil {
			dependencies.Logger.DebugContext(ctx, "command failed", "error", err)
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
