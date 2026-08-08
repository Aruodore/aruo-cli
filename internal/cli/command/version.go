package command

import (
	"fmt"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	"github.com/spf13/cobra"
)

func newVersion(build buildinfo.Info) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the Aruo version",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(command.OutOrStdout(), "aruo version %s\n", build.Version)
			return err
		},
	}
}
