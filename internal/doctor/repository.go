package doctor

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

const maxAuditFileBytes = 2 << 20

// Repository is a bounded, read-only view of a repository root.
type Repository struct{ fs fs.FS }

// NewRepository constructs a repository view.
func NewRepository(source fs.FS) (Repository, error) {
	if source == nil {
		return Repository{}, errors.New("repository filesystem is required")
	}
	return Repository{fs: source}, nil
}

// Exists reports whether a regular file or directory exists.
func (r Repository) Exists(name string) (bool, error) {
	info, err := fs.Stat(r.fs, name)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular() || info.IsDir(), nil
}

// ReadText returns a bounded text file.
func (r Repository) ReadText(name string) (string, error) {
	file, err := r.fs.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, maxAuditFileBytes+1))
	if err != nil {
		return "", err
	}
	if len(content) > maxAuditFileBytes {
		return "", fmt.Errorf("%s exceeds audit read limit", name)
	}
	return string(content), nil
}

// Glob returns sorted matches using repository-relative slash paths.
func (r Repository) Glob(pattern string) ([]string, error) { return fs.Glob(r.fs, pattern) }

// TestFiles returns likely language-native test files without entering dependency trees.
func (r Repository) TestFiles() ([]string, error) {
	var result []string
	err := fs.WalkDir(r.fs, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && name != "." {
			base := path.Base(name)
			if base == ".git" || base == "vendor" || base == "node_modules" || base == ".venv" || base == "dist" {
				return fs.SkipDir
			}
		}
		if entry.IsDir() {
			return nil
		}
		base := strings.ToLower(path.Base(name))
		if strings.HasSuffix(base, "_test.go") || strings.HasSuffix(base, "_test.py") ||
			strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") ||
			strings.HasSuffix(base, "_test.rs") {
			result = append(result, name)
		}
		return nil
	})
	return result, err
}
