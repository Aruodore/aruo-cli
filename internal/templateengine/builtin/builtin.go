// Package builtin exposes first-party template bundles embedded in Aruo.
package builtin

import (
	"embed"
	"io/fs"

	"github.com/aruodore/aruo/internal/templateengine"
)

//go:embed templates
var templates embed.FS

// GoLibrary returns the built-in Go library proof bundle.
func GoLibrary() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/go-library",
		Language: "go",
		Files: []templateengine.FileSpec{
			{Source: "go/library/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "go/library/go.mod.tmpl", Destination: "go.mod", Template: true},
			{Source: "go/library/gitignore", Destination: ".gitignore"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "foundation/docs-README.md.tmpl", Destination: "docs/README.md", Template: true},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "go/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "go/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "go/library/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "go/library/Makefile", Destination: "Makefile"},
			{Source: "go/library/library.go.tmpl", Destination: "{{ .Variables.PackageName }}.go", Template: true},
			{Source: "go/library/library_test.go.tmpl", Destination: "{{ .Variables.PackageName }}_test.go", Template: true},
		},
	}
}

// JSLibrary returns the built-in JavaScript library proof bundle.
func JSLibrary() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/js-library",
		Language: "javascript",
		Files: []templateengine.FileSpec{
			{Source: "js/library/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "js/library/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "js/library/gitignore", Destination: ".gitignore"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "foundation/docs-README.md.tmpl", Destination: "docs/README.md", Template: true},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "js/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "js/library/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "js/library/Makefile", Destination: "Makefile"},
			{Source: "js/library/index.js.tmpl", Destination: "src/index.js", Template: true},
			{Source: "js/library/index.test.js.tmpl", Destination: "src/index.test.js", Template: true},
		},
	}
}

// PythonLibrary returns the built-in Python library proof bundle.
func PythonLibrary() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/python-library",
		Language: "python",
		Files: []templateengine.FileSpec{
			{Source: "python/library/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "python/library/pyproject.toml.tmpl", Destination: "pyproject.toml", Template: true},
			{Source: "python/library/gitignore", Destination: ".gitignore"},
			{Source: "python/library/python-version", Destination: ".python-version"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "foundation/docs-README.md.tmpl", Destination: "docs/README.md", Template: true},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "python/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "python/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "python/library/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "python/library/Makefile", Destination: "Makefile"},
			{Source: "python/library/init.py.tmpl", Destination: "src/{{ .Variables.PackageName }}/__init__.py", Template: true},
			{Source: "python/library/test_library.py.tmpl", Destination: "tests/test_{{ .Variables.PackageName }}.py", Template: true},
		},
	}
}
