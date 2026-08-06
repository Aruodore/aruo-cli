# Project structure

```text
aruo/
├── .design/                  internal product memory; never distributed downstream
├── .devcontainer/            reproducible editor/container development environment
├── .github/                  collaboration, ownership, security, CI, and releases
├── cmd/aruo/                 process entry point and dependency wiring only
├── internal/buildinfo/       immutable linker-provided build identity
├── internal/cli/             invocation composition and error boundary
│   ├── command/              constructor-built Cobra presentation adapters
│   └── iostreams/            explicit process streams and future capabilities
├── internal/templateengine/  pure bounded renderer and embedded built-in bundles
├── internal/catalog/         qualified template discovery and built-in entries
├── internal/create/          catalog-neutral creation service and atomic writer
├── internal/doctor/          read-only checks, scoring, evidence, and recommendations
├── pkg/                      intentionally public Go APIs; empty by default
├── api/                      versioned schemas and plugin protocol contracts
├── blueprints/               first-party composable project definitions
├── templates/                text/document templates used by blueprints
├── plugins/                  first-party plugin manifests and SDK fixtures
├── policies/                 machine-readable engineering policy packs
├── configs/                  development/tool configuration fragments
├── docs/                     user and contributor documentation source
├── examples/                 executable consumer workflows, later tested
├── tests/                    cross-module integration/conformance fixtures
├── benchmarks/               harness definitions and governed results
├── scripts/                  thin documented automation, no business logic
├── decisions/                immutable architecture decision records
├── rfcs/                     substantial proposals and RFC process artifacts
├── research/                 permanent evidence, evaluations, and open questions
├── website/                  docs-site application/configuration, not content truth
├── assets/                   source-controlled shared brand/diagram assets
├── go.mod                    module identity and Go compatibility/toolchain
├── Makefile                  discoverable local quality task facade
├── .golangci.yml             formatting and curated static-analysis policy
└── .goreleaser.yaml          cross-platform release artifact specification
```

## Rules

- Do not create empty directories for appearance. During design, each future implementation directory contains a README defining its boundary.
- Add a top-level directory only when it has a distinct audience, lifecycle, ownership, or tooling boundary that cannot be expressed in an existing directory.
- `cmd/` wires dependencies; reusable behavior belongs in `internal/`.
- `pkg/` is not a dumping ground. Moving code there creates a compatibility promise and requires an ADR/API review.
- User docs live in `docs/`; `website/` contains presentation/build machinery only.
- `research/` preserves evidence and rejected options; `decisions/` states accepted choices; `rfcs/` proposes change. They are not interchangeable.
- `.design/` preserves informal product thinking and meeting context. It is committed in Aruo but excluded from generated repositories, packages, release archives, and the public docs build.
- Generated outputs, caches, dependencies, secrets, raw datasets, binaries, and benchmark scratch files are never committed.
- Distribution allowlists MUST be used instead of relying only on ignore patterns; tests MUST fail if `.design/` appears in any downstream artifact.
- Language-native colocated tests remain beside code; `tests/` holds cross-package, compatibility, E2E, and fixture suites.

The executable shell above is implemented; future domain package names remain provisional until their implementation RFCs. [Architecture](README.md) and the [CLI application architecture](cli-application.md) define responsibility boundaries that package names must preserve.
