// Package builtin constructs Aruo's compiled, qualified project catalog.
package builtin

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/templateengine"
	templatebuiltin "github.com/aruodore/aruo-cli/internal/templateengine/builtin"
)

var nonIdentifier = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

// New returns the compiled first-party catalog.
func New() (*catalog.Memory, error) {
	goSource, goBlueprint := templatebuiltin.GoLibrary()
	jsSource, jsBlueprint := templatebuiltin.JSLibrary()
	tsSource, tsBlueprint := templatebuiltin.TSLibrary()
	pySource, pyBlueprint := templatebuiltin.PythonLibrary()
	reactSource, reactBlueprint := templatebuiltin.ReactApp()
	reactLeanSource, reactLeanBlueprint := templatebuiltin.ReactLean()
	nuxtSource, nuxtBlueprint := templatebuiltin.NuxtApp()
	nuxtLeanSource, nuxtLeanBlueprint := templatebuiltin.NuxtLean()
	vueSource, vueBlueprint := templatebuiltin.VueLibrary()
	vueAppSource, vueAppBlueprint := templatebuiltin.VueApp()
	vueLeanSource, vueLeanBlueprint := templatebuiltin.VueLean()
	nextSource, nextBlueprint := templatebuiltin.NextApp()
	nextLeanSource, nextLeanBlueprint := templatebuiltin.NextLean()
	return catalog.NewMemory(
		catalog.Entry{
			ID:             "go-library",
			Name:           "Go library",
			Language:       "go",
			Kind:           "library",
			Description:    "A production-ready Go library with zero external dependencies, tests, CI, governance, security, and documentation",
			Color:          "#39C6E8",
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
			Description:    "A production-ready JavaScript library with zero npm dependencies, tests, CI, governance, security, and documentation",
			Color:          "#F7DF1E",
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
			Description:    "A production-ready TypeScript library with strict type-checking, tests, CI, governance, security, and documentation",
			Color:          "#65AADD",
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
			ID:             "react",
			Name:           "React",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A domain-neutral React modular monolith with a Node server, PostgreSQL, production controls, tests, and explicit limitations",
			Color:          "#61DAFB",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         reactSource,
			Blueprint:      reactBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "react",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"cp .env.example .env && npm install",
				"docker compose up -d db && npm run db:migrate",
				"npm run dev",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		webAppEntry("react-lean", "React (lean)", "A lean React baseline with explicit production omissions", "#61DAFB", reactLeanSource, reactLeanBlueprint),
		catalog.Entry{
			ID:             "nuxt",
			Name:           "Nuxt",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A Nuxt modular monolith with PostgreSQL, validation, structured logging, health checks, containers, tests, and an AI development contract",
			Color:          "#38E59D",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         nuxtSource,
			Blueprint:      nuxtBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "nuxt",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"cp .env.example .env && npm install",
				"docker compose up -d db && npm run db:migrate",
				"npm run dev",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		webAppEntry("nuxt-lean", "Nuxt (lean)", "A lean Nuxt baseline with explicit production omissions", "#38E59D", nuxtLeanSource, nuxtLeanBlueprint),
		catalog.Entry{
			ID:             "vue-library",
			Name:           "Vue library",
			Language:       "typescript",
			Kind:           "library",
			Description:    "A production-ready Vue 3 component library built with Vite library mode, tests, CI, governance, security, and documentation",
			Color:          "#63D6AA",
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
			ID:             "vue",
			Name:           "Vue",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A domain-neutral Vue modular monolith with a Node server, PostgreSQL, production controls, tests, and explicit limitations",
			Color:          "#63D6AA",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         vueAppSource,
			Blueprint:      vueAppBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "vue",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"cp .env.example .env && npm install",
				"docker compose up -d db && npm run db:migrate",
				"npm run dev",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		webAppEntry("vue-lean", "Vue (lean)", "A lean Vue baseline with explicit production omissions", "#63D6AA", vueLeanSource, vueLeanBlueprint),
		catalog.Entry{
			ID:             "next",
			Name:           "Next.js",
			Language:       "typescript",
			Kind:           "app",
			Description:    "A domain-neutral Next.js modular monolith with PostgreSQL, production controls, tests, and explicit limitations",
			Color:          "#B79AF4",
			Licenses:       []string{"MIT"},
			DefaultLicense: "MIT",
			Source:         nextSource,
			Blueprint:      nextBlueprint,
			Defaults: map[string]any{
				"IncludeInstall": false,
				"TemplateID":     "next",
			},
			Prepare: prepareJSLibrary,
			NextSteps: []string{
				"cp .env.example .env && npm install",
				"docker compose up -d db && npm run db:migrate",
				"npm run dev",
				"git init && git add .",
			},
			Prompts: catalog.ProjectPrompts{
				ModuleLabel:       "npm package name",
				ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`,
				ModuleExample:     "my-app",
			},
		},
		webAppEntry("next-lean", "Next.js (lean)", "A lean Next.js baseline with explicit production omissions", "#B79AF4", nextLeanSource, nextLeanBlueprint),
		catalog.Entry{
			ID:             "python-library",
			Name:           "Python library",
			Language:       "python",
			Kind:           "library",
			Description:    "A production-ready Python library with zero pip dependencies, tests, CI, governance, security, and documentation",
			Color:          "#FFD43B",
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

func webAppEntry(id, name, description, color string, source fs.FS, blueprint templateengine.Blueprint) catalog.Entry {
	return catalog.Entry{
		ID: id, Name: name, Language: "typescript", Kind: "app", Description: description, Color: color,
		Licenses: []string{"MIT"}, DefaultLicense: "MIT", Source: source, Blueprint: blueprint,
		Defaults:  map[string]any{"IncludeInstall": false, "TemplateID": id},
		Prepare:   prepareJSLibrary,
		NextSteps: []string{"npm install", "npm run dev", "git init && git add ."},
		Prompts: catalog.ProjectPrompts{
			ModuleLabel: "npm package name", ModuleDescription: `This is written to package.json as "name" (the app is private and not published).`, ModuleExample: "my-app",
		},
	}
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
