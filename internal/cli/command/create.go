package command

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/templateengine"
	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/spf13/cobra"
)

type createOptions struct {
	destination    string
	templateID     string
	language       string
	kind           string
	name           string
	module         string
	description    string
	author         string
	license        string
	nonInteractive bool
	yes            bool
	variables      []string
}

func newCreate(factory sessionFactory, templateCatalog catalog.Catalog, creator *create.Service) *cobra.Command {
	options := createOptions{}
	command := &cobra.Command{
		Use:   "create [name-or-path]",
		Short: "Create a project from a built-in template",
		Example: `  aruo create my-library
  aruo create tools/my-library
  aruo create my-library --template go-library`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 1 {
				options.destination = args[0]
			}
			return runCreate(command.Context(), factory, templateCatalog, creator, options)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.templateID, "template", "", "template ID (for example, go-library)")
	flags.StringVar(&options.language, "language", "", "filter templates by language")
	flags.StringVar(&options.kind, "kind", "", "filter templates by project kind")
	flags.StringVar(&options.name, "name", "", "project name (normally inferred; useful when creating in '.')")
	flags.StringVar(&options.module, "module", "", "package or module identifier for the template's ecosystem (for example, a Go module path or an npm package name)")
	flags.StringVar(&options.description, "description", "", "one-line project description")
	flags.StringVar(&options.author, "author", "", "copyright holder and security contact")
	flags.StringVar(&options.license, "license", "", "SPDX license identifier (defaults to the template recommendation)")
	flags.StringSliceVar(&options.variables, "set", nil, "template variable as key=value (repeatable)")
	flags.BoolVar(&options.nonInteractive, "non-interactive", false, "disable prompts and fail on missing input")
	flags.BoolVarP(&options.yes, "yes", "y", false, "accept the final creation confirmation")
	_ = flags.MarkDeprecated("non-interactive", "use --no-input instead")
	registerCreateCompletions(command, templateCatalog)
	return command
}

// registerCreateCompletions offers template, language, and kind values from
// the local, offline template catalog. Shell completion must return quickly
// and never reach the network; the builtin catalog is embedded, so listing
// it is instant.
func registerCreateCompletions(command *cobra.Command, templateCatalog catalog.Catalog) {
	entries := func() ([]catalog.Entry, error) { return templateCatalog.List(command.Context()) }
	_ = command.RegisterFlagCompletionFunc("template", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		list, err := entries()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		completions := make([]string, 0, len(list))
		for _, entry := range list {
			completions = append(completions, entry.ID+"\t"+entry.Description)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})
	_ = command.RegisterFlagCompletionFunc("language", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		list, err := entries()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return uniqueValues(list, func(entry catalog.Entry) string { return entry.Language }), cobra.ShellCompDirectiveNoFileComp
	})
	_ = command.RegisterFlagCompletionFunc("kind", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		list, err := entries()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return uniqueValues(list, func(entry catalog.Entry) string { return entry.Kind }), cobra.ShellCompDirectiveNoFileComp
	})
}

func uniqueValues(entries []catalog.Entry, field func(catalog.Entry) string) []string {
	seen := make(map[string]struct{}, len(entries))
	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		value := field(entry)
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func runCreate(ctx context.Context, factory sessionFactory, templateCatalog catalog.Catalog, creator *create.Service, options createOptions) error {
	terminal, err := factory.build(ctx, options.nonInteractive)
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Close() }()
	prompter := terminal.Prompter()
	interactive := terminal.Policy().Mode == tux.ModeInteractive

	var entry catalog.Entry
	if interactive {
		entry, err = runGuide(ctx, prompter, templateCatalog, &options)
	} else {
		if options.destination == "" {
			options.name, err = resolveInput(ctx, prompter, options.name, tux.InputRequest{
				ID:          "name",
				Label:       "Project name",
				Description: "Aruo will create a new folder with this name.",
				Example:     "my-library",
			}, "name-or-path argument")
			if err != nil {
				return err
			}
			options.destination = options.name
		}
		entry, err = resolveTemplate(ctx, templateCatalog, &options)
		if err == nil {
			err = resolveProjectFields(ctx, prompter, entry, &options)
		}
	}
	if err != nil {
		return err
	}

	if options.name == "" {
		options.name = filepath.Base(filepath.Clean(options.destination))
	}
	if options.name == "." {
		return errors.New("cannot infer a project name for the current directory; provide --name")
	}
	if options.license == "" {
		options.license = entry.DefaultLicense
	}
	if !contains(entry.Licenses, options.license) {
		return fmt.Errorf("template %q does not support license %q; supported: %s", entry.ID, options.license, strings.Join(entry.Licenses, ", "))
	}
	variables, err := parseVariables(options.variables)
	if err != nil {
		return err
	}

	result, err := creator.Create(ctx, create.Request{
		Destination: options.destination,
		TemplateID:  entry.ID,
		Project: templateengine.Project{
			Name: options.name, Module: options.module, Description: options.description,
			Author: options.author, License: options.license, Language: entry.Language,
		},
		Variables: variables,
	})
	if err != nil {
		return err
	}
	return presentCreated(ctx, terminal.Presenter(), options.destination, result)
}

