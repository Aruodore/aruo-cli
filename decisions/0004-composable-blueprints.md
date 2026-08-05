# ADR-0004: Compose constrained blueprints

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

Monolithic conditional templates accumulate untestable combinations and own files too coarsely. Pure code generators allow arbitrary effects. Native initializers remain authoritative for ecosystem structure.

## Considered options

- Repository snapshots: transparent but update-hostile.
- General executable generators: flexible but unsafe and difficult to reason about.
- Capability blueprints plus restricted templates and semantic operations: more engineering, but testable and upgradeable.

## Decision

Compose foundation, language, workload, capability, and organization layers. Render text with a restricted deterministic engine; edit structured documents semantically; expose commands as explicit plan operations.

## Consequences

Artifact manifests need compatibility/conflict declarations and comprehensive fixtures. Users retain files; provenance is granular. Some transformations will intentionally stop for human conflict resolution.

## Validation

Qualify supported composition boundaries and measure upgrade conflict rates on realistically edited repositories.

