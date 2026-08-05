// Package templateengine renders validated repository template bundles into file plans.
package templateengine

import "io/fs"

// Blueprint describes a resolved, language-specific set of template files.
type Blueprint struct {
	ID       string
	Language string
	Files    []FileSpec
}

// FileSpec maps one bundle source to one generated destination.
type FileSpec struct {
	Source      string
	Destination string
	Mode        fs.FileMode
	Template    bool
}

// Project contains stable, commonly required repository metadata.
type Project struct {
	Name        string
	Module      string
	Description string
	Author      string
	License     string
	Language    string
}

// Data is the complete explicit input available to templates.
type Data struct {
	Project   Project
	Variables map[string]any
}

// File is one planned repository artifact.
type File struct {
	Path    string
	Content []byte
	Mode    fs.FileMode
	Source  string
}

// Plan is a deterministic, side-effect-free rendering result.
type Plan struct {
	BlueprintID string
	Files       []File
}
