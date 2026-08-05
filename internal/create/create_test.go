package create_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/aruodore/aruo/internal/catalog"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/templateengine"
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
	if string(content) != "hello World" || result.FileCount != 1 {
		t.Errorf("content=%q result=%+v", content, result)
	}
}

func TestOSWriterRefusesExistingDestination(t *testing.T) {
	t.Parallel()
	destination := t.TempDir()
	err := (create.OSWriter{}).Write(context.Background(), destination, templateengine.Plan{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Write() error = %v", err)
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
