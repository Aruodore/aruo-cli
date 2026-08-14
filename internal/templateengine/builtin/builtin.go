// Package builtin exposes first-party template bundles embedded in Aruo.
package builtin

import (
	"embed"
	"io/fs"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

//go:embed templates
var templates embed.FS

func appendFullstackFiles(files []templateengine.FileSpec, packageSource string) []templateengine.FileSpec {
	filtered := files[:0]
	for _, file := range files {
		if file.Destination != "AGENTS.md" && file.Destination != "aruo.yaml" && file.Destination != "package.json" && file.Destination != "README.md" {
			filtered = append(filtered, file)
		}
	}
	return append(filtered,
		fullstackFile("fullstack/server/env.ts", "server/env.ts"),
		fullstackFile("fullstack/server/errors.ts", "server/errors.ts"),
		fullstackFile("fullstack/server/logger.ts", "server/logger.ts"),
		fullstackFile("fullstack/server/database.ts", "server/database.ts"),
		fullstackFile("fullstack/server/rate-limit.ts", "server/rate-limit.ts"),
		fullstackFile("fullstack/server/index.ts", "server/index.ts"),
		fullstackFile("fullstack/server/schema.ts", "server/db/schema.ts"),
		fullstackFile("fullstack/server/migration.sql", "server/db/migrations/0000_baseline.sql"),
		fullstackFile("fullstack/server/env.test.ts", "server/env.test.ts"),
		fullstackFile("fullstack/env.example", ".env.example"),
		fullstackFile("fullstack/compose.yaml", "compose.yaml"),
		fullstackFile("fullstack/Dockerfile", "Dockerfile"),
		fullstackFile("fullstack/dockerignore", ".dockerignore"),
		fullstackFile("fullstack/tsconfig.server.json", "tsconfig.server.json"),
		fullstackFile("fullstack/drizzle.config.ts", "drizzle.config.ts"),
		fullstackFile("fullstack/AGENTS.md", "AGENTS.md"),
		fullstackFile("fullstack/docs-architecture.md", "docs/architecture.md"),
		fullstackFile("fullstack/docs-operations.md", "docs/operations.md"),
		fullstackTemplate("fullstack/aruo.yaml.tmpl", "aruo.yaml"),
		fullstackTemplate(packageSource, "package.json"),
		fullstackTemplate("fullstack/README.md.tmpl", "README.md"),
	)
}

func fullstackFile(source, destination string) templateengine.FileSpec {
	return templateengine.FileSpec{Source: source, Destination: destination}
}

func fullstackTemplate(source, destination string) templateengine.FileSpec {
	return templateengine.FileSpec{Source: source, Destination: destination, Template: true}
}

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

// TSLibrary returns the built-in TypeScript library proof bundle.
func TSLibrary() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/ts-library",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "ts/library/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "ts/library/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "ts/library/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "ts/library/gitignore", Destination: ".gitignore"},
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
			{Source: "ts/library/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "ts/library/Makefile", Destination: "Makefile"},
			{Source: "ts/library/index.ts.tmpl", Destination: "src/index.ts", Template: true},
			{Source: "ts/library/index.test.ts.tmpl", Destination: "src/index.test.ts", Template: true},
		},
	}
}

// ReactLean returns the reduced React application baseline.
func ReactLean() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/react-lean",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "react/app/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "react/app/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "react/app/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "react/app/tsconfig.app.json", Destination: "tsconfig.app.json"},
			{Source: "react/app/tsconfig.node.json", Destination: "tsconfig.node.json"},
			{Source: "react/app/vite.config.ts", Destination: "vite.config.ts"},
			{Source: "react/app/index.html.tmpl", Destination: "index.html", Template: true},
			{Source: "react/app/gitignore", Destination: ".gitignore"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "js/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "react/app/main.tsx", Destination: "src/main.tsx"},
			{Source: "frontend/app/main.css", Destination: "src/main.css"},
			{Source: "react/app/App.tsx.tmpl", Destination: "src/app.tsx", Template: true},
			{Source: "react/app/App.test.tsx.tmpl", Destination: "src/app.test.tsx", Template: true},
			{Source: "frontend/app/AGENTS.md", Destination: "AGENTS.md"},
			{Source: "frontend/app/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "frontend/app/docs-README.md", Destination: "docs/README.md"},
			{Source: "frontend/app/prettierignore", Destination: ".prettierignore"},
			{Source: "frontend/app/eslint-react.config.mjs", Destination: "eslint.config.mjs"},
			{Source: "frontend/app/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "frontend/app/Makefile", Destination: "Makefile"},
		},
	}
}

