# ADR-0008: Constructor-built Cobra presentation layer

- Status: Accepted
- Date: 2026-08-04
- Owners: Aruo maintainers

## Context

Aruo needs mature command parsing and help today without allowing a CLI framework to shape the future domain architecture. Cobra is established across the Go CLI ecosystem, but its common generator pattern uses package globals and `init()` registration. Those choices obscure dependencies and create shared state as a command tree grows.

## Decision

Use Cobra v1.10.2 inside `internal/cli/command` only. Build a fresh tree through constructors for each invocation. Keep `cmd/aruo` as the process composition boundary, inject streams and build metadata, return errors through `RunE`, and render them once in `internal/cli`. Application and domain packages must not import Cobra.

Use ordinary constructor injection rather than a dependency-injection framework. Load configuration lazily after parsing for commands that require it. Use `log/slog` behind an injected logger.

## Consequences

- Command tests are isolated and parallel-safe.
- Domain services remain reusable by future interfaces.
- Dependency wiring is explicit and searchable.
- Some Cobra examples and generators cannot be copied directly.
- Constructors require modest deliberate wiring as commands are added.
- Help text remains a public compatibility surface and must be reviewed.

## Alternatives considered

- Standard `flag`: minimal but insufficient for a growing nested command system and generated help.
- A custom parser: maximum control with unjustified maintenance and compatibility cost.
- Global Cobra commands with `init()`: less initial wiring but hidden mutation and poor isolation.
- Kong, urfave/cli, and other frameworks: capable, but Cobra has stronger alignment with the Go tools studied and can remain contained behind Aruo's boundary.
