# Technical specification

## Scope and requirements

The v1 system MUST create and adopt repositories, resolve layered configuration, inventory current state, generate a deterministic plan, apply approved changes transactionally, run policy checks, manage provenance, and emit human plus stable machine output. It MUST work offline after artifacts are installed.

It SHOULD support Go libraries/CLIs, Python packages, Rust crates/CLIs, TypeScript libraries, Node services, and React/Vue applications at v1. ML/CV and benchmark profiles may graduate after reproducibility gates are proven.

## Core domain model

- `Project`: identity, lifecycle, packages, owners, targets.
- `Capability`: desired outcome such as docs, CI, release, benchmarks, security.
- `PolicyPack`: versioned requirements, checks, fixes, severity, rationale.
- `Blueprint`: compatible capabilities and semantic operations; replaces monolithic template as primary abstraction.
- `Observation`: normalized fact from repository inspection.
- `Finding`: expected versus observed state, evidence, severity, remediation.
- `Plan`: ordered operations, preconditions, effects, risk, rollback, provenance.
- `Operation`: add/merge/patch/move/run; typed and idempotency-declared.
- `Artifact`: template/plugin/policy with origin, version, digest, signature, compatibility.
- `Evidence`: command/tool version, inputs, output summary, timestamp, commit, environment.
- `Exception`: scope, rationale, owner, expiry, approval.

## Command execution contract

All mutating workflows use `discover → resolve → plan → approve → apply → verify → record`. Planning MUST NOT execute project code. Application stages writes in a repository-local temporary area, validates preconditions and output, then commits file operations atomically where the filesystem permits. External effects are separate, explicitly identified steps and are never implied by file generation.

## Functional boundaries

- No network during `audit`, `plan`, or local checks unless a named check declares it and the user permits it.
- No secrets in project configuration, plans, logs, lockfiles, or evidence bundles.
- Config, plans, findings, plugin protocol, and artifact manifests have JSON Schemas and explicit format versions.
- Exit codes are stable: `0` success/conformant, `1` operational failure, `2` usage/config, `3` findings above threshold, `4` conflicts, `5` trust/security refusal.
- JSON output is versioned and contains no terminal prose. Diagnostics go to stderr.

## Quality attributes

- warm help/version under 50 ms target; cold local read-only command under 150 ms p50 on reference hardware;
- reproducible native binaries for tier-1 platforms;
- deterministic planning given repository snapshot, config, and artifact lock;
- cancellation and crash safety; no partial silent mutation;
- forward-compatible config parsing where unknown fields fail with useful location and suggestion;
- Windows, macOS, and Linux path/process tests; no shell assumed for core operations.

## Non-goals for v1

Hosted dashboards, arbitrary code generation, remote execution, deployment provider abstraction, full AST refactoring across all languages, and untrusted in-process plugins are deferred.
