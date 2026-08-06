# ADR-0010: Charm v2 terminal UX stack behind an Aruo adapter boundary

- Status: Accepted
- Date: 2026-08-06
- Owners: Aruo maintainers
- Related: [docs/cli/terminal-ux.md](../docs/cli/terminal-ux.md), [research/technology/2026-go-terminal-ux-stack.md](../research/technology/2026-go-terminal-ux-stack.md), ADR-0008

## Context

Aruo needs prompts, styled output, progress reporting, and signal handling
that behave correctly across interactive terminals, redirected pipes, CI,
SSH/tmux, screen readers, and Windows, without letting a UI library become
Aruo's programming model. The research report evaluated the 2026 Go terminal
ecosystem against maintainer continuity, terminal correctness, accessibility,
composability, operational cost, and replaceability, and recommended a
Charm v2 stack behind a strict adapter boundary. This ADR records the
decision actually implemented, with measured evidence, rather than
authorizing it in the abstract.

## Decision

Adopt `charm.land/huh/v2` v2.0.3, `charm.land/bubbletea/v2` v2.0.2,
`charm.land/bubbles/v2` v2.0.0, and `charm.land/lipgloss/v2` v2.0.5 as
replaceable rich-terminal adapters, alongside `golang.org/x/term` v0.45.0 for
low-level TTY/size/raw-mode primitives and Cobra v1.10.2 (ADR-0008) for
command grammar. No package outside `internal/tux/charm` and
`internal/tux/term` may import these third-party packages; command and
application code depend only on the Aruo-owned `tux.Prompter`,
`tux.Presenter`, `tux.ProgressSink`, and `tux.Session` ports.

`internal/tux/plain` is a hand-written, dependency-free adapter implementing
the same ports. It is not a degraded fallback: it is the semantic reference
implementation, and `internal/tux/session` selects it deterministically
whenever the rich adapter cannot degrade gracefully — non-interactive mode,
machine output formats, `--accessible`/`ARUO_ACCESSIBLE`, and
`TERM=dumb` (which clears `CursorAddressing` even on a real TTY, and Huh's
forms have no in-place degraded path for a terminal that cannot address the
cursor).

`internal/tux/lifecycle` owns process signals with the standard library
(`os/signal`, `context.WithCancelCause`) rather than Bubble Tea's own signal
handling, which is disabled (`tea.WithoutSignalHandler`) on every rich
adapter Aruo constructs. Progress rendering is Aruo's own task/event model
(`tux.TaskEvent`) rendered by a Bubble Tea adapter capped at 10 FPS, not a
general-purpose progress-bar library.

## Consequences

- Commands, tests, and any future plugin bridge depend on a small, stable,
  Aruo-shaped vocabulary; Charm types never cross the adapter boundary.
- The plain adapter gives a genuine, tested accessibility and CI/machine
  path independent of the rich stack's health.
- Binary size grows by approximately 2.93 MiB (stripped) attributable to the
  Charm v2 family, under the research report's 5 MiB review threshold; see
  [benchmarks/results/2026-08-06-terminal-ux-baseline.md](../benchmarks/results/2026-08-06-terminal-ux-baseline.md).
- `aruo version` remains fast (median ~30 µs in-process, ~8 ms full process
  wall time; both far under the 50 ms budget) because help/version never
  construct a `tux.Session`.
- Ctrl+C acknowledgement is not a raw-mode guarantee: real terminals clear
  `ISIG` under `MakeRaw`, so a Ctrl+C typed inside an active Huh form never
  reaches the OS signal handler and surfaces only as `tux.ErrCancelled`.
  `cli.Run` treats that and `context.Canceled` as the same "ordinary
  cancellation" outcome so neither is misrendered as a red operational
  error.
- Concentration risk in one dependency family (Charm) is accepted, mitigated
  by the plain adapter, pinned exact versions, and this ADR's exit strategy
  rather than avoided outright.
- Not yet delivered: a JSONL progress adapter for machine mode (progress
  currently falls back to the plain adapter's sparse durable lines in
  machine mode, not a versioned structured event stream), Glamour/Chroma
  Markdown rendering, and a `--interactive` path that opens `/dev/tty`
  explicitly when stdin carries piped data. These remain scoped future work
  under the same adapter boundary.

## Alternatives considered

- **A full Bubble Tea application shell**: rejected; it would make
  scripting, CI, and screen readers harder for commands that do not need a
  full-screen view, and Aruo's primary contract is a traditional CLI, not a
  TUI with subcommands attached.
- **Huh/Bubble Tea/Lip Gloss types as the command API**: rejected; it would
  force every command, test, and future plugin bridge to understand a
  terminal vendor's model instead of Aruo's own typed requests and results.
- **One "batteries included" library (pterm, go-pretty)**: rejected for
  cohesion; broad suites couple prompts, tables, progress, and styling to a
  single vendor vocabulary instead of Aruo's semantic model.
- **Independent best-of-breed libraries per widget** (Huh + mpb + go-pretty +
  termenv + a separate picker): rejected; multiple render loops, width/color
  interpretations, and restoration responsibilities cost more than a
  specialist would save without demonstrated need.
- **AlecAivazis/survey, promptui, go-prompt**: rejected as a foundation; none
  documents an equivalent first-class line-oriented accessible mode.
- **Viper for configuration**: rejected independently (out of scope for this
  ADR); Aruo's own typed resolver retains provenance and unknown-key
  rejection that a general registry does not provide.

## Validation

Re-run the adapter conformance suite (`internal/tux/...`, `internal/cli/...`,
`cmd/aruo/...`) against both the plain and rich adapters at least once per
major release cycle. Every two years, or sooner if a maintenance signal
warrants it, build a spike adapter for one representative prompt, table,
progress tree, cancellation flow, and accessible transcript without the
current rich package, to prove the ports stay honest. Review
`benchmarks/results/` for regressions greater than 20% on normalized
hardware before adding further Charm-family surface area (Glamour/Chroma).
