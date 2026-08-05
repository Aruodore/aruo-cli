# Pull request summary: maintainable CLI architecture and minimal scaffold

## Summary

This change converts Aruo's empty executable bootstrap into a deliberately minimal, testable CLI shell. It implements only `aruo`, `aruo help`, and `aruo version`. It adds no repository workflow, configuration parser, plugin runtime, template engine, completion command, or other product feature.

The design follows a study of Cobra, Hugo, Helm, Docker CLI, GitHub CLI, Terraform, uv, Bun, and Deno. The research record separates observed ecosystem patterns from Aruo's inferred decisions.

## Architecture

- `cmd/aruo` is the operating-system boundary. It creates a signal-aware context, system streams, logger, and build metadata, invokes the application once, and exits with the returned code.
- `internal/cli` is the composition and error boundary. It constructs a fresh command tree per invocation and renders an error exactly once.
- `internal/cli/command` contains Cobra adapters. Commands are constructor-built rather than registered through globals or `init()`, so dependency edges remain explicit and tests can run in parallel.
- `internal/cli/iostreams` carries stdin, stdout, and stderr. Terminal capabilities such as TTY, color, width, and interactivity will be added at this boundary only when required.
- `internal/buildinfo` exposes immutable version, commit, and build date values supplied through release linker flags.
- Future application and domain packages remain independent of Cobra, terminal libraries, configuration decoders, and process globals.

Ordinary constructor injection was selected over a dependency-injection framework. Interfaces will be defined by consumers only after a real substitution seam exists. This avoids both global state and speculative abstraction.

## Command behavior

- `aruo` prints generated root help and exits successfully.
- `aruo help` uses Cobra's standard help system.
- `aruo version` prints one stable line: `aruo version <version>`.
- Unknown commands are reported once on stderr and return a failure code.
- Cobra's automatic completion command is disabled because it is not part of this milestone.
- No configuration is loaded for help or version, preserving fast and resilient informational paths.

## Dependency decision

Cobra v1.10.2 is pinned as the sole direct runtime dependency. It supplies mature parsing, help, validation, and future extensibility, but is confined to the CLI presentation package. The project does not adopt Cobra's generator convention of package-global commands and `init()` mutation.

The standard `log/slog` package is injected for future diagnostics. No color, progress, configuration, or logging framework dependency was added before a demonstrated need.

## Documentation and records

- Adds a complete CLI application architecture covering package rules, dependency injection, commands, configuration loading, logging, terminal output, error handling, tests, and extension rules.
- Adds the modern CLI ecosystem research report with primary-source links and extracted lessons.
- Adds ADR-0008 accepting a constructor-built Cobra presentation layer.
- Updates root architecture, CLI status, project tree, getting started, and README content to distinguish the implemented shell from proposed commands.

## Release integration

GoReleaser now injects semantic version, Git commit, and build date into `internal/buildinfo` while retaining CGO-free, cross-platform, trimmed builds. Local builds report `dev` honestly rather than fabricating a release version.

## Testing and quality evidence

The CLI suite uses injected buffers and deterministic build metadata. It covers root help, explicit help, version output, unknown-command errors, stdout/stderr separation, and exit status. Fresh command construction permits parallel execution without shared command state.

Verification performed:

- `go mod tidy`
- `go test -race ./...`
- `go vet ./...`
- `go build -trimpath ./cmd/aruo`
- manual execution of root, help, version, and unknown-command paths
- `golangci-lint v2.11.4 run ./...` — zero issues
- `goreleaser v2.17.1 check` — configuration valid

## Deliberate exclusions

- No proposed product commands.
- No configuration implementation; only its loading contract is designed.
- No machine-output schema, color, prompts, spinners, or shell completion yet.
- No public package API.
- No application/domain service placeholders without behavior.
- No plugin or template loading.

## Review focus

Reviewers should validate that the dependency direction is enforceable, the current command surface contains only the three authorized paths, build metadata is reproducible, errors are not double-rendered, and the documentation clearly labels future behavior as proposed.

## Rollback

The scaffold is isolated from future domain behavior. Reverting it removes the Cobra dependency and command shell without changing repository standards, research history, or release governance. No external state or published interface beyond pre-0.1 help/version is affected.
