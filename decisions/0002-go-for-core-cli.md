# ADR-0002: Use Go for the core CLI

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead
- Related: `research/technology/rust-first-stack-evaluation.md`

## Context

Aruo needs fast startup, cross-platform standalone binaries, safe concurrency, strong tooling, and an approachable contributor/runtime model. Phase 0 initially recommended Rust subject to a spike; the repository design selected Go as the accepted direction.

## Considered options

- Rust: strongest type/memory guarantees and performance control; higher contributor and implementation complexity.
- Go: simple deployment/concurrency, fast builds, mature CLI ecosystem, good standard library; garbage collection and a less expressive type system.
- TypeScript/Python: excellent ecosystem access but runtime distribution/startup/dependency costs for a cross-language control plane.

## Decision

Use Go for the CLI and core modular monolith. Keep schemas/protocols language-neutral and plugins out of process.

## Consequences

Use Go-native profiling and bounded concurrency. Public code remains under `internal/` unless stability is deliberate. Performance budgets watch GC/allocation behavior. The historical Rust recommendation remains in research.

## Validation

Before implementation expands, a vertical slice must meet startup, binary-size, YAML round-trip, plugin protocol, and cross-compilation acceptance criteria.

