// Package command constructs Aruo's command tree.
package command

import (
	"context"

	"github.com/aruodore/aruo/internal/buildinfo"
	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
	"github.com/aruodore/aruo/internal/tux"
	"github.com/aruodore/aruo/internal/tux/session"
	"github.com/aruodore/aruo/internal/tux/term"
	"github.com/spf13/cobra"
)

// sessionFactory builds one terminal session per command invocation using the
// global flags resolved at parse time and the process resources supplied by
// the composition root.
type sessionFactory struct {
	streams     iostreams.IOStreams
	environment map[string]string
	probe       term.Probe
	global      *globalOptions
}

// build resolves the terminal session for one command's output format.
// forceNoInput additionally disables input, satisfying a command's own
// deprecated non-interactive flag without changing global policy.
func (f sessionFactory) build(ctx context.Context, format tux.OutputFormat, forceNoInput bool) (*session.Session, error) {
	overrides, err := f.global.overrides(format)
	if err != nil {
		return nil, err
	}
	if forceNoInput {
		never := tux.FeatureNever
		overrides.Input = &never
	}
	return session.New(ctx, f.streams.In, f.streams.Out, f.streams.ErrOut, f.environment, overrides, f.probe), nil
}

// NewRoot constructs a fresh command tree for each invocation.
func NewRoot(streams iostreams.IOStreams, build buildinfo.Info, templateCatalog catalog.Catalog, creator *create.Service, doctorService *doctor.Service, environment map[string]string, probe term.Probe) *cobra.Command {
	global := newGlobalOptions()
	factory := sessionFactory{streams: streams, environment: environment, probe: probe, global: global}

	root := &cobra.Command{
		Use:           "aruo",
		Short:         "Build and maintain production-quality software projects",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}

	root.SetIn(streams.In)
	root.SetOut(streams.Out)
	root.SetErr(streams.ErrOut)
	// Cobra's generated completion command only reads static command
	// metadata and the registered completion functions below; it never
	// prompts, mutates state, or reaches the network.
	root.CompletionOptions.DisableDefaultCmd = false
	global.register(root)
	root.AddCommand(newVersion(build))
	if templateCatalog != nil && creator != nil {
		root.AddCommand(newCreate(factory, templateCatalog, creator))
	}
	if doctorService != nil {
		root.AddCommand(newDoctor(factory, doctorService))
	}

	return root
}
