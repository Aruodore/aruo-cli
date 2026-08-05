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
	source, blueprint := templatebuiltin.GoLibrary()
	return catalog.NewMemory(catalog.Entry{
		ID:             "go-library",
		Name:           "Go library",
		Language:       "go",
		Kind:           "library",
		Description:    "A production-ready Go library with CI, governance, security, tests, and documentation",
		Licenses:       []string{"MIT"},
		DefaultLicense: "MIT",
		Source:         source,
		Blueprint:      blueprint,
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
	})
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
