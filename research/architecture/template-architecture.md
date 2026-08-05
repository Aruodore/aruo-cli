# Template architecture

## Philosophy

A template is a tested, versioned composition of capabilities—not a repository snapshot. Small primitives reduce combinatorial forks; blueprints select compatible primitives; semantic operations update structured files. Every option must change a meaningful contract and have a test fixture.

## Layers

1. **Foundation:** community health, editor, Git, docs skeleton, task vocabulary.
2. **Language:** native package/build/test/format/type conventions.
3. **Workload:** library, CLI, service, frontend, research, benchmark.
4. **Capabilities:** docs site, release, container, API schema, database, GPU, model card.
5. **Organization overlay:** ownership, domains, approved CI runners, policy exceptions.

Layers merge through declared semantic keys. They do not overwrite whole `package.json`, `pyproject.toml`, workflow, or README documents.

## Initial catalog

| Blueprint | Default spine | Distinct evidence |
|---|---|---|
| Go library | module, native tests/examples, golangci-lint | API compatibility and multi-platform test |
| Go/Rust CLI | Cobra-equivalent or minimal parser, completion, man docs | golden CLI tests and release archives |
| Python package | `src` layout, uv, Ruff, Pyright/mypy policy, pytest | wheel/sdist clean-install matrix |
| Rust crate | Cargo conventions, rustfmt/clippy, docs.rs metadata | MSRV and feature matrix |
| TypeScript library | strict TS, explicit exports, ESM-first, Vitest | packed-artifact/type/API checks |
| Node service | layered service, schemas, health/telemetry | integration/contract/container checks |
| React/Vue app | feature slices, accessibility and E2E | bundle/performance/accessibility budgets |
| CV/ML project | configs, cards, external artifact manifests | reproducibility, slice analysis, quality/performance pair |
| Benchmark project | workload schema, baselines, raw results | environment-normalized reproducibility |

Computer Vision and ML share a research foundation but remain separate profiles when their data/evaluation contracts differ.

## Artifact format

Each artifact declares ID, version, Aruo protocol range, license, maintainers, inputs with validation/help, provided/required/conflicting capabilities, files, semantic operations, checks, migrations, test fixtures, digest, and signature. Executable hooks are not templates; they are permissioned plugins or explicit plan processes.

## Qualification

Every blueprint is tested as a matrix: minimal/default/maximal supported composition; Linux/macOS/Windows; render determinism; native build/test/docs/package; upgrade from each supported artifact version; dirty-repository conflict fixtures; license/security scans. Support is bounded—Aruo publishes a compatibility window and retires stale variants rather than accumulating choices forever.
