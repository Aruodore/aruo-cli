// Package builtin constructs Aruo's compiled, qualified project catalog.
package builtin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/templateengine"
	templatebuiltin "github.com/aruodore/aruo/internal/templateengine/builtin"
)

var nonIdentifier = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

// New returns the compiled first-party catalog.
func New() (*catalog.Memory, error) {
	goSource, goBlueprint := templatebuiltin.GoLibrary()
	jsSource, jsBlueprint := templatebuiltin.JSLibrary()
	tsSource, tsBlueprint := templatebuiltin.TSLibrary()
	pySource, pyBlueprint := templatebuiltin.PythonLibrary()
	reactSource, reactBlueprint := templatebuiltin.ReactApp()
	nuxtSource, nuxtBlueprint := templatebuiltin.NuxtApp()
	vueSource, vueBlueprint := templatebuiltin.VueLibrary()
	return catalog.NewMemory(
		catalog.Entry{
			ID:             "go-library",
			Name:           "Go library",
			Language:       "go",
			Kind:           "library",
			Description:    "A production-ready Go library with CI, governance, security, tests, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         goSource,
			Blueprint:      goBlueprint,
			Defaults: map[string]any{
				"GoVersion":      "1.26.0",
				"IncludeInstall": true,
				"TemplateID":     "go-library",
			},
			Prepare: prepareGoLibrary,
			NextSteps: []string{
				"go test ./...",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "Go module path",
				ModuleDescription: "This is written to go.mod and becomes the import path for your library.",
				ModuleExample:     "github.com/your-name/my-library",
			},
		},
		catalog.Entry{
			ID:             "js-library",
			Name:           "JavaScript library",
			Language:       "javascript",
			Kind:           "library",
			Description:    "A production-ready JavaScript library with CI, governance, security, tests, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         jsSource,
			Blueprint:      jsBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": true,
				"TemplateID":     "js-library",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"node --test",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name".`,
				ModuleExample:     "my-library",
			},
		},
		catalog.Entry{
			ID:             "ts-library",
			Name:           "TypeScript library",
			Language:       "typescript",
			Kind:           "library",
			Description:    "A production-ready TypeScript library with CI, governance, security, tests, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         tsSource,
			Blueprint:      tsBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": true,
				"TemplateID":     "ts-library",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"npm install",
				"npm test",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name".`,
				ModuleExample:     "my-library",
			},
		},
		catalog.Entry{
			ID:             "react-app",
			Name:           "React application",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A production-ready React application with Vite, Vitest, CI, governance, security, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         reactSource,
			Blueprint:      reactBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "react-app",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"npm install",
				"npm run dev",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		catalog.Entry{
			ID:             "nuxt-app",
			Name:           "Nuxt application",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A production-ready Nuxt application with server-side rendering, CI, governance, security, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         nuxtSource,
			Blueprint:      nuxtBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "nuxt-app",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"npm install",
				"npm run dev",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		catalog.Entry{
			ID:             "vue-library",
			Name:           "Vue library",
			Language:       "typescript",
			Kind:           "library",
			Description:    "A production-ready Vue component library with CI, governance, security, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         vueSource,
			Blueprint:      vueBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": true,
				"TemplateID":     "vue-library",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"npm install",
				"npm test",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name".`,
				ModuleExample:     "my-library",
			},
		},
		catalog.Entry{
			ID:             "python-library",
			Name:           "Python library",
			Language:       "python",
			Kind:           "library",
			Description:    "A production-ready Python library with CI, governance, security, tests, and documentation",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         pySource,
			Blueprint:      pyBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": true,
				"TemplateID":     "python-library",
			},
			Prepare: preparePythonLibrary,
			NextSteps: []string{
				"make check",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "PyPI package name",
				ModuleDescription: "This is written to pyproject.toml as \"name\".",
				ModuleExample:     "my-library",
			},
		},
	)
}

func prepareGoLibrary(data templateengine.Data) (templateengine.Data, error) {
	if data.Project.Module == "" {
		return data, fmt.Errorf("go module path is required (use --module)")
	}
	name := strings.ToLower(nonIdentifier.ReplaceAllString(data.Project.Name, "_"))
	name = strings.Trim(name, "_")
	if name == "" || (name[0] >= '0' && name[0] <= '9') {
		return data, fmt.Errorf("project name %q cannot form a Go package identifier", data.Project.Name)
	}
	data.Variables["PackageName"] = name
	return data, nil
}

func prepareJSLibrary(data templateengine.Data) (templateengine.Data, error) {
	if data.Project.Module == "" {
		return data, fmt.Errorf("npm package name is required (use --module)")
	}
	if data.Project.Module != strings.ToLower(data.Project.Module) {
		return data, fmt.Errorf("npm package name %q must be lowercase", data.Project.Module)
	}
	return data, nil
}

func preparePythonLibrary(data templateengine.Data) (templateengine.Data, error) {
	if data.Project.Module == "" {
		return data, fmt.Errorf("PyPI package name is required (use --module)")
	}
	name := strings.ToLower(nonIdentifier.ReplaceAllString(data.Project.Name, "_"))
	name = strings.Trim(name, "_")
	if name == "" || (name[0] >= '0' && name[0] <= '9') {
		return data, fmt.Errorf("project name %q cannot form a Python package identifier", data.Project.Name)
	}
	data.Variables["PackageName"] = name
	return data, nil
}
