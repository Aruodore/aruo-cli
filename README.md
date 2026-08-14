# Aruo

[![CI](https://github.com/Aruodore/aruo-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Aruodore/aruo-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Aruodore/aruo-cli?include_prereleases)](https://github.com/Aruodore/aruo-cli/releases)
[![License](https://img.shields.io/github/license/Aruodore/aruo-cli)](LICENSE)

**A CLI that installs a durable AI engineering contract into real applications and audits whether production responsibilities remain explicit.**

> **Status: pre-1.0.** `aruo init` is implemented on the development branch after the v0.1.0 prerelease; expect breaking changes before 1.0. Existing `create` templates remain available but are no longer the product's architectural center.

```sh
git clone https://github.com/Aruodore/aruo-cli.git aruo && cd aruo
go build -o "$HOME/.local/bin/aruo" ./cmd/aruo
cd your-existing-application
aruo init --dry-run
aruo init --yes && aruo doctor
```

## Why Aruo

AI can produce working software quickly, but it cannot infer every application's architecture, threat model, operating environment, or definition of done. Aruo installs explicit rules that make those expectations available to AI agents and humans inside the repository. The installer owns and can eventually update the rules; the developer owns the application and its production intent. Doctor makes drift and unresolved responsibilities visible without pretending the application is certified production-ready.

## What works today

- **`aruo init`**: detects an existing repository's ecosystem, framework, and package manager; plans or installs `AGENTS.md`, a versioned contract and rules under `.aruo/`, managed-file hashes, and application-owned `aruo.yaml`. It adds no runtime dependency and refuses every existing target rather than overwriting it.
- **`aruo doctor`**: verifies the installed contract's integrity, audits application intent, and independently scores repository health against `aruo.repository-health/v1`.
- **`aruo create`**: the existing project-template catalog remains available during the pivot. New production guidance belongs in the installable contract, not increasingly large templates.
- Accessible and scriptable by design: a line-oriented fallback adapter, `NO_COLOR`/`--color`/`--motion` support, and clean non-interactive behavior for `init`, `create`, and `doctor`.

Everything else referenced elsewhere in this repository's docs (`inspect`, `check`, `adopt`, `plan`, `apply`, plugins, migrations, `aruo config`) is designed but not implemented. See [Roadmap](#roadmap).

## Installation

Aruo isn't on any package manager yet, but a tagged release with prebuilt binaries is available.

### Download a release binary

Grab the archive for your platform from the [v0.1.0 release](https://github.com/Aruodore/aruo-cli/releases/tag/v0.1.0) (Linux/macOS/Windows × amd64/arm64), verify it against `checksums.txt`, then extract and put `aruo` on your `PATH`. Each archive ships alongside an SBOM and a [GitHub build provenance attestation](https://github.com/Aruodore/aruo-cli/attestations).

```sh
curl -sSfL -o aruo.tar.gz https://github.com/Aruodore/aruo-cli/releases/download/v0.1.0/aruo_0.1.0_linux_amd64.tar.gz
tar xzf aruo.tar.gz
install aruo "$HOME/.local/bin/aruo"
```

`aruo version` on this binary prints the real tagged version, `aruo version 0.1.0`. GoReleaser stamps that in at build time; `go install` and local builds below don't, so they report `dev` instead.

### `go install`

```sh
go install github.com/aruodore/aruo-cli/cmd/aruo@latest
```

Works now that the module is public and tagged. Reports `aruo version dev`, like any source build.

### Build from source

**Prerequisites:** Go 1.26 or newer (the repository pins 1.26.5 via `go.mod`'s `toolchain` directive; a modern `go` command fetches it automatically). No CGO, no other system dependencies.

```sh
git clone https://github.com/Aruodore/aruo-cli.git aruo
cd aruo
go build -o "$HOME/.local/bin/aruo" ./cmd/aruo
```

- **Installed to:** wherever you point `-o`; the example above uses `$HOME/.local/bin/aruo`. Make sure that directory is on `PATH`.
- **Supported platforms:** Linux, macOS, and Windows all build and pass the full test suite in CI (`ubuntu-latest`, `macos-latest`, `windows-latest`); `amd64` and `arm64` both build via the pinned GoReleaser config, though only Linux/`amd64` has been run interactively during development.
- **Verify:** `aruo version` prints `aruo version dev`, matching every source build.
- **Upgrade:** `git pull && go build -o "$HOME/.local/bin/aruo" ./cmd/aruo`. There's no update mechanism beyond rebuilding.
- **Uninstall:** delete the binary you built (e.g. `rm "$HOME/.local/bin/aruo"`). Nothing else is installed on your system.

### Not yet available

None of these exist today. They're listed so you don't go looking for something that isn't there yet, not as a promise of when they'll land.

| Method | Status |
| --- | --- |
| Homebrew | Not configured. No `brews:` block in the release config. |
| Scoop / Winget | Not configured. |
| AUR / Nix | Not configured. |
| Install script (`curl \| sh`) | Not implemented. |

## Quick start

```sh
aruo version
cd your-existing-application
aruo init --dry-run
aruo init --yes
aruo doctor
```

That is inspect → install the contract → audit the repository, using only commands that exist. Run `aruo init` without `--yes` for confirmation, or use `--no-input --yes` in automation.

## Example workflow

A real, unedited run of the commands above:

```text
$ aruo create my-library --template go-library
✓ Created go-library with 24 files at ./my-library

Next steps:
  cd my-library
  go test ./...
  git init && git add .

$ cd my-library && go test ./...
ok  	my-library	0.002s

$ aruo doctor
Repository health: 99/100 (A)
/path/to/my-library
Category       Score
completeness   20/20
documentation  20/20
ci             15/15
tests          15/15
license        10/10
security       10/10
github          9/10
Recommendations:
  - Missing CODEOWNERS
    Declare real maintainers for review routing; do not add placeholders.
```

The one missing point is expected: a freshly generated project has no named maintainer yet, and Aruo won't invent one.

## Commands

Run `aruo <command> --help` for the authoritative, always-current flag list; this section is a summary, not a substitute.

### `aruo init [repository]`

Installs Aruo's AI engineering contract into an existing repository. It detects local stack evidence but does not execute project code, install dependencies, or modify application source.

```sh
aruo init --dry-run                 # inspect the exact file plan
aruo init --yes                     # initialize the current repository
aruo init ./application --format json --yes
```

Managed files are `AGENTS.md` and `.aruo/**`. `aruo.yaml` is created once and then owned by the application. Initialization refuses collisions, uses exclusive commit operations, and rolls back files it created when a commit cannot complete. Updating an existing installation is intentionally deferred to `aruo update`; rerunning `init` does not overwrite it.

### `aruo create [name-or-path]`

Creates a new project from a catalog template. Refuses a destination with real content in it; an existing but empty directory (most commonly `.`) is fine.

**Key flags:**

- `--template <id>`: skip the picker
- `--language` / `--kind`: filter the catalog
- `--module`: Go module path, npm package name, or PyPI name, depending on the template. Defaults to the project's own name; pass a real import path like `github.com/you/name` to override it
- `--description`, `--author`, `--license`: metadata for the generated project
- `--set key=value`: template variables
- `-y` / `--yes`: accept the confirmation
- `--no-input`: fail instead of prompting, for CI

```sh
aruo create                                   # fully interactive, guided
aruo create my-library --template go-library                       # --module defaults to my-library
aruo create my-tool --template python-library --no-input --yes     # --module defaults to my-tool
aruo create . --name my-tool --template go-library --module github.com/you/my-tool --no-input --yes
```

**Cancellation:** Ctrl+C once cancels cleanly. Writes are staged and only committed at the end, so no partial project is left behind. A second Ctrl+C forces immediate exit.

### `aruo doctor [repository]`

Scores a repository's engineering health and, when `aruo.yaml` declares production intent, audits its capability claims and unresolved responsibilities. Intent findings are separate and do not alter the versioned 100-point score. Defaults to the current directory when no path is given.

Doctor statically verifies supported evidence conventions such as complete npm quality scripts, framework build/type-check scripts, test files, dependency-audit CI steps, committed migrations, and health-route pairs. It does not run repository code. Runtime and provider-dependent evidence remains visibly `DECLARED`, never overstated as verified.

**Key flags:** `--format human|json`, `--minimum-score <0-100>` (default `80`). Doctor exits `3` when the score is below the threshold or the intent manifest has blocking findings, making it useful as a CI gate.

```sh
aruo doctor                              # score the current directory
aruo doctor ./some-project --format json
aruo doctor . --minimum-score 90
```

### `aruo version`

Prints the running build's version: the real tagged version on a release binary, `dev` on anything built locally. See [Installation](#installation).

### `aruo completion [bash|zsh|fish|powershell]`

Generates a shell completion script (Cobra's standard mechanism), including dynamic completion for `--template`/`--language`/`--kind`.

### Exit codes

Only what's actually implemented: `0` success, `1` operational failure, `3` `doctor` findings (score below `--minimum-score` or blocking production intent), `130`/`143` interrupted by Ctrl+C/SIGTERM.

## Generated project structure

`aruo create my-library --template go-library` produces:

```text
my-library/
├── .github/
│   ├── ISSUE_TEMPLATE/{bug.yml,feature.yml}
│   ├── workflows/{ci.yml,pr-title.yml,release.yml}
│   ├── dependabot.yml
│   └── pull_request_template.md
├── docs/README.md
├── aruo.yaml                  # template provenance and explicit production intent
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── go.mod
├── LICENSE
├── Makefile
├── my_library.go
├── my_library_test.go
├── README.md
├── release-please-config.json
├── .release-please-manifest.json
├── ROADMAP.md
└── SECURITY.md
```

Every template's exact file set differs by ecosystem (an npm-based template gets `package.json`/CI steps for `npm test`, for example) but follows the same shape: a working CI workflow, tests, and governance/security files, never placeholders.

## Configuration

There's no `aruo config` command or resolved project configuration today. Doctor reads only the versioned provenance and `intent.capabilities` contract from `aruo.yaml`; it does not treat the file as runtime configuration. The active configuration surface remains CLI flags and four environment variables:

| Variable | Effect |
| --- | --- |
| `NO_COLOR` | Disable color output (any value, including empty) |
| `ARUO_ACCESSIBLE` | Force the line-oriented accessible adapter |
| `ARUO_MOTION=never` | Disable animation |
| `ARUO_NO_INPUT` | Disable prompting; fail instead of asking |

The `aruo.yaml` written into every generated project records template provenance and explicit production intent. `aruo doctor` reads and audits those fields only. [`docs/configuration/README.md`](docs/configuration/README.md) describes a considerably larger planned configuration system (`aruo.yaml` as live config, `aruo config explain`, org policy, profiles); that configuration system does not exist yet.

## Documentation

- [`docs/README.md`](docs/README.md) — documentation index
- [`docs/architecture/README.md`](docs/architecture/README.md) — system architecture
- [`docs/architecture/create-command.md`](docs/architecture/create-command.md) — how `create` actually works
- [`docs/cli/README.md`](docs/cli/README.md) — command surface and terminal UX contract
- [`docs/cli/copy-style-guide.md`](docs/cli/copy-style-guide.md) — prompt/error/status wording rules
- [`docs/templates/README.md`](docs/templates/README.md) — the template catalog

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md), the [RFC process](rfcs/README.md), and the [Code of Conduct](CODE_OF_CONDUCT.md). Security reports go through [SECURITY.md](SECURITY.md), not public issues.

## Roadmap

`init`, `create`, and `doctor` are real. Managed contract updates, organization policy, plugins, and resolved application configuration are not built. See [ROADMAP.md](ROADMAP.md) for what's next.

## License and security

Aruo is licensed under [Apache-2.0](LICENSE). Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).
