// Package catalog describes qualified project templates available to creation workflows.
package catalog

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

// Entry is a resolved, qualified project template.
type Entry struct {
	ID          string
	Name        string
	Language    string
	Kind        string
	Description string
	// Color is the ecosystem's real brand color as a truecolor hex value
	// (for example Go's "#00ADD8"), an identity fact like Language rather
	// than a UI choice. Consumers with no color rendering ignore it;
	// color-capable renderers apply it only when the active color policy
	// permits color at all.
	Color          string
	Licenses       []string
	DefaultLicense string
	Source         fs.FS
	Blueprint      templateengine.Blueprint
	Defaults       map[string]any
	Prepare        func(templateengine.Data) (templateengine.Data, error)
	NextSteps      []string
	Prompts        ProjectPrompts
}

// ProjectPrompts gives each template concrete, ecosystem-native language for
// project metadata instead of exposing Aruo's generic internal field names.
type ProjectPrompts struct {
	ModuleLabel       string
	ModuleDescription string
	ModuleExample     string
}

// Catalog resolves templates without exposing source-specific discovery to callers.
type Catalog interface {
	List(context.Context) ([]Entry, error)
	Resolve(context.Context, string) (Entry, error)
}

// Memory is an immutable built-in catalog implementation.
type Memory struct{ entries map[string]Entry }

// NewMemory validates and constructs a catalog.
func NewMemory(entries ...Entry) (*Memory, error) {
	result := &Memory{entries: make(map[string]Entry, len(entries))}
	for _, entry := range entries {
		if entry.ID == "" || entry.Name == "" || entry.Language == "" || entry.Kind == "" || entry.Source == nil || entry.DefaultLicense == "" {
			return nil, errors.New("catalog entries require ID, name, language, kind, and source")
		}
		if _, exists := result.entries[entry.ID]; exists {
			return nil, fmt.Errorf("duplicate catalog entry %q", entry.ID)
		}
		result.entries[entry.ID] = entry
	}
	return result, nil
}

// List returns entries grouped by kind, then ordered by stable ID within
// each kind, so a picker over many entries reads as related groups (every
// library together, every app together) rather than one flat alphabetical
// list interleaving unrelated kinds.
func (c *Memory) List(ctx context.Context) ([]Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(c.entries))
	for _, entry := range c.entries {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}
		return entries[i].ID < entries[j].ID
	})
	return entries, nil
}

// Resolve returns the entry with ID.
func (c *Memory) Resolve(ctx context.Context, id string) (Entry, error) {
	if err := ctx.Err(); err != nil {
		return Entry{}, err
	}
	entry, exists := c.entries[id]
	if !exists {
		return Entry{}, fmt.Errorf("unknown template %q", id)
	}
	return entry, nil
}
