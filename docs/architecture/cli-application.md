# CLI Application Architecture

## Purpose

This document defines the maintainability boundaries of the Aruo executable. It is an application architecture, not a catalogue of future commands.

## Runtime shape

```text
operating system
      │ args, stdio, signals
      ▼
cmd/aruo                 process composition and exit
      │ explicit dependencies
      ▼
internal/cli             execute once; render an error once
      ├── command        syntax, flags, help, argument validation
      └── iostreams      input/output capabilities
      │ application request
      ▼
future application services
      │ ports
      ▼
domain core ───────── adapters (filesystem, Git, network, process)
```

Dependencies point inward. Domain and application packages must not import Cobra, terminal libraries, configuration decoders, or OS process globals.

## Package layout

| Package | Responsibility | Must not do |
| --- | --- | --- |
| `cmd/aruo` | Read process resources, compose dependencies, choose exit code | Contain command behavior or domain rules |
| `internal/cli` | Construct one invocation and own the top-level error boundary | Become a service locator |
| `internal/cli/command` | Define command grammar and adapt inputs to application requests | Read global streams, call `os.Exit`, or contain workflows |
| `internal/cli/iostreams` | Carry stdin, stdout, stderr and later terminal capabilities | Decide business output |
| `internal/buildinfo` | Expose linker-injected immutable build identity | Discover Git state at runtime |
| future `internal/config` | Resolve typed layers with source provenance | Expose decoder-specific maps to the domain |
| future `internal/app` | Orchestrate use cases and transactions | Format terminal output |
| future domain packages | Implement deterministic rules and types | Import infrastructure or presentation packages |

New command families may receive their own child package when they have multiple commands. A file-per-command flat package is preferred until cohesion or compile-time boundaries justify another package.

## Dependency injection

Aruo uses constructor injection and ordinary Go values. `main` creates process dependencies; `cli.Run` is the composition root; each command constructor receives only what it needs. Commands are constructed fresh for every invocation.

No package-level command variables, `init()` registration, global configuration singleton, or dependency-injection framework is allowed. This keeps the graph searchable, avoids order-dependent tests, and makes parallel tests safe. Interfaces are defined at the consuming boundary only when a second implementation or test seam exists.

## Command system

Cobra is a presentation dependency. It provides mature parsing, nested help, suggestions, and future completion generation, but it does not define Aruo's core API. Constructors return `*cobra.Command`; handlers use `RunE`; leaf commands reject unexpected positional arguments.

The root owns stable global behavior. Commands own their local flags. Persistent flags are rare because they create implicit coupling. Automatic completion is disabled until Aruo deliberately designs and tests that public surface.

The current pre-0.1 public tree is intentionally limited to:

```text
aruo
├── create
├── doctor
├── help
└── version
```

`aruo` prints help. `aruo help` uses Cobra's generated help command. `aruo version` emits one script-friendly line. `aruo create` delegates to the catalog-neutral creation service described in [Create Command Architecture](create-command.md). `aruo doctor` delegates to the read-only check engine in [Repository Doctor Architecture](doctor.md). No application configuration is loaded for informational paths.

## Configuration loading

Configuration will be loaded lazily after command parsing and only for commands that require it. The planned pipeline is:

```text
defaults → organization → workspace → project → local → environment → flags
            resolve + retain provenance → validate typed model → application request
```

Parsing, merging, migration, and validation are separate stages. Unknown keys fail with a source location and suggestion. Secrets remain references. A command must be able to explain the winning source of a value. Help and version must remain fast and functional when project configuration is invalid.

## Logging and terminal output

Operational logs and user output are different channels:

- User results go to stdout.
- Diagnostics, warnings, progress, and errors go to stderr.
- Structured logs use the standard library's `log/slog` and are injected; libraries never configure a global logger.
- Machine output, when introduced, is rendered from typed results rather than scraped human text.
- Color, animation, and width decisions belong to `iostreams` capabilities and must respect redirection, `NO_COLOR`, accessibility, and non-interactive execution.

The first scaffold intentionally has no color or progress dependency. Complexity will be introduced only with a command that needs it.

## Error handling

Errors travel upward with context and are rendered exactly once in `cli.Run`. Deep packages neither print nor terminate the process. Future typed errors may carry a stable category, safe user message, remediation hint, and exit class while preserving a wrapped cause for logs.

The intended exit contract is `0` success, `1` operation failure, `2` invalid invocation, and distinct documented codes only when automation has a demonstrated need. Secrets and untrusted payloads must never enter user-facing or structured error fields.

## Testing strategy

- Command tests construct a fresh tree with buffers and deterministic build metadata.
- Parser and help tests verify public command grammar and golden output only where layout stability matters.
- Application services use fakes at narrow consumer-owned ports; domain tests remain framework-free.
- Filesystem behavior uses `t.TempDir`; process and network boundaries use explicit adapters.
- Integration tests run the compiled binary for exit codes, signal behavior, and stdout/stderr separation.
- Race tests are mandatory; tests may run in parallel because command and stream state is not global.
- Startup time and binary size become release benchmarks before additional UI/configuration dependencies are accepted.

## Extension rules

To add a future command: define its user contract, create an application service independent of Cobra, construct a thin command adapter, inject narrow dependencies, add command and application tests, then update generated reference documentation. A command handler that accumulates workflow decisions is an architecture defect.
