# CLI Application Architecture

## Purpose

This document defines the maintainability boundaries of the Aruo executable. It is an application architecture, not a catalogue of future commands.

## Runtime shape

```text
operating system
      │ args, stdio, signals
      ▼
cmd/aruo                 process composition; owns lifecycle.Manager and exit
      │ explicit dependencies
      ▼
internal/cli             execute once; render an error once
      ├── command        syntax, flags, help, argument validation, sessionFactory
      └── iostreams      input/output capabilities
      │ application request                    │ tux.Session (per invocation)
      ▼                                         ▼
application services (create, doctor)   internal/tux
      │ ports                                   ├── model/ports   semantic values, Prompter/Presenter/ProgressSink
      ▼                                         ├── term/policy   capability detection, mode resolution
domain core ─── adapters (filesystem, Git, ...) ├── plain         dependency-free reference adapter
                                                 ├── charm         Huh/Bubble Tea/Lip Gloss, isolated
                                                 └── lifecycle     signals, cancellation, restoration
```

Dependencies point inward. Domain and application packages must not import Cobra, terminal libraries, configuration decoders, or OS process globals. `internal/tux` is the one exception to "no terminal libraries": it is the sanctioned isolation boundary, and only its `charm` and `term` subpackages may import Charm v2 or `golang.org/x/term`. Everything else — including `internal/cli/command` — depends only on `internal/tux`'s own ports (`tux.Prompter`, `tux.Presenter`, `tux.ProgressSink`, `tux.Session`). See [ADR-0010](../../decisions/0010-charm-v2-terminal-ux-stack.md) and [the terminal UX specification](../cli/terminal-ux.md).

## Package layout

| Package | Responsibility | Must not do |
| --- | --- | --- |
| `cmd/aruo` | Read process resources, own `lifecycle.Manager`, compose dependencies, choose exit code | Contain command behavior or domain rules |
| `internal/cli` | Construct one invocation, own the top-level error boundary, classify cancellation vs. operational errors | Become a service locator |
| `internal/cli/command` | Define command grammar, resolve one `tux.Session` per command via `sessionFactory`, adapt inputs to application requests | Read global streams, call `os.Exit`, import Charm/Huh/Bubble Tea directly, or contain workflows |
| `internal/cli/iostreams` | Carry stdin, stdout, stderr | Decide business output |
| `internal/tux` | Semantic terminal contracts (`model.go`, `ports.go`, `errors.go`) consumed by commands | Import third-party terminal libraries itself |
| `internal/tux/term` | Detect capabilities via `golang.org/x/term`; the one place that may | Infer capabilities from `TERM` substrings alone |
| `internal/tux/policy` | Resolve deterministic mode/feature policy from capabilities, environment, and flags | Read process globals directly (takes an environment map) |
| `internal/tux/plain` | Dependency-free reference `Prompter`/`Presenter`/`ProgressSink` adapter; accessibility source of truth | Depend on any other `internal/tux` adapter |
| `internal/tux/charm` | Isolated Huh/Bubble Tea/Lip Gloss adapters behind the same ports | Leak Charm types past its own package boundary |
| `internal/tux/session` | Assemble capabilities, policy, and adapter selection once per invocation | Own process signals (that is `internal/tux/lifecycle`) |
| `internal/tux/lifecycle` | Own signal subscription, two-tier Ctrl+C, SIGTERM, idempotent restoration | Render user-facing output |
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

The current public tree is intentionally limited to:

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
- Diagnostics, warnings, prompts, and progress go to stderr, rendered through `tux.Presenter`/`tux.Prompter`/`tux.ProgressSink`, never raw `fmt` calls in `internal/cli/command`.
- Structured logs use the standard library's `log/slog` and are injected; libraries never configure a global logger.
- `doctor --format json` is rendered directly from the typed `doctor.Report` (never through a `tux.Session`, and never scraped from human text); only the `human` path builds a session.
- Color, animation, and width decisions belong to `tux.Capabilities`/`tux.Policy`, resolved once per invocation by `internal/tux/session`, and must respect redirection, `NO_COLOR`, `TERM=dumb`, accessibility, CI, and non-interactive execution. See [the terminal UX specification's implementation status](../cli/terminal-ux.md#implementation-status-2026-08-06) for exactly what is wired today.

Help and version never construct a `tux.Session`; they remain dependency-free and fast (see the benchmark evidence in [benchmarks/results/](../../benchmarks/README.md)).

## Error handling

Errors travel upward with context and are rendered exactly once in `cli.Run`. Deep packages neither print nor terminate the process. `cli.Run` classifies an error as ordinary cancellation — printing `Cancelled.` rather than a red `Error:` line — when it matches `context.Canceled` (the `lifecycle.Manager`-driven signal path) or `tux.ErrCancelled` (a rich adapter's own Ctrl+C key handling, which real terminals can reach even when the OS-level signal never fires, because raw mode clears `ISIG`). Future typed errors may carry a stable category, safe user message, remediation hint, and exit class while preserving a wrapped cause for logs.

The intended exit contract is `0` success, `1` operation failure, `2` invalid invocation, `130` interrupted (Ctrl+C, cooperative or forced), `143` terminated (SIGTERM), and further documented codes only when automation has a demonstrated need. `cmd/aruo` prefers `lifecycle.Manager.ExitCode()` over `cli.Run`'s own return value whenever a signal actually fired. Secrets and untrusted payloads must never enter user-facing or structured error fields.

## Testing strategy

- Command tests construct a fresh tree with buffers and deterministic build metadata.
- Parser and help tests verify public command grammar and golden output only where layout stability matters.
- Application services use fakes at narrow consumer-owned ports; domain tests remain framework-free.
- Filesystem behavior uses `t.TempDir`; process and network boundaries use explicit adapters.
- Integration tests run the compiled binary for exit codes, signal behavior, and stdout/stderr separation (`cmd/aruo/main_test.go` sends real SIGINT/SIGTERM to the built binary via a test-only `ARUO_TEST_SIGNAL_MODE` hook, since neither `os.Exit` nor real signal delivery is observable in-process).
- Terminal-capability tests fake a real TTY by wrapping a stream with a `Fd() uintptr` method and pairing it with an injected `term.Probe`, rather than requiring an actual PTY (`internal/tux/term`, `internal/tux/session`, `internal/cli`). This exercises the real capability-detection and adapter-selection code deterministically; it does not replace genuine PTY keystroke testing, which remains a known gap (see the terminal UX specification's implementation status).
- Race tests are mandatory; tests may run in parallel because command and stream state is not global.
- Startup time, Ctrl+C acknowledgement latency, and binary size are release benchmarks, recorded under `benchmarks/results/`, before additional UI/configuration dependencies are accepted.

## Extension rules

To add a future command: define its user contract, create an application service independent of Cobra, construct a thin command adapter, inject narrow dependencies, add command and application tests, then update generated reference documentation. A command handler that accumulates workflow decisions is an architecture defect.
