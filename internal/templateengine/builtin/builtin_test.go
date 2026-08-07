package builtin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aruodore/aruo/internal/templateengine"
	"github.com/aruodore/aruo/internal/templateengine/builtin"
)

func TestGoLibrary(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.GoLibrary()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example.com/project",
			Description: "An example Go library.",
			License:     "MIT",
			Language:    "go",
		},
		Variables: map[string]any{"IncludeInstall": true, "GoVersion": "1.26.0", "TemplateID": "go-library", "PackageName": "example"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 24 {
		t.Errorf("file count = %d, want 24", len(files))
	}
	if readme := files["README.md"]; !strings.Contains(readme, "go get example.com/project") {
		t.Errorf("README = %q", readme)
	}
	if module := files["go.mod"]; module != "module example.com/project\n\ngo 1.26.0\n" {
		t.Errorf("go.mod = %q", module)
	}
}

func TestJSLibrary(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.JSLibrary()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-project",
			Description: "An example JavaScript library.",
			License:     "MIT",
			Language:    "javascript",
		},
		Variables: map[string]any{"IncludeInstall": true, "TemplateID": "js-library"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 24 {
		t.Errorf("file count = %d, want 24", len(files))
	}
	if readme := files["README.md"]; !strings.Contains(readme, "npm install example-project") {
		t.Errorf("README = %q", readme)
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-project"`) {
		t.Errorf("package.json = %q", pkg)
	}
}

func TestTSLibrary(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.TSLibrary()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-project",
			Description: "An example TypeScript library.",
			License:     "MIT",
			Language:    "typescript",
		},
		Variables: map[string]any{"IncludeInstall": true, "TemplateID": "ts-library"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 25 {
		t.Errorf("file count = %d, want 25", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-project"`) || !strings.Contains(pkg, `"typescript"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if tsconfig := files["tsconfig.json"]; !strings.Contains(tsconfig, `"strict": true`) {
		t.Errorf("tsconfig.json = %q", tsconfig)
	}
}

func TestReactApp(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.ReactApp()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-app",
			Description: "An example React app.",
			License:     "MIT",
			Language:    "typescript",
		},
		Variables: map[string]any{"IncludeInstall": false, "TemplateID": "react-app"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 30 {
		t.Errorf("file count = %d, want 30", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-app"`) || !strings.Contains(pkg, `"react"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if app := files["src/App.tsx"]; !strings.Contains(app, "<h1>Example</h1>") {
		t.Errorf("src/App.tsx = %q", app)
	}
	if test := files["src/App.test.tsx"]; !strings.Contains(test, `name: "Example"`) {
		t.Errorf("src/App.test.tsx = %q", test)
	}
	if index := files["index.html"]; !strings.Contains(index, "<title>Example</title>") {
		t.Errorf("index.html = %q", index)
	}
}

func TestNuxtApp(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.NuxtApp()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-app",
			Description: "An example Nuxt app.",
			License:     "MIT",
			Language:    "typescript",
		},
		Variables: map[string]any{"IncludeInstall": false, "TemplateID": "nuxt-app"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 27 {
		t.Errorf("file count = %d, want 27", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-app"`) || !strings.Contains(pkg, `"nuxt"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if app := files["app/app.vue"]; !strings.Contains(app, "<h1>Example</h1>") {
		t.Errorf("app/app.vue = %q", app)
	}
	if test := files["tests/app.test.ts"]; !strings.Contains(test, `toContain("Example")`) {
		t.Errorf("tests/app.test.ts = %q", test)
	}
}

func TestPythonLibrary(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.PythonLibrary()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-project",
			Description: "An example Python library.",
			License:     "MIT",
			Language:    "python",
		},
		Variables: map[string]any{"IncludeInstall": true, "TemplateID": "python-library", "PackageName": "example"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 25 {
		t.Errorf("file count = %d, want 25", len(files))
	}
	if readme := files["README.md"]; !strings.Contains(readme, "pip install example-project") {
		t.Errorf("README = %q", readme)
	}
	if pyproject := files["pyproject.toml"]; !strings.Contains(pyproject, `name = "example-project"`) {
		t.Errorf("pyproject.toml = %q", pyproject)
	}
	if init := files["src/example/__init__.py"]; !strings.Contains(init, `VERSION = "0.1.0"`) {
		t.Errorf("src/example/__init__.py = %q", init)
	}
}
