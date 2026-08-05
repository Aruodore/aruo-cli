package doctor_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/aruodore/aruo/internal/catalog/builtin"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
	"github.com/aruodore/aruo/internal/templateengine"
)

func TestGeneratedRepositoryScoresA(t *testing.T) {
	t.Parallel()
	templateCatalog, err := builtin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "project")
	_, err = creator.Create(context.Background(), create.Request{
		Destination: destination, TemplateID: "go-library",
		Project: templateengine.Project{
			Name: "Healthy", Module: "example.com/healthy", Description: "A healthy production-ready Go library.",
			Author: "Healthy Maintainers <security@example.com>", License: "MIT", Language: "go",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	service := newService(t)
	report, err := service.Audit(context.Background(), destination)
	if err != nil {
		t.Fatal(err)
	}
	if report.MaxScore != 100 || report.Score < 90 || report.Grade != "A" {
		t.Errorf("score = %d/%d grade=%s, want at least 90/A", report.Score, report.MaxScore, report.Grade)
	}
	if len(report.Categories) != 7 || len(report.Assessments) != 7 {
		t.Errorf("categories=%d assessments=%d", len(report.Categories), len(report.Assessments))
	}
}

func TestEmptyRepositoryProducesActionableZeroScore(t *testing.T) {
	t.Parallel()
	report, err := newService(t).Audit(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if report.Score != 0 || report.MaxScore != 100 || report.Grade != "F" {
		t.Errorf("report = %+v", report)
	}
	recommendations := 0
	for _, assessment := range report.Assessments {
		recommendations += len(assessment.Recommendations)
	}
	if recommendations < 10 {
		t.Errorf("recommendations = %d, want at least 10", recommendations)
	}
}

func TestServiceDoesNotFollowEscapingSymlink(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	outside := filepath.Join(parent, "outside.md")
	if err := os.WriteFile(outside, []byte(strings.Repeat("secret", 100)), 0o600); err != nil {
		t.Fatal(err)
	}
	repository := filepath.Join(parent, "repository")
	if err := os.Mkdir(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repository, "README.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, err := newService(t).Audit(context.Background(), repository)
	if err == nil {
		t.Fatal("Audit() error = nil; escaping symlink should not be read")
	}
}

func TestEngineRejectsInvalidCheckResults(t *testing.T) {
	t.Parallel()
	invalid := fakeCheck{}
	engine, err := doctor.NewEngine(invalid)
	if err != nil {
		t.Fatal(err)
	}
	repository, err := doctor.NewRepository(fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Run(context.Background(), ".", repository)
	if err == nil || !strings.Contains(err.Error(), "invalid assessment") {
		t.Fatalf("Run() error = %v", err)
	}
	if _, err := doctor.NewEngine(invalid, invalid); err == nil {
		t.Fatal("NewEngine(duplicate) error = nil")
	}
}

type fakeCheck struct{}

func (fakeCheck) ID() string                { return "plugin.invalid" }
func (fakeCheck) Category() doctor.Category { return doctor.CategorySecurity }
func (fakeCheck) MaxPoints() int            { return 2 }
func (fakeCheck) Run(context.Context, doctor.Repository) (doctor.Assessment, error) {
	return doctor.Assessment{ID: "plugin.invalid", Category: doctor.CategorySecurity, Points: 3, MaxPoints: 2}, nil
}

func newService(t *testing.T) *doctor.Service {
	t.Helper()
	engine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		t.Fatal(err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatal(err)
	}
	return service
}
