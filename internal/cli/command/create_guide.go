package command

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/tux"
)

// runGuide drives the entire interactive create flow -- project name,
// kind, template, module, description, author, and the final confirmation
// -- as one Prompter.Guide session, so the user can move backward across
// any of those screens rather than only within one. Each step is skipped
// when the corresponding flag already answers it, matching what the
// non-interactive path (resolveTemplate/resolveProjectFields) already
// requires via flags instead of prompts.
func runGuide(ctx context.Context, prompter tux.Prompter, templateCatalog catalog.Catalog, options *createOptions) (catalog.Entry, error) {
	entries, err := templateCatalog.List(ctx)
	if err != nil {
		return catalog.Entry{}, err
	}
	// Only offer kinds that actually have a template under the requested
	// --language filter, so a user can never pick a kind that narrows the
	// template step to zero options.
	languageFiltered := filterEntries(entries, options.language, "")
	kinds := uniqueValues(languageFiltered, func(entry catalog.Entry) string { return entry.Kind })
	gitAuthor := detectGitAuthor(ctx)

	resolveEntry := func(answers tux.Answers) catalog.Entry {
		id := options.templateID
		if value, ok := answers["template"].(tux.OptionID); ok && value != "" {
			id = string(value)
		}
		// id is always a real catalog ID here: either a pre-validated
		// --template flag, or a value the template step itself offered
		// from this same catalog. Resolve isn't expected to fail; a
		// closure has no error return to propagate one anyway.
		entry, _ := templateCatalog.Resolve(ctx, id)
		return entry
	}
	effectiveDestination := func(answers tux.Answers) string {
		if options.destination != "" {
			return options.destination
		}
		if value, ok := answers["name"].(string); ok {
			return value
		}
		return ""
	}

	steps := []tux.Step{
		{
			ID:   "name",
			Kind: tux.StepInput,
			Skip: func() bool { return options.destination != "" },
			Input: func(tux.Answers) tux.InputRequest {
				return tux.InputRequest{
					ID:          "name",
					Label:       "Project name",
					Description: "Aruo will create a new folder with this name.",
					Example:     "my-library",
				}
			},
		},
		{
			ID:   "kind",
			Kind: tux.StepSelect,
			Skip: func() bool { return options.templateID != "" || options.kind != "" || len(kinds) <= 1 },
			Select: func(tux.Answers) tux.SelectRequest {
				kindOptions := make([]tux.Option, len(kinds))
				for index, kind := range kinds {
					kindOptions[index] = tux.Option{ID: tux.OptionID(kind), Label: kindLabel(kind), Description: kindEntryNames(languageFiltered, kind)}
				}
				return tux.SelectRequest{ID: "kind", Label: "What are you building?", Options: kindOptions}
			},
		},
		{
			ID:   "template",
			Kind: tux.StepSelect,
			Skip: func() bool { return options.templateID != "" },
			Select: func(answers tux.Answers) tux.SelectRequest {
				kind := options.kind
				if value, ok := answers["kind"].(tux.OptionID); ok && value != "" {
					kind = string(value)
				}
				matching := filterEntries(languageFiltered, "", kind)
				return tux.SelectRequest{ID: "template", Label: "Template", Options: templateOptions(matching)}
			},
		},
		{
			ID:   "module",
			Kind: tux.StepInput,
			Skip: func() bool { return options.module != "" },
			Input: func(answers tux.Answers) tux.InputRequest {
				entry := resolveEntry(answers)
				return tux.InputRequest{
					ID:          "module",
					Label:       moduleLabel(entry),
					Description: fmt.Sprintf("Creating: %s\n\n%s", entry.Name, entry.Prompts.ModuleDescription),
					Example:     entry.Prompts.ModuleExample,
					Placeholder: entry.Prompts.ModuleExample,
				}
			},
		},
		{
			ID:   "description",
			Kind: tux.StepInput,
			Skip: func() bool { return options.description != "" },
			Input: func(answers tux.Answers) tux.InputRequest {
				entry := resolveEntry(answers)
				return tux.InputRequest{
					ID:          "description",
					Label:       "Short description (Optional)",
					Description: "One sentence explaining what the project does.",
					Example:     "A Go library for reliable configuration loading",
					Placeholder: "A Go library for reliable configuration loading",
					Optional:    true,
					Default:     &entry.Description,
				}
			},
		},
		{
			ID:   "author",
			Kind: tux.StepInput,
			Skip: func() bool { return options.author != "" },
			Input: func(tux.Answers) tux.InputRequest {
				return tux.InputRequest{
					ID:          "author",
					Label:       "Author or organization (Optional)",
					Description: "Used in the license and project metadata.",
					Example:     "Jane Doe or Acme, Inc.",
					Placeholder: "Jane Doe or Acme, Inc.",
					Optional:    true,
					Default:     &gitAuthor,
				}
			},
		},
		{
			ID:   "confirm",
			Kind: tux.StepConfirm,
			Skip: func() bool { return options.yes },
			Confirm: func(answers tux.Answers) tux.ConfirmRequest {
				entry := resolveEntry(answers)
				name := options.name
				if value, ok := answers["name"].(string); ok && value != "" {
					name = value
				}
				destination := effectiveDestination(answers)
				if name == "" {
					name = filepath.Base(filepath.Clean(destination))
				}
				module := options.module
				if value, ok := answers["module"].(string); ok {
					module = value
				}
				license := options.license
				if license == "" {
					license = entry.DefaultLicense
				}
				summary := fmt.Sprintf(
					"Project summary:\n  Name:        %s\n  Template:    %s\n  Location:    %s\n  %s: %s\n  License:     %s",
					name, entry.Name, destination, moduleLabel(entry), module, license,
				)
				return tux.ConfirmRequest{ID: "confirm", Label: "Create this project?", Description: summary}
			},
		},
	}

	answers, err := prompter.Guide(ctx, steps)
	if err != nil {
		return catalog.Entry{}, err
	}

	if options.destination == "" {
		if value, ok := answers["name"].(string); ok {
			options.name = value
			options.destination = value
		}
	}
	if options.templateID == "" {
		if value, ok := answers["template"].(tux.OptionID); ok {
			options.templateID = string(value)
		}
	}
	if options.module == "" {
		if value, ok := answers["module"].(string); ok {
			options.module = value
		}
	}
	if options.description == "" {
		if value, ok := answers["description"].(string); ok {
			options.description = value
		}
	}
	if options.author == "" {
		if value, ok := answers["author"].(string); ok {
			options.author = value
		}
	}
	if !options.yes {
		if confirmed, ok := answers["confirm"].(bool); ok && !confirmed {
			return catalog.Entry{}, errors.New("creation cancelled")
		}
	}

	entry, err := templateCatalog.Resolve(ctx, options.templateID)
	if err != nil {
		return catalog.Entry{}, err
	}
	return entry, nil
}

func kindLabel(kind string) string {
	switch kind {
	case "app":
		return "Application"
	case "library":
		return "Library"
	default:
		return kind
	}
}

func kindEntryNames(entries []catalog.Entry, kind string) string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Kind == kind {
			names = append(names, entry.Name)
		}
	}
	return strings.Join(names, ", ")
}
