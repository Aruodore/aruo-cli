package command

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/templateengine"
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

func newCreate(streams iostreams.IOStreams, templateCatalog catalog.Catalog, creator *create.Service) *cobra.Command {
	options := createOptions{}
	command := &cobra.Command{
		Use:   "create [name-or-path]",
		Short: "Create a production-ready project",
		Example: `  aruo create my-library
  aruo create tools/my-library
  aruo create my-library --template go-library`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 1 {
				options.destination = args[0]
			}
			return runCreate(command.Context(), streams, templateCatalog, creator, options)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.templateID, "template", "", "template ID (for example, go-library)")
	flags.StringVar(&options.language, "language", "", "filter templates by language")
	flags.StringVar(&options.kind, "kind", "", "filter templates by project kind")
	flags.StringVar(&options.name, "name", "", "project name (normally inferred; useful when creating in '.')")
	flags.StringVar(&options.module, "module", "", "Go module path for Go templates (for example, github.com/you/project)")
	flags.StringVar(&options.description, "description", "", "one-line project description")
	flags.StringVar(&options.author, "author", "", "copyright holder and security contact")
	flags.StringVar(&options.license, "license", "", "SPDX license identifier (defaults to the template recommendation)")
	flags.StringSliceVar(&options.variables, "set", nil, "template variable as key=value (repeatable)")
	flags.BoolVar(&options.nonInteractive, "non-interactive", false, "disable prompts and fail on missing input")
	flags.BoolVarP(&options.yes, "yes", "y", false, "accept the final creation confirmation")
	return command
}

func runCreate(ctx context.Context, streams iostreams.IOStreams, templateCatalog catalog.Catalog, creator *create.Service, options createOptions) error {
	prompter := &linePrompter{scanner: bufio.NewScanner(streams.In), out: streams.ErrOut, disabled: options.nonInteractive}
	var err error
	if options.destination == "" {
		options.name, err = prompter.requiredWithHint("Project name", options.name, "", "Aruo will create a new folder with this name.", "my-library")
		if err != nil {
			return missingFlag(err, "name-or-path argument")
		}
		options.destination = options.name
	}
	entries, err := templateCatalog.List(ctx)
	if err != nil {
		return err
	}
	entries = filterEntries(entries, options.language, options.kind)
	if options.templateID == "" {
		switch {
		case len(entries) == 0:
			return errors.New("no templates match the requested language/kind filters")
		case len(entries) == 1:
			options.templateID = entries[0].ID
		case options.nonInteractive:
			return errors.New("--template is required when multiple templates match")
		default:
			for _, entry := range entries {
				_, _ = fmt.Fprintf(streams.ErrOut, "  %s — %s\n", entry.ID, entry.Description)
			}
			options.templateID, err = prompter.required("Template", "", "")
			if err != nil {
				return err
			}
		}
	}
	entry, err := templateCatalog.Resolve(ctx, options.templateID)
	if err != nil {
		return err
	}
	if options.language != "" && entry.Language != options.language || options.kind != "" && entry.Kind != options.kind {
		return fmt.Errorf("template %q does not match the requested language/kind filters", entry.ID)
	}
	if !options.nonInteractive {
		_, _ = fmt.Fprintf(streams.ErrOut, "\nCreating: %s\n%s\n\n", entry.Name, entry.Description)
	}
	if options.name == "" {
		options.name = filepath.Base(filepath.Clean(options.destination))
	}
	if options.name == "" {
		options.name = filepath.Base(filepath.Clean(options.destination))
	}
	if options.name == "." {
		return errors.New("cannot infer a project name for the current directory; provide --name")
	}
	moduleLabel := entry.Prompts.ModuleLabel
	if moduleLabel == "" {
		moduleLabel = "Package or module path"
	}
	options.module, err = prompter.requiredWithHint(moduleLabel, options.module, "", entry.Prompts.ModuleDescription, entry.Prompts.ModuleExample)
	if err != nil {
		return missingFlag(err, "--module")
	}
	options.description, err = prompter.requiredWithHint("Short description", options.description, "", "One sentence explaining what the project does.", "A Go library for reliable configuration loading")
	if err != nil {
		return missingFlag(err, "--description")
	}
	options.author, err = prompter.requiredWithHint("Author or organization", options.author, "", "Used in the license and project metadata.", "Jane Doe or Acme, Inc.")
	if err != nil {
		return missingFlag(err, "--author")
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
	if !options.yes && !options.nonInteractive {
		_, _ = fmt.Fprintln(streams.ErrOut, "\nProject summary:")
		_, _ = fmt.Fprintf(streams.ErrOut, "  Name:        %s\n", options.name)
		_, _ = fmt.Fprintf(streams.ErrOut, "  Template:    %s\n", entry.Name)
		_, _ = fmt.Fprintf(streams.ErrOut, "  Location:    %s\n", options.destination)
		_, _ = fmt.Fprintf(streams.ErrOut, "  %s: %s\n", moduleLabel, options.module)
		_, _ = fmt.Fprintf(streams.ErrOut, "  License:     %s\n\n", options.license)
		confirmed, confirmErr := prompter.confirm("Create this project?")
		if confirmErr != nil {
			return confirmErr
		}
		if !confirmed {
			return errors.New("creation cancelled")
		}
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
	_, _ = fmt.Fprintf(streams.Out, "Created %s with %d files at %s\n", result.TemplateID, result.FileCount, result.Destination)
	_, _ = fmt.Fprintln(streams.Out, "\nNext steps:")
	_, _ = fmt.Fprintf(streams.Out, "  cd %s\n", options.destination)
	for _, step := range result.NextSteps {
		_, _ = fmt.Fprintf(streams.Out, "  %s\n", step)
	}
	return nil
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

type linePrompter struct {
	scanner  *bufio.Scanner
	out      interface{ Write([]byte) (int, error) }
	disabled bool
}

func (p *linePrompter) required(label, value, defaultValue string) (string, error) {
	return p.requiredWithHint(label, value, defaultValue, "", "")
}

func (p *linePrompter) requiredWithHint(label, value, defaultValue, description, example string) (string, error) {
	if value != "" {
		return value, nil
	}
	if p.disabled {
		return "", fmt.Errorf("%s is required in non-interactive mode", label)
	}
	if description != "" {
		_, _ = fmt.Fprintf(p.out, "%s\n", description)
	}
	if example != "" {
		_, _ = fmt.Fprintf(p.out, "Example: %s\n", example)
	}
	if defaultValue != "" {
		_, _ = fmt.Fprintf(p.out, "%s [%s]: ", label, defaultValue)
	} else {
		_, _ = fmt.Fprintf(p.out, "%s: ", label)
	}
	if !p.scanner.Scan() {
		return "", errors.New("input ended before required answers were provided")
	}
	answer := strings.TrimSpace(p.scanner.Text())
	if answer == "" {
		answer = defaultValue
	}
	if answer == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	return answer, nil
}

func (p *linePrompter) confirm(question string) (bool, error) {
	_, _ = fmt.Fprintf(p.out, "%s [y/N]: ", question)
	if !p.scanner.Scan() {
		return false, errors.New("input ended before confirmation")
	}
	answer := strings.ToLower(strings.TrimSpace(p.scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

func missingFlag(err error, flag string) error {
	return fmt.Errorf("%w; provide %s", err, flag)
}