// ReactApp returns the comprehensive React modular-monolith baseline.
func ReactApp() (fs.FS, templateengine.Blueprint) {
	source, blueprint := ReactLean()
	blueprint.ID = "aruo/react"
	blueprint.Files = appendFullstackFiles(blueprint.Files, "react/full/package.json.tmpl")
	return source, blueprint
}

// NuxtApp returns the built-in Nuxt application proof bundle.
func NuxtApp() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/nuxt",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "nuxt/production/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "nuxt/production/AGENTS.md", Destination: "AGENTS.md"},
			{Source: "nuxt/production/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "nuxt/production/nuxt.config.ts", Destination: "nuxt.config.ts"},
			{Source: "nuxt/production/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "nuxt/production/eslint.config.mjs", Destination: "eslint.config.mjs"},
			{Source: "nuxt/production/vitest.config.ts", Destination: "vitest.config.ts"},
			{Source: "nuxt/production/drizzle.config.ts", Destination: "drizzle.config.ts"},
			{Source: "nuxt/production/prettierignore", Destination: ".prettierignore"},
			{Source: "nuxt/production/gitignore", Destination: ".gitignore"},
			{Source: "nuxt/production/env.example", Destination: ".env.example"},
			{Source: "nuxt/production/app.vue.tmpl", Destination: "app/app.vue", Template: true},
			{Source: "nuxt/app/main.css", Destination: "app/assets/css/main.css"},
			{Source: "nuxt/production/env.ts", Destination: "server/utils/env.ts"},
			{Source: "nuxt/production/errors.ts", Destination: "server/utils/errors.ts"},
			{Source: "nuxt/production/request-id.ts", Destination: "server/utils/request-id.ts"},
			{Source: "nuxt/production/logger.ts", Destination: "server/utils/logger.ts"},
			{Source: "nuxt/production/request-context.ts", Destination: "server/middleware/request-context.ts"},
			{Source: "nuxt/production/security-headers.ts", Destination: "server/middleware/security-headers.ts"},
			{Source: "nuxt/production/db.ts", Destination: "server/db/client.ts"},
			{Source: "nuxt/production/schema.ts", Destination: "server/db/schema.ts"},
			{Source: "nuxt/production/migrate.ts", Destination: "server/db/migrate.ts"},
			{Source: "nuxt/production/migration.sql", Destination: "server/db/migrations/0000_initial.sql"},
			{Source: "nuxt/production/migration-journal.json", Destination: "server/db/migrations/meta/_journal.json"},
			{Source: "nuxt/production/health-live.ts", Destination: "server/api/health/live.get.ts"},
			{Source: "nuxt/production/health-ready.ts", Destination: "server/api/health/ready.get.ts"},
			{Source: "nuxt/production/env.test.ts", Destination: "tests/env.test.ts"},
			{Source: "nuxt/production/errors.test.ts", Destination: "tests/errors.test.ts"},
			{Source: "nuxt/production/Dockerfile", Destination: "Dockerfile"},
			{Source: "nuxt/production/dockerignore", Destination: ".dockerignore"},
			{Source: "nuxt/production/compose.yaml", Destination: "compose.yaml"},
			{Source: "nuxt/production/Makefile", Destination: "Makefile"},
			{Source: "nuxt/production/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "nuxt/production/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "nuxt/production/docs-architecture.md", Destination: "docs/architecture.md"},
			{Source: "nuxt/production/docs-operations.md", Destination: "docs/operations.md"},
			{Source: "nuxt/production/docs-README.md", Destination: "docs/README.md"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
		},
	}
}

// NuxtLean returns the reduced Nuxt application baseline.
func NuxtLean() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/nuxt-lean",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "nuxt/app/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "nuxt/app/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "nuxt/app/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "nuxt/app/nuxt.config.ts", Destination: "nuxt.config.ts"},
			{Source: "nuxt/app/vitest.config.ts", Destination: "vitest.config.ts"},
			{Source: "nuxt/app/gitignore", Destination: ".gitignore"},
			{Source: "nuxt/app/app.vue.tmpl", Destination: "app/app.vue", Template: true},
			{Source: "nuxt/app/main.css", Destination: "app/assets/css/main.css"},
			{Source: "nuxt/app/app.test.ts.tmpl", Destination: "tests/app.test.ts", Template: true},
			{Source: "nuxt/app/Makefile", Destination: "Makefile"},
			{Source: "nuxt/app/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "frontend/app/AGENTS.md", Destination: "AGENTS.md"},
			{Source: "nuxt/app/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "frontend/app/docs-README.md", Destination: "docs/README.md"},
			{Source: "nuxt/production/eslint.config.mjs", Destination: "eslint.config.mjs"},
			{Source: "nuxt/production/prettierignore", Destination: ".prettierignore"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
		},
	}
}

