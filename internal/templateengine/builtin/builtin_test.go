package builtin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/templateengine"
	"github.com/aruodore/aruo-cli/internal/templateengine/builtin"
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
		Variables: map[string]any{"IncludeInstall": false, "TemplateID": "react"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 33 {
		t.Errorf("file count = %d, want 33", len(files))
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
		Variables: map[string]any{"IncludeInstall": false, "TemplateID": "nuxt"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 48 {
		t.Errorf("file count = %d, want 48", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-app"`) || !strings.Contains(pkg, `"nuxt"`) || !strings.Contains(pkg, `"check"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if app := files["app/app.vue"]; !strings.Contains(app, "<h1>Example</h1>") {
		t.Errorf("app/app.vue = %q", app)
	}
	for _, path := range []string{
		"AGENTS.md", ".env.example", "Dockerfile", "compose.yaml", "aruo.yaml",
		"server/api/health/live.get.ts", "server/api/health/ready.get.ts", "server/middleware/security-headers.ts",
		"server/db/migrations/0000_initial.sql", ".github/workflows/ci.yml",
	} {
		if _, exists := files[path]; !exists {
			t.Errorf("missing production contract file %q", path)
		}
	}
	if manifest := files["aruo.yaml"]; !strings.Contains(manifest, "status: REQUIRED") || !strings.Contains(manifest, "status: SOLVED") {
		t.Errorf("aruo.yaml does not expose capabilities and limitations: %q", manifest)
	}
	if agents := files["AGENTS.md"]; !strings.Contains(agents, "Never disable validation") || !strings.Contains(agents, "server-side") {
		t.Errorf("AGENTS.md lacks required safety constraints: %q", agents)
	}
}

func TestVueLibrary(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.VueLibrary()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-library",
			Description: "An example Vue library.",
			License:     "MIT",
			Language:    "typescript",
		},
		Variables: map[string]any{"IncludeInstall": true, "TemplateID": "vue-library"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 28 {
		t.Errorf("file count = %d, want 28", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-library"`) || !strings.Contains(pkg, `"vue"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if greeting := files["src/Greeting.vue"]; !strings.Contains(greeting, "{{ name }}") {
		t.Errorf("src/Greeting.vue = %q", greeting)
	}
}

func TestNextApp(t *testing.T) {
	t.Parallel()

	source, blueprint := builtin.NextApp()
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Render(context.Background(), blueprint, templateengine.Data{
		Project: templateengine.Project{
			Name:        "Example",
			Module:      "example-app",
			Description: "An example Next.js app.",
			License:     "MIT",
			Language:    "typescript",
		},
		Variables: map[string]any{"IncludeInstall": false, "TemplateID": "next"},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	files := make(map[string]string, len(plan.Files))
	for _, file := range plan.Files {
		files[file.Path] = string(file.Content)
	}
	if len(files) != 31 {
		t.Errorf("file count = %d, want 31", len(files))
	}
	if pkg := files["package.json"]; !strings.Contains(pkg, `"name": "example-app"`) || !strings.Contains(pkg, `"next"`) {
		t.Errorf("package.json = %q", pkg)
	}
	if page := files["app/page.tsx"]; !strings.Contains(page, "<h1>Example</h1>") {
		t.Errorf("app/page.tsx = %q", page)
	}
	if layout := files["app/layout.tsx"]; !strings.Contains(layout, `title: "Example"`) {
		t.Errorf("app/layout.tsx = %q", layout)
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
