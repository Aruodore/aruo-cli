# Engineering decisions

This page is the index and reader-friendly summary. Normative decisions live as immutable records in [`decisions/`](decisions/README.md).

| ADR | Decision | Status |
|---|---|---|
| [0001](decisions/0001-local-first-repository-control-plane.md) | Build a local-first repository lifecycle control plane | Accepted |
| [0002](decisions/0002-go-for-core-cli.md) | Implement the core CLI in Go | Accepted |
| [0003](decisions/0003-out-of-process-plugins.md) | Run third-party plugins out of process | Accepted |
| [0004](decisions/0004-composable-blueprints.md) | Compose constrained blueprints instead of monolithic templates | Accepted |
| [0005](decisions/0005-outcome-based-repository-standards.md) | Enforce outcomes while preserving native layouts | Accepted |
| [0006](decisions/0006-apache-2-license.md) | License Aruo under Apache-2.0 | Accepted |
| [0007](decisions/0007-repository-toolchain-and-automation.md) | Use native Go tools, Make, GitHub Actions, and release PRs | Accepted |

## Why Go

Go produces portable static binaries, has fast startup and compilation, a strong standard library, simple concurrency, first-class cross-compilation, and an established CLI ecosystem. It lowers operational and contributor complexity for a tool that orchestrates other tools. The cost is less expressive type modeling and fewer memory-safety guarantees than Rust; ADR-0002 records why simplicity and delivery risk win here.

## Why a CLI

Repositories live in local files and CI. A CLI works offline, composes with scripts, exposes deterministic exit/output contracts, and keeps user code under user control. Future APIs and IDEs share the application layer rather than replace it.

## Why plugins and templates

Plugins isolate ecosystem/provider variability behind versioned capabilities. Templates provide human-readable starting material, while semantic operations and provenance make them maintainable. Neither is allowed to bypass planning, trust, or policy boundaries.

## Why Apache-2.0

Apache-2.0 permits broad commercial and open-source reuse while providing an explicit contributor patent grant and patent-termination protection. It is longer than MIT and imposes notice obligations, but those costs are appropriate for developer infrastructure intended to become an ecosystem foundation. This is a project recommendation, not legal advice.

## Creating future decisions

Use [`decisions/0000-template.md`](decisions/0000-template.md). ADRs address one consequential, difficult-to-reverse choice. Accepted ADRs are not edited to match later reality; a new ADR supersedes them.
