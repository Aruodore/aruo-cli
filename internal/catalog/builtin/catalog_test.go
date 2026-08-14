package builtin_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/templateengine"
)

func TestEveryApplicationFrameworkHasOneCanonicalTemplate(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"react", "vue", "next", "nuxt"} {
		entry, resolveErr := templateCatalog.Resolve(context.Background(), id)
		if resolveErr != nil {
			t.Errorf("Resolve(%q): %v", id, resolveErr)
			continue
		}
		if entry.Kind != "app" {
			t.Errorf("%s kind = %q, want app", id, entry.Kind)
		}
		if entry.Defaults["TemplateID"] != id {
			t.Errorf("%s manifest ID default = %v", id, entry.Defaults["TemplateID"])
		}
	}
	for _, removed := range []string{"react-lean", "vue-lean", "next-lean", "nuxt-lean"} {
		if _, resolveErr := templateCatalog.Resolve(context.Background(), removed); resolveErr == nil {
			t.Errorf("removed template %q still resolves", removed)
		}
	}
	entries, err := templateCatalog.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	appCount := 0
	for _, entry := range entries {
		if entry.Kind == "app" {
			appCount++
		}
	}
	if appCount != 4 {
		t.Fatalf("application template count = %d, want 4", appCount)
	}
}

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
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOCACHE="+filepath.Join(t.TempDir(), "gocache"))
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
// same reason as TestTSLibraryHasRequiredFiles: React's dependencies
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
		Destination: destination, TemplateID: "react",
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
		"src/main.tsx", "src/app.tsx", "src/app.test.tsx", "Makefile",
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
	assertCanonicalApplicationManifest(t, destination, "react")
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
		Destination: destination, TemplateID: "nuxt",
		Project: templateengine.Project{
			Name: "Example", Module: "example-app", Description: "An example app.",
			Author: "Example Authors", License: "MIT", Language: "typescript",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"README.md", "AGENTS.md", "AGENTS.local.md", "aruo.yaml", "package.json", "nuxt.config.ts",
		"docs/README.md", "app/app.vue", "tests/app.test.ts", ".github/workflows/ci.yml",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}
	for _, retired := range []string{".env.example", "Dockerfile", "compose.yaml", "docs/operations.md", "server/db/schema.ts"} {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(retired))); !os.IsNotExist(err) {
			t.Errorf("retired comprehensive file %s still exists", retired)
		}
	}
	assertCanonicalApplicationManifest(t, destination, "nuxt")
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
		"vitest.config.ts", "src/index.ts", "src/greeting.vue", "src/__tests__/greeting.test.ts", "Makefile",
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
		Destination: destination, TemplateID: "vue",
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
		"vitest.config.ts", "index.html", "src/main.ts", "src/app.vue", "src/app.test.ts", "Makefile",
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
	assertCanonicalApplicationManifest(t, destination, "vue")
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
		Destination: destination, TemplateID: "next",
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
	assertCanonicalApplicationManifest(t, destination, "next")
}

func assertCanonicalApplicationManifest(t *testing.T, destination, templateID string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(destination, "aruo.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	manifest := string(content)
	if !strings.Contains(manifest, "id: "+templateID) {
		t.Errorf("aruo.yaml does not use canonical template ID %q: %s", templateID, manifest)
	}
	if strings.Contains(manifest, "profile:") || strings.Contains(manifest, "-lean") {
		t.Errorf("aruo.yaml still exposes a removed application profile: %s", manifest)
	}
	if _, err := os.Stat(filepath.Join(destination, "AGENTS.local.md")); err != nil {
		t.Errorf("application-owned stack guidance is missing: %v", err)
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
