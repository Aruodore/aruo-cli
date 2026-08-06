# Contributing to Aruo

Thank you for helping build durable open-source engineering infrastructure. During Phase 0, changes should improve research, requirements, schemas, ADRs, RFCs, fixtures, or documentation. Production implementation starts only after its RFC is accepted.

## Before you begin

Read [VISION.md](VISION.md), [PHILOSOPHY.md](PHILOSOPHY.md), [the architecture guide](docs/architecture/README.md), and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). Search issues/RFCs before opening a large proposal. Security reports follow [SECURITY.md](SECURITY.md), not public issues.

## Development setup

Install the versions in `.tool-versions` (Go and golangci-lint are required; GoReleaser and Syft are release-only), then run:

```sh
make bootstrap
make check
```

Use `make help` for the supported task list and [`docs/development/toolchain.md`](docs/development/toolchain.md) for setup, container use, and troubleshooting. CI invokes equivalent native commands on Linux, macOS, and Windows.

## Workflow

1. Discuss substantial or compatibility-affecting work through the [RFC process](RFC_PROCESS.md).
2. Branch from current `main`: `type/short-description`, where type is `feature`, `fix`, `docs`, `refactor`, `test`, `perf`, or `chore`.
3. Make one cohesive change with tests/evidence and documentation.
4. Add release intent when user-visible behavior changes.
5. Open a draft PR early for risky work; complete the PR template before review.

Commits SHOULD follow `type(scope): imperative summary` using Conventional Commit types for readability. Maintainers squash when appropriate; commit syntax is not the sole source of release truth.

## Pull requests and review

PRs explain problem, approach, alternatives, user impact, risk, verification, compatibility, docs, and rollback. Generated or benchmark changes include reproducible commands. Authors respond to review without erasing useful discussion. At least one authorized maintainer approves; security, public API, protocol, policy, and architecture changes require the relevant owner. Authors do not merge their own high-risk changes without another maintainer.

Review evaluates correctness, simplicity, contracts, security, cross-platform behavior, accessibility, performance, tests, documentation, and maintainability. Be specific and kind; distinguish blocking concerns from suggestions.

## Standards

Follow [STYLE_GUIDE.md](STYLE_GUIDE.md), [TESTING.md](TESTING.md), and [DOCUMENTATION.md](DOCUMENTATION.md). Never commit credentials, private data, generated dependencies, or unlicensed assets. Contributions are licensed under Apache-2.0; a Developer Certificate of Origin sign-off may be introduced before accepting code.
