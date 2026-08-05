package builtin_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aruodore/aruo/internal/catalog/builtin"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/templateengine"
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
