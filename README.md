# Aruo

**Aruo is an open-source engineering control plane for creating, understanding, maintaining, and evolving trustworthy software repositories.**

Most generators stop after copying files. Aruo treats creation as the beginning: it models repository intent, checks evidence against versioned standards, previews safe changes, reconciles upgrades, and coordinates documentation, tests, benchmarks, and releases through native ecosystem tools.

> Project status: **pre-0.1**. `aruo create` provides two qualified catalog entries (Go and JavaScript libraries), and `aruo doctor` performs versioned local repository-health assessment. Other lifecycle commands remain intentionally unavailable.

## Why Aruo

Software teams repeatedly assemble CI, security policy, documentation, release automation, tests, and community files. The pieces drift because each tool owns only one lifecycle stage. Aruo connects them through a local-first repository model and an explainable plan:

```text
intent → inspect → plan → apply → verify → record → evolve
```

## Planned capabilities

- create new projects and adopt existing repositories;
- inspect repository structure and explain detected conventions;
- check versioned policy packs and emit human, JSON, or SARIF findings;
- preview transactional fixes and provenance-aware upgrades;
- compose language-native blueprints, templates, and capabilities;
- run permissioned out-of-process plugins;
- coordinate documentation, benchmarks, releases, and migrations;
- preserve evidence without requiring a hosted service or AI.

## Quick start

With the pinned Go toolchain installed:

```sh
go run ./cmd/aruo
go run ./cmd/aruo help
go run ./cmd/aruo version
go run ./cmd/aruo create --help
go run ./cmd/aruo doctor .
```

### Install this development checkout

Aruo has not published a production release yet. To build the current checkout
and install it for your user account:

```sh
go build -o "$HOME/.local/bin/aruo" ./cmd/aruo
aruo version
aruo help
```

Ensure `$HOME/.local/bin` is on `PATH`. This installs an unversioned development
binary (`aruo version` reports `dev`); rebuild it after pulling changes. Do not
install packages claiming to be official Aruo releases before this repository
publishes verified release instructions.

Review the [create architecture](docs/architecture/create-command.md), [system architecture](docs/architecture/README.md), and [CLI specification](docs/cli/README.md) before extending the command or catalog.

## Roadmap

The first vertical slice supports read-only inspection and policy checking for Go and Python. Safe create/adopt operations follow only after plan determinism and recovery behavior are proven. See [ROADMAP.md](ROADMAP.md).

## Contributing

Early contributions should improve requirements, research, ADRs, RFCs, fixtures, architecture, and acceptance criteria—not add product commands without an accepted implementation RFC. Start with [CONTRIBUTING.md](CONTRIBUTING.md), the [RFC process](rfcs/README.md), and the [Code of Conduct](CODE_OF_CONDUCT.md).

## License and security

Aruo is licensed under [Apache-2.0](LICENSE). Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).
