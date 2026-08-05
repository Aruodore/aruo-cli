# Historical technology-stack evaluation

## Core decision: Rust

Build the CLI/core as a Rust modular monolith. Cargo provides native package/workspace conventions; Rust supports predictable standalone binaries, strong types for versioned plans/config, safe concurrency, and good cross-platform performance. The tradeoff is contributor learning curve and longer compile times. Go is the strongest alternative and should be reconsidered after a spike measuring binary size, startup, YAML comment preservation, plugin protocol ergonomics, and cross-platform packaging.

## Components

| Concern | Recommendation | Reason / constraint |
|---|---|---|
| CLI | clap | generated help/completion, typed hierarchy; keep domain independent |
| async | Tokio only at adapter boundary | concurrent I/O; avoid infecting pure planning code |
| config | serde + schema generation + comment-preserving YAML CST library selected by spike | typed validation plus round-trip edits |
| diagnostics | miette-style structured diagnostics | source spans and actionable errors; stable own diagnostic model |
| Git | invoke system Git first | honors user credential/config behavior; libgit only for proven gaps |
| HTTP | rustls-based client | portable TLS; network always explicit/cacheable |
| protocol | JSON Lines over stdio, JSON Schema | language-neutral, inspectable; benchmark before binary protocol |
| templates | restricted MiniJinja-compatible expressions | familiar syntax without arbitrary execution |
| structured edits | dedicated JSON/TOML/YAML/XML adapters; language AST tools as plugins | preserve semantics/comments and limit core breadth |
| local store | files/SQLite only when query need is proven | keep repository state transparent |
| docs | VitePress + local search | fast static docs; framework-neutral Markdown source |
| CI | GitHub Actions initially | matches current handbook; pin actions and abstract policy |
| releases | GoReleaser-equivalent Rust workflow or cargo-dist evaluated; GitHub immutable release, checksums, signatures, SBOM/provenance | build once and verify artifacts |

## Native tool defaults

- Go: gofmt/go vet plus curated golangci-lint.
- Python: uv, Ruff, pytest, and Pyright or mypy chosen per compatibility needs.
- Rust: cargo fmt/clippy/test/doc; cargo-deny/audit policy.
- JS/TS: one pinned package manager, strict TypeScript, Biome for greenfield when ecosystem plugins are unnecessary; ESLint + Prettier where framework/plugin semantics require them.

Aruo wraps these tools with a normalized task/result contract but retains native config and direct invocation. It MUST print the underlying command and version in verbose/evidence output.

## Performance budget

No network, plugin discovery, repository-wide hash, or runtime initialization for `--help`/`--version`. Lazy-load adapters; cache content hashes keyed by tool/config/version; bounded concurrency by resource class; stream subprocess output; never cache secrets or nondeterministic checks. Target compressed binary under 25 MB initially and measure rather than promise. Tier 1: Linux x86_64/aarch64, macOS x86_64/arm64, Windows x86_64; other targets community-supported until CI and release signing exist.
