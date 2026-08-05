# Development toolchain

## Requirements

- Go 1.26.5, selected through `.go-version` or `.tool-versions`.
- GNU Make 4 or a compatible Make implementation for the task facade.
- golangci-lint 2.11.4 for formatting and linting.
- GoReleaser 2.17.1 only for release validation.
- Syft 1.44.0 only for release SBOM generation.
- Git 2.40 or newer recommended.

Use a version manager that understands `.tool-versions`, install verified upstream binaries, or use the provided development container. Do not install tools with an unversioned `@latest` command in automation.

## Bootstrap

```sh
go version
make bootstrap
make check
```

`make bootstrap` downloads module dependencies and verifies external tools. It intentionally does not modify shell profiles, install system packages, or run remote install scripts.

## Daily loop

```sh
make fmt
make test
make check
```

`make help` is the task reference. CI runs the same underlying commands directly so Windows jobs do not depend on Make.

## Development container

`.devcontainer/devcontainer.json` supplies the Go base environment and editor defaults. External tools still use the pinned versions in `.tool-versions`; automatic extension/tool updates are disabled to reduce local/CI drift.

## Troubleshooting

- Tool version mismatch: compare `go version`, `golangci-lint version`, and `.tool-versions`.
- Go attempts a toolchain download: install `.go-version` or allow Go’s authenticated toolchain download.
- Lint differs from CI: remove the golangci-lint cache and verify version 2.11.4.
- `release-check` fails before application work: the release scaffold expects `cmd/aruo`; it must never publish from a snapshot command.
