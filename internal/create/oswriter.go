package create

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aruodore/aruo/internal/templateengine"
)

// OSWriter materializes plans through same-parent staging and rename.
type OSWriter struct{}

// Write creates destination without overwriting any existing path.
func (OSWriter) Write(ctx context.Context, destination string, plan templateengine.Plan) error {
	if _, err := os.Lstat(destination); err == nil {
		return fmt.Errorf("destination already exists: %s", destination)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect destination: %w", err)
	}
	parent := filepath.Dir(destination)
	if info, err := os.Stat(parent); err != nil || !info.IsDir() {
		return fmt.Errorf("destination parent is unavailable: %s", parent)
	}
	staging, err := os.MkdirTemp(parent, ".aruo-create-*")
	if err != nil {
		return fmt.Errorf("create staging directory: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(staging)
		}
	}()

	for _, file := range plan.Files {
		if err := ctx.Err(); err != nil {
			return err
		}
		target := filepath.Join(staging, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create parent for %s: %w", file.Path, err)
		}
		if err := os.WriteFile(target, file.Content, file.Mode); err != nil {
			return fmt.Errorf("write %s: %w", file.Path, err)
		}
	}
	if err := os.Rename(staging, destination); err != nil {
		return fmt.Errorf("commit repository: %w", err)
	}
	committed = true
	return nil
}
