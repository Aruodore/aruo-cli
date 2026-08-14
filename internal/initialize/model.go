// Package initialize installs Aruo's managed engineering contract into an existing repository.
package initialize

// Request describes one repository adoption operation.
type Request struct {
	Target string
}

// Stack is the deterministic local stack observation written by init.
type Stack struct {
	Ecosystem      string   `json:"ecosystem"`
	Frameworks     []string `json:"frameworks"`
	PackageManager string   `json:"packageManager,omitempty"`
}

// FileChange describes one proposed managed file.
type FileChange struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

// Plan is a side-effect-free initialization proposal.
type Plan struct {
	Repository string       `json:"repository"`
	Stack      Stack        `json:"stack"`
	Changes    []FileChange `json:"changes"`
	Conflicts  []string     `json:"conflicts,omitempty"`
	files      []plannedFile
}

// Result describes a completed or dry-run initialization.
type Result struct {
	Repository string       `json:"repository"`
	Stack      Stack        `json:"stack"`
	Changes    []FileChange `json:"changes"`
	DryRun     bool         `json:"dryRun"`
}

type plannedFile struct {
	path    string
	content []byte
}
