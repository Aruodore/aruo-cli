# Architecture decision records

ADRs capture one consequential decision, its context, considered options, and consequences. The format follows the established Nygard core and MADR’s emphasis on alternatives. Status is `proposed`, `accepted`, `deprecated`, or `superseded by ADR-NNNN`.

Use `0000-template.md`. Never recycle numbers or rewrite an accepted decision’s outcome; create a superseding ADR. Small factual/link corrections are allowed without changing meaning.

## Index

| ADR | Decision | Status |
|---|---|---|
| [0001](0001-local-first-repository-control-plane.md) | Build a local-first repository lifecycle control plane | Accepted |
| [0002](0002-go-for-core-cli.md) | Implement the core CLI in Go | Accepted |
| [0003](0003-out-of-process-plugins.md) | Run third-party plugins out of process | Accepted |
| [0004](0004-composable-blueprints.md) | Compose constrained blueprints instead of monolithic templates | Accepted |
| [0005](0005-outcome-based-repository-standards.md) | Enforce outcomes while preserving native layouts | Accepted |
| [0006](0006-apache-2-license.md) | License Aruo under Apache-2.0 | Accepted |
| [0007](0007-repository-toolchain-and-automation.md) | Use native Go tools, Make, GitHub Actions, and release PRs | Accepted |

The current sequence continues with [ADR-0008: Constructor-built Cobra presentation layer](0008-constructor-built-cobra-cli.md), [ADR-0009](0009-pure-standard-library-template-renderer.md) (the pure standard-library template renderer), and [ADR-0010](0010-charm-v2-terminal-ux-stack.md) (the Charm v2 terminal UX stack behind the `internal/tux` adapter boundary).

## Why Go

Go produces portable static binaries, has fast startup and compilation, a strong standard library, simple concurrency, first-class cross-compilation, and an established CLI ecosystem. It lowers operational and contributor complexity for a tool that orchestrates other tools. The cost is less expressive type modeling and fewer memory-safety guarantees than Rust; ADR-0002 records why simplicity and delivery risk win here.

## Why a CLI

Repositories live in local files and CI. A CLI works offline, composes with scripts, exposes deterministic exit/output contracts, and keeps user code under user control. Future APIs and IDEs share the application layer rather than replace it.

## Why plugins and templates

Plugins isolate ecosystem/provider variability behind versioned capabilities. Templates provide human-readable starting material, while semantic operations and provenance make them maintainable. Neither is allowed to bypass planning, trust, or policy boundaries.

## Why Apache-2.0

Apache-2.0 permits broad commercial and open-source reuse while providing an explicit contributor patent grant and patent-termination protection. It is longer than MIT and imposes notice obligations, but those costs are appropriate for developer infrastructure intended to become an ecosystem foundation. This is a project recommendation, not legal advice.
