package create

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

// OSWriter materializes plans through same-parent staging and rename.
type OSWriter struct{}

// Write creates destination without overwriting any existing content. An
// already-existing destination is accepted only when it's an empty
// directory (for example ".", or a directory the caller just created) --
// everything else (a file, a non-empty directory, a path Lstat can't read)
// is refused.
func (OSWriter) Write(ctx context.Context, destination string, plan templateengine.Plan) error {
	adopting, err := destinationIsEmptyDir(destination)
	if err != nil {
		return err
	}

	parent := filepath.Dir(destination)
	if info, statErr := os.Stat(parent); statErr != nil || !info.IsDir() {
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

	if adopting {
		if err := adoptStaging(staging, destination); err != nil {
			return fmt.Errorf("commit repository: %w", err)
		}
	} else if err := os.Rename(staging, destination); err != nil {
		return fmt.Errorf("commit repository: %w", err)
	}
	committed = true
	return nil
}

// destinationIsEmptyDir reports whether Write should adopt an
// already-existing destination: false and no error when destination
// doesn't exist yet (the common case, created fresh via a single rename);
// true when it exists and is an empty directory (populated in place
// instead, since a brand-new directory can't be renamed onto it -- see
// adoptStaging). Anything else is refused with an error.
func destinationIsEmptyDir(destination string) (bool, error) {
	info, err := os.Lstat(destination)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect destination: %w", err)
	}
	if !info.IsDir() {
		return false, fmt.Errorf("destination already exists: %s", destination)
	}
	entries, err := os.ReadDir(destination)
	if err != nil {
		return false, fmt.Errorf("inspect destination: %w", err)
	}
	if len(entries) > 0 {
		return false, fmt.Errorf("destination already exists and is not empty: %s", destination)
	}
	return true, nil
}

// adoptStaging moves every top-level staged entry into an already-existing
// empty destination directory, then removes the now-empty staging
// directory. Individual renames, not one atomic directory swap: renaming a
// new directory onto an existing one isn't reliably supported across
// platforms, and destination may be "." itself, which can never be
// replaced at all. If a rename partway through fails, every entry already
// moved is best-effort rolled back so destination is left exactly as
// empty as it started, rather than partially populated.
func adoptStaging(staging, destination string) error {
	entries, err := os.ReadDir(staging)
	if err != nil {
		return err
	}
	moved := make([]string, 0, len(entries))
	for _, entry := range entries {
		oldPath := filepath.Join(staging, entry.Name())
		newPath := filepath.Join(destination, entry.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			for _, name := range moved {
				_ = os.RemoveAll(filepath.Join(destination, name))
			}
			return err
		}
		moved = append(moved, entry.Name())
	}
	return os.Remove(staging)
}
