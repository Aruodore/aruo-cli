// Package command constructs Aruo's command tree.
package command

import (
	"github.com/aruodore/aruo/internal/buildinfo"
	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
	"github.com/spf13/cobra"
)

// NewRoot constructs a fresh command tree for each invocation.
func NewRoot(streams iostreams.IOStreams, build buildinfo.Info, templateCatalog catalog.Catalog, creator *create.Service, doctorService *doctor.Service) *cobra.Command {
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
	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(newVersion(build))
	if templateCatalog != nil && creator != nil {
		root.AddCommand(newCreate(streams, templateCatalog, creator))
	}
	if doctorService != nil {
		root.AddCommand(newDoctor(streams, doctorService))
	}

	return root
}
