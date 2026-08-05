package doctor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Service opens repository roots and delegates policy evaluation.
type Service struct{ engine *Engine }

// NewService constructs a doctor service.
func NewService(engine *Engine) (*Service, error) {
	if engine == nil {
		return nil, fmt.Errorf("doctor engine is required")
	}
	return &Service{engine: engine}, nil
}

// Audit evaluates path without executing repository content or using network access.
func (s *Service) Audit(ctx context.Context, target string) (Report, error) {
	absolute, err := filepath.Abs(target)
	if err != nil {
		return Report{}, fmt.Errorf("resolve repository path: %w", err)
	}
	root, err := os.OpenRoot(absolute)
	if err != nil {
		return Report{}, fmt.Errorf("open repository: %w", err)
	}
	defer root.Close()
	repository, err := NewRepository(root.FS())
	if err != nil {
		return Report{}, err
	}
	return s.engine.Run(ctx, absolute, repository)
}
