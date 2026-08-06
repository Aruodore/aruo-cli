# Design principles

This section will illustrate how documentation-first work, opinionated defaults, native ecosystem respect, performance budgets, safe automation, user ownership, and evidence-based claims change concrete product decisions.

## Documentation first

Documentation is a versioned interface. User-visible behavior, examples, limitations, migrations, and architecture change with the product and are tested with it.

## Opinionated defaults

Defaults should produce a coherent supported path and minimize premature choices. Escape hatches are explicit, documented, owned, and time-bounded rather than hidden forks.

## Performance matters

Startup, steady-state time, memory, artifact size, and benchmark methodology are product requirements. Optimization follows measurement, not folklore.

## Developer experience is a feature

Clear help, useful errors, accessible output, safe previews, quick feedback, and local/CI parity are correctness concerns. A feature that users cannot discover or recover is incomplete.

## Automation over repetition

Repeated work should become deterministic automation. Automation must state its inputs, effects, permissions, evidence, and rollback.

## Simple APIs

Public contracts should be small, typed, stable, and unsurprising. Internal complexity is not exported merely because it exists.

## Composable architecture

Capabilities compose through explicit contracts and dependency rules. Core policy is separated from language/provider adapters, and cycles are rejected rather than managed by convention.

## Quality over quantity

A small catalog with owners, tests, upgrades, and support policy is better than many stale templates. Checks measure outcomes and risk, not ceremony.

## Long-term maintainability

Dependencies, formats, decisions, and deprecations carry ownership and migration paths. Accepted decisions are superseded, never rewritten to erase history.

## Native ecosystems deserve respect

Aruo standardizes outcomes while preserving language-native structure and tooling. A Go project should feel like Go; a Python project should feel like Python.

## Trust is a feature

Offline operation, least privilege, signed artifacts, secret redaction, plugin isolation, and reversible mutations are defaults. AI may propose; deterministic engines verify; people approve consequential effects.

## Every repository should teach

Structure, commands, decisions, examples, failure modes, and contribution routes should be understandable without private institutional memory.

## Related pages

- [Vision](vision.md)
