# ADR-0001: Build a local-first repository control plane

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

Existing tools optimize generation or one lifecycle stage. Repository source, native tools, and CI remain the durable source of engineering truth. Requiring a service would add privacy, availability, cost, and lock-in risks before proving the local workflow.

## Considered options

- Template generator only: simple, but cannot explain drift or maintain projects.
- Hosted control plane first: enables fleet views, but weakens offline ownership and expands scope.
- Local lifecycle engine with optional future service: preserves ownership and composes with CI.

## Decision

Build a local-first CLI that models intent, inspects state, creates deterministic plans, verifies evidence, and records provenance. Hosted collaboration remains optional.

## Consequences

Core workflows work offline and schemas must be portable. Cross-repository history is deferred. Local security and crash recovery become first-class responsibilities.

## Validation

The v0.1 acceptance suite runs without accounts/network and produces identical plans from identical inputs.

