package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aruodore/aruo-cli/internal/initialize"
	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/spf13/cobra"
)

type initOptions struct {
	dryRun bool
	yes    bool
	format string
}

func newInit(factory sessionFactory, service *initialize.Service) *cobra.Command {
	options := initOptions{format: "human"}
	command := &cobra.Command{
		Use:   "init [repository]",
		Short: "Install Aruo's AI engineering contract",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}
			return runInit(command.Context(), factory, service, target, options)
		},
	}
	command.Flags().BoolVar(&options.dryRun, "dry-run", false, "show the initialization plan without writing files")
	command.Flags().BoolVarP(&options.yes, "yes", "y", false, "apply without confirmation")
	command.Flags().StringVar(&options.format, "format", "human", "output format: human or json")
	_ = command.RegisterFlagCompletionFunc("format", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"human", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	return command
}

func runInit(ctx context.Context, factory sessionFactory, service *initialize.Service, target string, options initOptions) error {
	if options.format != "human" && options.format != "json" {
		return fmt.Errorf("unsupported format %q; use human or json", options.format)
	}
	plan, err := service.Plan(ctx, initialize.Request{Target: target})
	if err != nil {
		return err
	}
	if len(plan.Conflicts) > 0 {
		return fmt.Errorf("refusing to overwrite existing paths: %s", strings.Join(plan.Conflicts, ", "))
	}
	if options.dryRun {
		return renderInitResult(ctx, factory, initialize.Result{Repository: plan.Repository, Stack: plan.Stack, Changes: plan.Changes, DryRun: true}, options.format)
	}
	terminal, err := factory.build(ctx, false)
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Close() }()
	if !options.yes {
		confirmed, confirmErr := terminal.Prompter().Confirm(ctx, tux.ConfirmRequest{
			ID:          "confirm-init",
			Label:       "Install Aruo's engineering contract?",
			Description: fmt.Sprintf("Create %d files in %s. Application code and dependencies will not change.", len(plan.Changes), plan.Repository),
			Default:     true,
		})
		if confirmErr != nil {
			return confirmErr
		}
		if !confirmed {
			return errors.New("initialization cancelled")
		}
	}
	result, err := service.Apply(ctx, plan)
	if err != nil {
		return err
	}
	if options.format == "json" {
		return encodeInitJSON(factory, result)
	}
	return renderInitHuman(ctx, terminal.Presenter(), result)
}

func renderInitResult(ctx context.Context, factory sessionFactory, result initialize.Result, format string) error {
	if format == "json" {
		return encodeInitJSON(factory, result)
	}
	terminal, err := factory.build(ctx, true)
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Close() }()
	return renderInitHuman(ctx, terminal.Presenter(), result)
}

func encodeInitJSON(factory sessionFactory, result initialize.Result) error {
	encoder := json.NewEncoder(factory.streams.Out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func renderInitHuman(ctx context.Context, presenter tux.Presenter, result initialize.Result) error {
	frameworks := strings.Join(result.Stack.Frameworks, ", ")
	if frameworks == "" {
		frameworks = "none detected"
	}
	mode := "Installed"
	if result.DryRun {
		mode = "Dry run"
	}
	if err := presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: fmt.Sprintf("%s: Aruo engineering contract\n%s\nStack: %s; frameworks: %s; package manager: %s", mode, result.Repository, result.Stack.Ecosystem, frameworks, result.Stack.PackageManager)}); err != nil {
		return err
	}
	table := tux.Table{Columns: []tux.Column{{ID: "action", Heading: "Action"}, {ID: "path", Heading: "Path"}}}
	for _, change := range result.Changes {
		table.Rows = append(table.Rows, []string{change.Action, change.Path})
	}
	if err := presenter.Table(ctx, table); err != nil {
		return err
	}
	if result.DryRun {
		return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: "No files were written. Run again with --yes to apply this plan."})
	}
	return presenter.Message(ctx, tux.Message{Kind: tux.MessageSuccess, Text: "Aruo initialized. AI agents must begin with AGENTS.md; run aruo doctor to inspect production intent."})
}
