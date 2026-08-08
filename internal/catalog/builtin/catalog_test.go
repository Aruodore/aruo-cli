package builtin_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/templateengine"
)

func TestGoLibraryIsProductionReadyAndBuilds(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "library")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "go-library",
		Project: templateengine.Project{
			Name: "Example", Module: "example.com/example", Description: "An example library.",
			Author: "Example Authors", License: "MIT", Language: "go",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "go.mod", "example.go", "example_test.go", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
	command := exec.CommandContext(context.Background(), "go", "test", "./...")
	command.Dir = destination
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOCACHE=/tmp/aruo-generated-gocache")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated project go test: %v\n%s", err, output)
	}
}

func TestJSLibraryIsProductionReadyAndBuilds(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not installed")
	}
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "library")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "js-library",
		Project: templateengine.Project{
			Name: "Example", Module: "example-library", Description: "An example library.",
			Author: "Example Authors", License: "MIT", Language: "javascript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "src/index.js", "src/index.test.js", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
	command := exec.CommandContext(context.Background(), "node", "--test")
	command.Dir = destination
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated project node --test: %v\n%s", err, output)
	}
}

// TestTSLibraryHasRequiredFiles checks the generated file plan only. Unlike
// go-library/js-library/python-library, ts-library has a real devDependency
// (typescript) that must come from the network via npm install; running that
// here would make this package's tests non-hermetic, violating "no default
// test depends on external availability" (docs/development/testing.md). The
// generated project's own CI (.github/workflows/ci.yml) does run npm ci and
// npm test, where a real network-connected CI job is expected.
func TestTSLibraryHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "library")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "ts-library",
		Project: templateengine.Project{
			Name: "Example", Module: "example-library", Description: "An example library.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "tsconfig.json",
		"src/index.ts", "src/index.test.ts", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

// TestReactAppHasRequiredFiles checks the generated file plan only, for the
// same reason as TestTSLibraryHasRequiredFiles: react-app's dependencies
// (react, vite, vitest, jsdom, ...) must come from the network via npm
// install, so running that here would make this package's tests
// non-hermetic. The generated project's own CI does npm ci + npm test +
// npm run build for real.
func TestReactAppHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "app")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "react-app",
		Project: templateengine.Project{
			Name: "Example", Module: "example-app", Description: "An example app.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "tsconfig.json", "tsconfig.app.json",
		"tsconfig.node.json", "vite.config.ts", "index.html",
		"src/main.tsx", "src/App.tsx", "src/App.test.tsx", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

// TestNuxtAppHasRequiredFiles checks the generated file plan only, for the
// same reason as TestReactAppHasRequiredFiles: nuxt/vue/vitest/happy-dom
// must come from the network via npm install. The generated project's own
// CI does npm ci + npm test + npm run build for real.
func TestNuxtAppHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "app")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "nuxt-app",
		Project: templateengine.Project{
			Name: "Example", Module: "example-app", Description: "An example app.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "nuxt.config.ts", "tsconfig.json",
		"vitest.config.ts", "app/app.vue", "app/assets/css/main.css", "tests/app.test.ts", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

// TestVueLibraryHasRequiredFiles checks the generated file plan only, for
// the same reason as TestTSLibraryHasRequiredFiles: vue/vite/vitest/
// happy-dom must come from the network via npm install. The generated
// project's own CI does npm ci + npm test + npm run build for real.
func TestVueLibraryHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "library")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "vue-library",
		Project: templateengine.Project{
			Name: "Example", Module: "example-library", Description: "An example library.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "tsconfig.json", "vite.config.ts",
		"vitest.config.ts", "src/index.ts", "src/Greeting.vue", "src/__tests__/Greeting.test.ts", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

// TestVueAppHasRequiredFiles checks the generated file plan only, for the
// same reason as TestReactAppHasRequiredFiles: vue/vite/vitest/happy-dom
// must come from the network via npm install. The generated project's own
// CI does npm ci + npm test + npm run build for real.
func TestVueAppHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "app")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "vue-app",
		Project: templateengine.Project{
			Name: "Example", Module: "example-app", Description: "An example app.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "tsconfig.json", "vite.config.ts",
		"vitest.config.ts", "index.html", "src/main.ts", "src/App.vue", "src/App.test.ts", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

// TestNextAppHasRequiredFiles checks the generated file plan only, for the
// same reason as TestReactAppHasRequiredFiles: next/react/vitest/jsdom must
// come from the network via npm install. The generated project's own CI
// does npm ci + npm test + npm run build for real.
func TestNextAppHasRequiredFiles(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "app")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "next-app",
		Project: templateengine.Project{
			Name: "Example", Module: "example-app", Description: "An example app.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "package.json", "tsconfig.json", "next.config.ts",
		"vitest.config.ts", "app/layout.tsx", "app/page.tsx", "app/page.test.tsx", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
}

func TestPythonLibraryIsProductionReadyAndBuilds(t *testing.T) {
	t.Parallel()
	pythonBinary, err := exec.LookPath("python3")
	if err != nil {
		pythonBinary, err = exec.LookPath("python")
	}
	if err != nil {
		t.Skip("python is not installed")
	}
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "library")
	_, err = service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "python-library",
		Project: templateengine.Project{
			Name: "Example", Module: "example-library", Description: "An example library.",
			Author: "Example Authors", License: "MIT", Language: "python",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "LICENSE", "CHANGELOG.md", "ROADMAP.md", "SECURITY.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "aruo.yaml", "pyproject.toml", ".python-version",
		"src/example/__init__.py", "tests/test_example.py", "Makefile",
		"docs/README.md", ".github/workflows/ci.yml", ".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug.yml", ".github/ISSUE_TEMPLATE/feature.yml",
		".github/dependabot.yml", ".github/workflows/pr-title.yml", ".github/workflows/release.yml",
		"release-please-config.json", ".release-please-manifest.json",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
	command := exec.CommandContext(context.Background(), pythonBinary, "-m", "unittest", "discover", "-s", "tests")
	command.Dir = destination
	command.Env = append(os.Environ(), "PYTHONPATH=src")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated project python -m unittest: %v\n%s", err, output)
	}
}
