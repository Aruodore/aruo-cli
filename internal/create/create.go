// Package create orchestrates safe creation of new repositories.
package create

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/aruodore/aruo-cli/internal/catalog"
	"github.com/aruodore/aruo-cli/internal/initialize"
	"github.com/aruodore/aruo-cli/internal/templateengine"
)

// Request contains resolved command input independent of presentation.
type Request struct {
	Destination string
	TemplateID  string
	Project     templateengine.Project
	Variables   map[string]any
}

// Result describes the created repository.
type Result struct {
	Destination string
	TemplateID  string
	FileCount   int
	NextSteps   []string
}

// Writer atomically materializes a complete file plan.
type Writer interface {
	Write(context.Context, string, templateengine.Plan) error
}

// Service creates repositories from catalog entries.
type Service struct {
	catalog catalog.Catalog
	writer  Writer
}

// NewService constructs a creation service.
func NewService(templateCatalog catalog.Catalog, writer Writer) (*Service, error) {
	if templateCatalog == nil || writer == nil {
		return nil, errors.New("create service requires catalog and writer")
	}
	return &Service{catalog: templateCatalog, writer: writer}, nil
}

// Create renders, validates, and atomically writes one new repository.
func (s *Service) Create(ctx context.Context, request Request) (Result, error) {
	if request.Destination == "" {
		return Result{}, errors.New("destination is required")
	}
	entry, err := s.catalog.Resolve(ctx, request.TemplateID)
	if err != nil {
		return Result{}, err
	}
	data := templateengine.Data{Project: request.Project, Variables: maps.Clone(entry.Defaults)}
	if data.Variables == nil {
		data.Variables = map[string]any{}
	}
	maps.Copy(data.Variables, request.Variables)
	if entry.Prepare != nil {
		data, err = entry.Prepare(data)
		if err != nil {
			return Result{}, fmt.Errorf("prepare template %q: %w", entry.ID, err)
		}
	}
	engine, err := templateengine.New(entry.Source, templateengine.Options{})
	if err != nil {
		return Result{}, err
	}
	plan, err := engine.Render(ctx, entry.Blueprint, data)
	if err != nil {
		return Result{}, err
	}
	plan, err = initialize.AugmentPlan(plan)
	if err != nil {
		return Result{}, fmt.Errorf("add Aruo engineering contract: %w", err)
	}
	destination, err := filepath.Abs(request.Destination)
	if err != nil {
		return Result{}, fmt.Errorf("resolve destination: %w", err)
	}
	if err := s.writer.Write(ctx, destination, plan); err != nil {
		return Result{}, err
	}
	return Result{
		Destination: destination,
		TemplateID:  entry.ID,
		FileCount:   len(plan.Files),
		NextSteps:   append([]string(nil), entry.NextSteps...),
	}, nil
}