// VueLibrary returns the built-in Vue component library proof bundle.
func VueLibrary() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/vue-library",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "vue/library/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "vue/library/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "vue/library/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "vue/library/vite.config.ts", Destination: "vite.config.ts"},
			{Source: "vue/library/vitest.config.ts", Destination: "vitest.config.ts"},
			{Source: "vue/library/gitignore", Destination: ".gitignore"},
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
			{Source: "vue/library/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "vue/library/Makefile", Destination: "Makefile"},
			{Source: "vue/library/index.ts.tmpl", Destination: "src/index.ts", Template: true},
			{Source: "vue/library/Greeting.vue.tmpl", Destination: "src/greeting.vue", Template: true},
			{Source: "vue/library/Greeting.test.ts", Destination: "src/__tests__/greeting.test.ts"},
		},
	}
}

// VueLean returns the reduced Vue application baseline.
func VueLean() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/vue-lean",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "vue/app/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "vue/app/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "vue/app/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "vue/app/vite.config.ts", Destination: "vite.config.ts"},
			{Source: "vue/app/vitest.config.ts", Destination: "vitest.config.ts"},
			{Source: "vue/app/index.html.tmpl", Destination: "index.html", Template: true},
			{Source: "vue/app/gitignore", Destination: ".gitignore"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "js/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "vue/app/main.ts", Destination: "src/main.ts"},
			{Source: "vue/app/env.d.ts", Destination: "src/env.d.ts"},
			{Source: "frontend/app/main.css", Destination: "src/main.css"},
			{Source: "vue/app/App.vue.tmpl", Destination: "src/app.vue", Template: true},
			{Source: "vue/app/App.test.ts.tmpl", Destination: "src/app.test.ts", Template: true},
			{Source: "frontend/app/AGENTS.md", Destination: "AGENTS.md"},
			{Source: "frontend/app/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "frontend/app/docs-README.md", Destination: "docs/README.md"},
			{Source: "frontend/app/prettierignore", Destination: ".prettierignore"},
			{Source: "frontend/app/eslint.config.mjs", Destination: "eslint.config.mjs"},
			{Source: "frontend/app/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "frontend/app/Makefile", Destination: "Makefile"},
		},
	}
}

// VueApp returns the comprehensive Vue modular-monolith baseline.
func VueApp() (fs.FS, templateengine.Blueprint) {
	source, blueprint := VueLean()
	blueprint.ID = "aruo/vue"
	blueprint.Files = appendFullstackFiles(blueprint.Files, "vue/full/package.json.tmpl")
	return source, blueprint
}

