# ADR-0007: Bootstrap with native Go tools, Make, GitHub Actions, and release PRs

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

The repository needs one reproducible contributor loop, current Go quality checks, cross-platform CI, dependency maintenance, conventional release intent, and verifiable artifacts before business logic. Extra runtimes and overlapping task systems increase bootstrap cost.

## Considered options

- Make plus Go-native tools: broadly available, transparent, small abstraction; Windows contributors may rely on CI/devcontainer or install Make.
- Task/Mage/Just: stronger cross-platform ergonomics in places, but adds another binary or Go code before need is proven.
- Python pre-commit: rich ecosystem and staged-file speed, but adds Python and remains bypassable.
- semantic-release: maximal automation, but releases directly from commit interpretation with less review than a release PR.
- release-please plus GoReleaser: reviewable versions/changelog and strong Go artifact pipeline.

## Decision

Use Go 1.26, Make as the single task facade, gofmt/goimports, curated golangci-lint v2, native tests/race/coverage, GitHub Actions pinned by SHA, Dependabot, Conventional Commit PR titles, release-please release PRs, and GoReleaser artifacts. Do not require pre-commit hooks or a second task runner.

## Consequences

CI is the enforcement boundary. Contributors install one external linter for the full local gate; GoReleaser is release-only. Windows CI invokes native commands directly. Tool upgrades are visible PRs. Release automation requires protected environment/token setup.

## Validation

All configured checks must pass on the empty bootstrap, action pins must be immutable, snapshot release must build without publishing, and the contributor setup must remain documented and reproducible.