// resolveTemplate applies options.language/options.kind filters and
// resolves options.templateID for a non-interactive session, where an
// empty or ambiguous match can only be resolved through the --template
// flag -- there's no prompter to fall back on. The interactive session
// resolves the same fields through runGuide's "kind" and "template" steps
// instead.
func resolveTemplate(ctx context.Context, templateCatalog catalog.Catalog, options *createOptions) (catalog.Entry, error) {
	entries, err := templateCatalog.List(ctx)
	if err != nil {
		return catalog.Entry{}, err
	}
	entries = filterEntries(entries, options.language, options.kind)
	if options.templateID == "" {
		switch len(entries) {
		case 0:
			return catalog.Entry{}, errors.New("no templates match the requested language/kind filters")
		case 1:
			options.templateID = entries[0].ID
		default:
			return catalog.Entry{}, errors.New("--template is required when multiple templates match")
		}
	}
	entry, err := templateCatalog.Resolve(ctx, options.templateID)
	if err != nil {
		return catalog.Entry{}, err
	}
	if options.language != "" && entry.Language != options.language || options.kind != "" && entry.Kind != options.kind {
		return catalog.Entry{}, fmt.Errorf("template %q does not match the requested language/kind filters", entry.ID)
	}
	return entry, nil
}

// resolveProjectFields fills in the module, description, and author
// fields for a non-interactive session. --module never prompts even
// interactively; it silently defaults to the project's own name, since
// asking for it separately is redundant -- pass --module explicitly for a
// real Go import path like github.com/you/name. description/author use
// resolveInput, which requires each via its flag unless the field is
// optional, in which case a real default is derived automatically
// (detectGitAuthor, catalog.Entry.Description).
func resolveProjectFields(ctx context.Context, prompter tux.Prompter, entry catalog.Entry, options *createOptions) error {
	if options.module == "" {
		name := options.name
		if name == "" {
			name = filepath.Base(filepath.Clean(options.destination))
		}
		options.module = name
	}
	var err error
	options.description, err = resolveInput(ctx, prompter, options.description, tux.InputRequest{
		ID:          "description",
		Label:       "Short description (Optional)",
		Description: "One sentence explaining what the project does.",
		Example:     "A Go library for reliable configuration loading",
		Placeholder: "A Go library for reliable configuration loading",
		Optional:    true,
		Default:     &entry.Description,
	}, "--description")
	if err != nil {
		return err
	}
	gitAuthor := detectGitAuthor(ctx)
	options.author, err = resolveInput(ctx, prompter, options.author, tux.InputRequest{
		ID:          "author",
		Label:       "Author or organization (Optional)",
		Description: "Used in the license and project metadata.",
		Example:     "Jane Doe or Acme, Inc.",
		Placeholder: "Jane Doe or Acme, Inc.",
		Optional:    true,
		Default:     &gitAuthor,
	}, "--author")
	return err
}

// detectGitAuthor best-effort reads the local `git config user.name` so
// --author defaults to something real instead of a placeholder that needs
// manual follow-up. Any failure (git missing, unset, timeout) yields an
// empty default rather than an error; the caller treats that the same as
// no default at all.
func detectGitAuthor(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	output, err := exec.CommandContext(ctx, "git", "config", "--get", "user.name").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func moduleLabel(entry catalog.Entry) string {
	if entry.Prompts.ModuleLabel != "" {
		return entry.Prompts.ModuleLabel
	}
	return "Package or module path"
}

func presentCreated(ctx context.Context, presenter tux.Presenter, destination string, result create.Result) error {
	lines := []string{fmt.Sprintf("Created %s with %d files at %s", result.TemplateID, result.FileCount, result.Destination)}
	if len(result.NextSteps) > 0 {
		lines = append(lines, "", "Next steps:", "  cd "+destination)
		for _, step := range result.NextSteps {
			lines = append(lines, "  "+step)
		}
	}
	return presenter.Message(ctx, tux.Message{Kind: tux.MessageSuccess, Text: strings.Join(lines, "\n")})
}

func resolveInput(ctx context.Context, prompter tux.Prompter, value string, request tux.InputRequest, flag string) (string, error) {
	if value != "" {
		return value, nil
	}
	resolved, err := prompter.Input(ctx, request)
	switch {
	case err == nil:
		return resolved, nil
	case errors.Is(err, tux.ErrUnavailable):
		if request.Optional {
			if request.Default != nil {
				return *request.Default, nil
			}
			return "", nil
		}
		return "", fmt.Errorf("%s is required; provide %s", request.Label, flag)
	default:
		return "", err
	}
}

func templateOptions(entries []catalog.Entry) []tux.Option {
	options := make([]tux.Option, len(entries))
	for index, entry := range entries {
		options[index] = tux.Option{ID: tux.OptionID(entry.ID), Label: entry.Name, Color: entry.Color}
	}
	return options
}

func filterEntries(entries []catalog.Entry, language, kind string) []catalog.Entry {
	result := make([]catalog.Entry, 0, len(entries))
	for _, entry := range entries {
		if language != "" && entry.Language != language || kind != "" && entry.Kind != kind {
			continue
		}
		result = append(result, entry)
	}
	return result
}

func parseVariables(values []string) (map[string]any, error) {
	result := make(map[string]any, len(values))
	for _, value := range values {
		key, raw, found := strings.Cut(value, "=")
		if !found || key == "" {
			return nil, fmt.Errorf("invalid --set %q; expected key=value", value)
		}
		if boolean, err := strconv.ParseBool(raw); err == nil {
			result[key] = boolean
		} else {
			result[key] = raw
		}
	}
	return result, nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