// NextLean returns the reduced Next.js application baseline.
func NextLean() (fs.FS, templateengine.Blueprint) {
	source, err := fs.Sub(templates, "templates")
	if err != nil {
		panic("embedded template subtree is invalid: " + err.Error())
	}
	return source, templateengine.Blueprint{
		ID:       "aruo/next-lean",
		Language: "typescript",
		Files: []templateengine.FileSpec{
			{Source: "next/app/README.md.tmpl", Destination: "README.md", Template: true},
			{Source: "next/app/package.json.tmpl", Destination: "package.json", Template: true},
			{Source: "next/app/tsconfig.json", Destination: "tsconfig.json"},
			{Source: "next/app/next.config.ts", Destination: "next.config.ts"},
			{Source: "next/app/vitest.config.ts", Destination: "vitest.config.ts"},
			{Source: "next/app/gitignore", Destination: ".gitignore"},
			{Source: "foundation/LICENSE.tmpl", Destination: "LICENSE", Template: true},
			{Source: "foundation/CHANGELOG.md", Destination: "CHANGELOG.md"},
			{Source: "foundation/ROADMAP.md", Destination: "ROADMAP.md"},
			{Source: "foundation/SECURITY.md.tmpl", Destination: "SECURITY.md", Template: true},
			{Source: "foundation/CONTRIBUTING.md", Destination: "CONTRIBUTING.md"},
			{Source: "foundation/CODE_OF_CONDUCT.md", Destination: "CODE_OF_CONDUCT.md"},
			{Source: "foundation/editorconfig", Destination: ".editorconfig"},
			{Source: "next/app/aruo.yaml.tmpl", Destination: "aruo.yaml", Template: true},
			{Source: "next/app/docs-README.md", Destination: "docs/README.md"},
			{Source: "foundation/issue-bug.yml", Destination: ".github/ISSUE_TEMPLATE/bug.yml"},
			{Source: "foundation/issue-feature.yml", Destination: ".github/ISSUE_TEMPLATE/feature.yml"},
			{Source: "foundation/pull-request.md", Destination: ".github/pull_request_template.md"},
			{Source: "js/library/dependabot.yml", Destination: ".github/dependabot.yml"},
			{Source: "foundation/pr-title.yml", Destination: ".github/workflows/pr-title.yml"},
			{Source: "foundation/release.yml", Destination: ".github/workflows/release.yml"},
			{Source: "js/library/release-please-config.json", Destination: "release-please-config.json"},
			{Source: "foundation/release-please-manifest.json", Destination: ".release-please-manifest.json"},
			{Source: "next/app/ci.yml", Destination: ".github/workflows/ci.yml"},
			{Source: "next/app/Makefile", Destination: "Makefile"},
			{Source: "next/app/layout.tsx.tmpl", Destination: "app/layout.tsx", Template: true},
			{Source: "frontend/app/main.css", Destination: "app/main.css"},
			{Source: "next/app/page.tsx.tmpl", Destination: "app/page.tsx", Template: true},
			{Source: "next/app/page.test.tsx.tmpl", Destination: "app/page.test.tsx", Template: true},
			{Source: "next/app/AGENTS.md", Destination: "AGENTS.md"},
			{Source: "next/app/eslint.config.mjs", Destination: "eslint.config.mjs"},
			{Source: "next/app/prettierignore", Destination: ".prettierignore"},
		},
	}
}

// NextApp returns the comprehensive Next.js modular-monolith baseline.
func NextApp() (fs.FS, templateengine.Blueprint) {
	source, blueprint := NextLean()
	blueprint.ID = "aruo/next"
	filtered := blueprint.Files[:0]
	for _, file := range blueprint.Files {
		if file.Destination != "AGENTS.md" && file.Destination != "aruo.yaml" && file.Destination != "package.json" && file.Destination != "README.md" {
			filtered = append(filtered, file)
		}
	}
	filtered = append(filtered,
		fullstackFile("fullstack/server/env.ts", "server/env.ts"),
		fullstackFile("fullstack/server/errors.ts", "server/errors.ts"),
		fullstackFile("fullstack/server/logger.ts", "server/logger.ts"),
		fullstackFile("fullstack/server/database.ts", "server/database.ts"),
		fullstackFile("fullstack/server/rate-limit.ts", "server/rate-limit.ts"),
		fullstackFile("fullstack/server/schema.ts", "server/db/schema.ts"),
		fullstackFile("fullstack/server/migration.sql", "server/db/migrations/0000_baseline.sql"),
		fullstackFile("fullstack/server/env.test.ts", "server/env.test.ts"),
		fullstackFile("fullstack/env.example", ".env.example"),
		fullstackFile("fullstack/compose.yaml", "compose.yaml"),
		fullstackFile("fullstack/drizzle.config.ts", "drizzle.config.ts"),
		fullstackFile("fullstack/AGENTS.md", "AGENTS.md"),
		fullstackFile("fullstack/docs-architecture.md", "docs/architecture.md"),
		fullstackFile("fullstack/docs-operations.md", "docs/operations.md"),
		fullstackFile("fullstack/dockerignore", ".dockerignore"),
		fullstackFile("next/full/Dockerfile", "Dockerfile"),
		fullstackFile("next/full/health-live.ts", "app/api/health/live/route.ts"),
		fullstackFile("next/full/health-ready.ts", "app/api/health/ready/route.ts"),
		fullstackFile("next/full/proxy.ts", "proxy.ts"),
		fullstackTemplate("fullstack/aruo.yaml.tmpl", "aruo.yaml"),
		fullstackTemplate("next/full/package.json.tmpl", "package.json"),
		fullstackTemplate("fullstack/README.md.tmpl", "README.md"),
	)
	blueprint.Files = filtered
	return source, blueprint
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
