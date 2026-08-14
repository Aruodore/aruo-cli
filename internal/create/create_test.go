package create_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/templateengine"
)

func TestServiceCreatesRepository(t *testing.T) {
	t.Parallel()
	entry := catalog.Entry{
		ID: "generic", Name: "Generic", Language: "test", Kind: "library",
		Licenses: []string{"MIT"}, DefaultLicense: "MIT",
		Source: fstest.MapFS{"file": {Data: []byte("hello {{ .Project.Name }}")}},
		Blueprint: templateengine.Blueprint{ID: "generic", Language: "test", Files: []templateengine.FileSpec{
			{Source: "file", Destination: "README.md", Template: true},
		}},
	}
	templateCatalog, err := catalog.NewMemory(entry)
	if err != nil {
		t.Fatal(err)
	}
	service, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "project")
	result, err := service.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "generic",
		Project: templateengine.Project{Name: "World", Language: "test"},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(destination, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello World" || result.FileCount != 13 {
		t.Errorf("content=%q result=%+v", content, result)
	}
	for _, name := range []string{"AGENTS.md", "aruo.yaml", ".aruo/contract.yaml", ".aruo/managed.json", ".aruo/rules/security.md"} {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); err != nil {
			t.Errorf("created repository is missing %s: %v", name, err)
		}
	}
}

func TestOSWriterRefusesExistingFileDestination(t *testing.T) {
	t.Parallel()
	destination := filepath.Join(t.TempDir(), "already-a-file")
	if err := os.WriteFile(destination, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := (create.OSWriter{}).Write(context.Background(), destination, templateengine.Plan{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Write() error = %v", err)
	}
}

func TestOSWriterRefusesNonEmptyExistingDirectory(t *testing.T) {
	t.Parallel()
	destination := t.TempDir()
	if err := os.WriteFile(filepath.Join(destination, "existing"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := (create.OSWriter{}).Write(context.Background(), destination, templateengine.Plan{
		Files: []templateengine.File{{Path: "README.md", Content: []byte("hi"), Mode: 0o644}},
	})
	if err == nil || !strings.Contains(err.Error(), "already exists and is not empty") {
		t.Fatalf("Write() error = %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(destination, "README.md")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("README.md stat error = %v, want the refused write to leave the directory untouched", statErr)
	}
}

// TestOSWriterAdoptsExistingEmptyDirectory covers the "aruo create ." case:
// the destination already exists (the current directory always does) but
// is empty, so Write must populate it in place instead of refusing it.
func TestOSWriterAdoptsExistingEmptyDirectory(t *testing.T) {
	t.Parallel()
	destination := t.TempDir()
	err := (create.OSWriter{}).Write(context.Background(), destination, templateengine.Plan{
		Files: []templateengine.File{
			{Path: "README.md", Content: []byte("hi"), Mode: 0o644},
			{Path: "src/main.go", Content: []byte("package main"), Mode: 0o644},
		},
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(destination, "README.md"))
	if err != nil || string(content) != "hi" {
		t.Fatalf("README.md = %q, err = %v", content, err)
	}
	if _, statErr := os.Stat(filepath.Join(destination, "src", "main.go")); statErr != nil {
		t.Fatalf("src/main.go stat error = %v", statErr)
	}
	entries, err := os.ReadDir(destination)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".aruo-create-") {
			t.Errorf("staging directory %q leaked into the adopted destination", entry.Name())
		}
	}
}

func TestOSWriterCancellationLeavesNoDestination(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	destination := filepath.Join(t.TempDir(), "project")
	err := (create.OSWriter{}).Write(ctx, destination, templateengine.Plan{Files: []templateengine.File{{Path: "a", Content: []byte("a"), Mode: 0o644}}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Write() error = %v", err)
	}
	if _, statErr := os.Stat(destination); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("destination stat error = %v", statErr)
	}
}
