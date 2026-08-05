# Go terminal UX technology stack

**Status:** Proposed technology decision

**Research date:** 5 August 2026

**Decision owner:** Aruo maintainers

**Scope:** Terminal interaction infrastructure; no product or domain behavior

## Executive decision

Aruo should remain a traditional, composable Cobra CLI and use rich terminal UI selectively. The recommended stack is:

| Layer | Selection | Use |
| --- | --- | --- |
| Command grammar | Cobra 1.x | Commands, flags, help, completion generation |
| Lifecycle | Go standard library | `context`, `os/signal`, `os/exec`, bounded cleanup |
| Terminal primitives | `golang.org/x/term` | TTY checks, size, password/raw-mode primitives |
| Prompt adapter | Huh v2 | Inputs, secrets, confirms, selects, filtering, validation, forms |
| Stateful UI runtime | Bubble Tea v2 | Only complex stateful or full-screen interactions |
| Components | Bubbles v2, transitively through Huh | Text editing, viewports, key maps; direct use only where justified |
| Styling and layout | Lip Gloss v2 | Pure styles, width-aware layout, tables, trees, ANSI degradation |
| Progress | Aruo task/event model rendered by an Aruo Bubble Tea adapter | Spinners, bars, concurrent/nested tasks, durable line fallback |
| Markdown | Glamour v2, lazy path | Finite human documentation and release-note previews |
| Highlighting | Chroma v2 through Glamour | Small, explicitly requested code samples only |
| Logging | `log/slog` | Injected structured diagnostics; JSON and text handlers |
| Configuration | Aruo-owned typed resolver; Koanf v2 only as an internal merge adapter if a spike proves useful | Layer acquisition/merge, never domain access |
| Shell integration | Cobra completion APIs | Bash, zsh, fish, PowerShell; static and bounded dynamic completion |
| Testing | Aruo fakes, golden tests, subprocess/PTY fixtures | Contract, rendering, lifecycle, accessibility, and platform tests |

The most important selection is architectural: **no package outside `internal/tux/adapters` may import a third-party terminal UI package**. Command and application code depend on Aruo-owned semantic interfaces. Huh, Bubble Tea, Lip Gloss, and Glamour are replaceable adapters, not the Aruo programming model.

This is a coherent stack because Huh v2 is implemented on Bubble Tea v2 and Bubbles v2, its themes use Lip Gloss v2, and Glamour v2 also uses Lip Gloss v2. Aruo therefore has one rendering family, one keyboard/event runtime, one width/color model, and one adapter boundary—not four unrelated UI systems. The cost is concentration risk in the Charm ecosystem and a non-trivial transitive graph. The isolation layer, capability tests, pinned releases, and a line-mode reference adapter are mandatory mitigations.

## Decision criteria

Popularity is not a criterion by itself. Candidates were judged on:

1. **Maintainer continuity:** recent releases, migration guides, security posture, and more than one meaningful contributor.
2. **Terminal correctness:** Windows/ConPTY, redirected handles, resize, paste, Unicode graphemes/display width, restoration, and cancellation.
3. **Accessibility:** a usable line-oriented screen-reader path, keyboard-only operation, no-color/no-motion support, and semantic output independent of presentation.
4. **Composability:** injectable streams and context, absence of unavoidable global state, embeddable components, and a path from inline prompts to stateful views.
5. **Operational cost:** startup, memory, binary size, dependency graph, idle CPU, render frequency, and failure modes.
6. **Replaceability:** whether Aruo can hide the dependency behind stable, product-shaped contracts.

The evaluation is based on current official documentation, repositories, release records, and module manifests. Performance claims are treated as hypotheses until measured in Aruo's own release build on its reference machines.

## Ecosystem state in 2026

The Go terminal ecosystem changed materially in 2026. Bubble Tea v2 reached stable release and introduced a new renderer, declarative terminal state, progressive keyboard enhancements, built-in color downsampling, and clearer ownership of terminal I/O. Huh v2 and Glamour v2 moved to the same `charm.land` v2 family. This removes several v1-era reasons to combine unrelated prompt, styling, and rendering packages.

The current Huh source exposes normal, password, and no-echo input; suggestions; inline validation; dynamic fields; filtered selection; custom key maps; and a first-class accessible mode that replaces cursor-addressed forms with ordinary prompts. That accessible path is decisive. Most older Go prompt packages offer useful controls but do not establish equivalent screen-reader behavior as a supported rendering mode.

The current module manifests also expose the cost. Huh is not small: it brings Bubble Tea, Bubbles, Lip Gloss, terminal/conpty support, width/grapheme packages, color handling, synchronization, and PTY-related packages. Glamour adds Goldmark, Chroma, Bluemonday, and `x/text`. These dependencies are acceptable on interactive code paths only after Aruo measures startup and binary impact; they should not be initialized for `aruo version`, cached help, or machine-only commands.

## Technology comparison matrix

Ratings are relative to Aruo's requirements: **strong**, **adequate**, **weak**, or **not applicable**. “Activity” reflects evidence available on the research date, not a guarantee of future maintenance.

### Prompt libraries

| Candidate | Capability | Accessibility | Composition | Activity in 2026 | Aruo assessment |
| --- | --- | --- | --- | --- | --- |
| Huh v2 | Strong: text, secret, confirm, select, multi-select, filtering, suggestions, validation, dynamic forms | Strong: documented line-oriented accessible mode | Strong with Bubble Tea; customizable themes/key maps | Active v2 family | **Primary adapter** |
| AlecAivazis/survey | Broad classic prompt set | Weak evidence for a first-class screen-reader mode | Traditional blocking API | Mature but no longer the leading active design | Reject for new foundation |
| manifoldco/promptui | Input/select with templates and search | Weak | Simple but narrower; templating couples presentation | Maintenance risk relative to current alternatives | Reject |
| c-bata/go-prompt | Excellent REPL completion use case | Weak; not a form system | Strong for REPLs, poor fit for wizards | Successor projects explicitly cite it as unmaintained | Reject |
| nao1215/prompt | Modern go-prompt replacement | Not yet demonstrated at Huh's level | Strong for interactive prompt/REPL scenarios | Active | Watch; not primary |
| `golang.org/x/term` alone | Line/password primitives only | Line mode is assistive-technology friendly | Excellent low-level seam | Go project maintained | **Fallback primitive**, not full prompt suite |
| Hand-built prompts | Exactly scoped | Potentially strongest if carefully designed | Total control | Aruo bears all maintenance | Use only for reference line adapter |

Huh v2 wins because it covers the required form interactions and exposes an accessibility mode within the same conceptual API. It does not win because its default visuals are attractive. Aruo must still supply its own semantics, cancellation policy, validation messages, defaults, and stream ownership.

Fuzzy search deserves a distinction. Huh provides filtered selection, which is sufficient for template/language/profile lists. A fuzzy-ranking dependency should not be added until real datasets show substring filtering is inadequate. If added later, ranking belongs behind Aruo's `Matcher` interface so prompt rendering does not own search policy.

### Rendering, styling, tables, and trees

| Candidate | Strengths | Costs and risks | Decision |
| --- | --- | --- | --- |
| Lip Gloss v2 | Pure declarative styles; ANSI-aware wrapping/width; 16/256/true-color profiles; downsampling; layout, table, list, and tree packages; same family as Huh/Bubble Tea | Large enough to require benchmarks; style-heavy APIs invite leaking vendor types | **Selected behind renderer adapter** |
| `termenv` directly | Focused color/profile output and hyperlinks | Duplicates capability/color work already owned by the selected stack | Do not add directly |
| go-pretty v6 | Capable tables, lists, progress, text; broad formatting | A second rendering/width/style system and older dependency baseline | Reject for cohesion |
| pterm | Broad “batteries included” presentation | Large opinionated surface; risks making library vocabulary the application architecture | Reject |
| ANSI literals plus `fmt` | Tiny and fast | Easy to corrupt redirection, width, reset, accessibility, and Windows behavior | Prohibit outside terminal adapters |

Aruo should expose semantic styles such as `Emphasis`, `Success`, `Warning`, `Danger`, `Muted`, and `Code`; it must not expose colors or Lip Gloss `Style`. Tables are typed data with column importance and alignment. The adapter decides whether the result is a bordered table, compact columns, key/value list, or JSON. Narrow terminals should drop low-priority columns before wrapping identifiers into unreadable output.

### Progress libraries

| Candidate | Concurrent/nested work | Line fallback | Cancellation model | Assessment |
| --- | --- | --- | --- | --- |
| Aruo event model + Bubble Tea adapter | Designed to Aruo's task hierarchy | Yes, a separate durable renderer | Context and typed terminal states | **Selected** |
| Huh spinner | Single indeterminate action | Limited | Context-aware | Use only inside the Huh adapter for a local prompt transition |
| Bubbles spinner/progress | Components for custom models | Aruo supplies it | Integrates with Bubble Tea messages | Building block, not public API |
| mpb v8 | Strong multiple bars, dynamic totals, elapsed/ETA, EWMA | Can target writers but separate fallback still needed | Rendering cancellation supported | Strong standalone alternative; reject initially for two render loops |
| schollz/progressbar | Simple progress bars | Adequate | Basic | Too narrow for task trees |
| pterm/go-pretty progress | Integrated with their presentation suite | Varies | Varies | Rejected with their parallel rendering stacks |

Progress is domain telemetry, not a spinner call. Operations emit task events: started, advanced, message, retrying, completed, failed, cancelled. One renderer owns the terminal region. Interactive mode can display concurrent/nested tasks; non-interactive mode emits sparse durable lines; machine mode emits versioned JSONL. This prevents progress libraries from infecting long-running application services and makes remote/plugin progress possible later.

### Markdown and syntax highlighting

| Candidate | Markdown | Highlighting | Cost | Decision |
| --- | --- | --- | --- | --- |
| Glamour v2 | GFM-oriented terminal rendering, wrapping, styles | Chroma integration | Heavy optional graph including Goldmark, Chroma, Bluemonday | **Selected only for explicit finite previews** |
| Goldmark + custom renderer | Excellent parser | Must integrate separately | More Aruo rendering code | Reject initially |
| Plain source output | Fully accessible and copyable | None | Minimal | Default fallback and machine behavior |
| External pager/renderer | User preference | Tool-dependent | Process/security/platform complexity | Allow pager for finite output; do not require external Markdown tool |

Terminal Markdown should be supported for `aruo docs show`, release-note previews, and finite help/tutorial content—not for ordinary status output, errors, CI logs, or machine formats. Render only trusted/local Markdown unless sanitization rules are explicit. Cap input size, width, and render time. Links must retain visible URLs when hyperlinks are unavailable.

Syntax highlighting improves comprehension for short code/config examples where language is known. It harms startup and output stability when applied indiscriminately. Use Chroma only through the Markdown adapter, only for bounded snippets, and disable it in plain, accessible, no-color, and machine modes. Do not add Chroma as a general command dependency.

### Configuration libraries

| Candidate | Layering | Typed result/provenance | Dependency posture | Decision |
| --- | --- | --- | --- | --- |
| Aruo resolver + focused YAML decoder | Exactly matches documented precedence | Can preserve source/provenance and reject unknown keys | Smallest, most controlled | **Required public architecture** |
| Koanf v2 internal adapter | Modular providers/parsers; explicit merge | Still map-oriented unless wrapped | Active v2 releases; providers are separately selectable | Candidate after spike |
| Viper | Defaults, env, files, flags, remote stores, watching | Global/registry patterns and `Get*` zero values are easy to misuse; no deep merge of complex values | Broad graph and much more surface than Aruo needs | Reject |
| env-only libraries | Simple environment mapping | Cannot represent project/workspace/org provenance | Small | Insufficient |

Configuration is included here because terminal behavior depends on color, motion, accessibility, input, pager, and output preferences. Business code receives a validated `TerminalPolicy`, never Koanf/Viper. Aruo should implement its documented precedence and provenance itself. A short Koanf v2 spike may be used inside `internal/config/adapters/koanf`, but adoption requires proof that it preserves Aruo's error locations, unknown-key policy, merge semantics, and `config explain` provenance without leaking maps. Viper is intentionally rejected despite Cobra affinity: Aruo does not need remote key/value stores, live watching, aliases, case-insensitive keys, or a broad mutable registry.

## Recommended interaction architecture

### Traditional CLI first, selective TUI

Aruo should not become a full-screen application with subcommands attached. Its primary contract remains:

```text
arguments/config -> deterministic plan -> optional approval -> effect -> durable result
```

Use line-oriented output for help, version, doctor findings, command results, errors, CI, pipes, screen readers, and most long-running tasks. Use inline Huh prompts only to resolve missing human intent. Use Bubble Tea directly only when all of these are true:

- users must navigate or filter a changing collection;
- the interaction has meaningful state that cannot be expressed clearly as sequential prompts;
- a complete flag/config/non-interactive equivalent exists;
- a line-oriented accessible equivalent exists;
- resize, cancellation, restoration, and PTY tests exist;
- the added binary/startup cost is already paid or justified.

Likely future candidates are a large audit-results explorer, plugin browser, or live benchmark dashboard. `aruo create` does not require a full-screen TUI.

### Capability and policy before rendering

The composition root computes a terminal session once from injected streams, environment, flags, configuration, and platform probes. It retains independent facts:

- stdin/stdout/stderr terminal status;
- dimensions of the presentation stream;
- color profile, Unicode/display-width confidence, hyperlinks, cursor addressing, alternate screen;
- interaction eligibility, CI, SSH, multiplexer, dumb terminal;
- accessible, reduced-motion, no-color, and no-input policy;
- output format, verbosity, pager, and redaction policy.

Do not encode this as one `Interactive bool`. Do not infer accessibility or safety from TTY status. Environment detection for `CI`, GitHub Actions, Docker, SSH, tmux, screen, Windows Terminal, VS Code, or JetBrains is advisory. User flags and Aruo configuration win. Unknown environments degrade to append-only, 80-column, no-probe behavior.

Use `golang.org/x/term` as the stable low-level default for `IsTerminal`, size, raw mode, restoration, and no-echo password primitives. Bubble Tea may own raw mode while its adapter is running; two libraries must never compete for terminal ownership. Aruo owns signal policy above either implementation.

## The Aruo terminal abstraction

### Boundary rule

Third-party imports are permitted only below these adapter packages:

```text
internal/tux/
  model/          semantic values: capabilities, style roles, tables, tasks, diagnostics
  ports/          small interfaces consumed by command/application packages
  policy/         mode selection and user overrides
  session/        stream ownership, lifecycle, cancellation, restoration
  adapters/
    plain/        stdlib reference renderer and accessible prompts
    charm/        Huh, Bubble Tea, Bubbles, Lip Gloss adapters
    markdown/     Glamour adapter
    term/         x/term and platform capability probes
    slog/         diagnostic log handlers/redaction
  testing/        fakes, event recorders, terminal fixtures
```

Neither domain packages nor `internal/create`, `internal/doctor`, plugins, template code, or configuration models may import these third-party packages. Cobra commands may depend on Aruo ports but not Charm types. Plugins exchange versioned semantic events; they never receive renderer handles.

### Interface design principles

- Interfaces describe Aruo outcomes, not library mechanisms.
- Define interfaces at the consuming boundary and keep them small.
- Pass `context.Context` on every blocking interaction.
- Return typed cancellation/EOF/unsupported errors; never vendor errors.
- Data objects contain semantic roles, not ANSI, colors, or rendering types.
- Keep requested results separate from diagnostics and progress streams.
- Prefer one `Session` assembled at the composition root over package globals.
- Provide a standard-library plain adapter as the executable specification and disaster fallback.

The following examples are architectural sketches, not implementation committed by this research task.

```go
// Aruo-owned values. No third-party types cross this boundary.
type Capabilities struct {
    InputTTY, OutputTTY, ErrorTTY bool
    Width, Height                 int
    Color                         ColorLevel
    Unicode, Hyperlinks           bool
    CursorAddressing, AltScreen   bool
}

type Policy struct {
    Mode       Mode // interactive-human, non-interactive-human, machine
    Accessible bool
    Motion     MotionPolicy
    Color      FeaturePolicy
    Unicode    FeaturePolicy
    Input      InputPolicy
    Format     OutputFormat
}

type Session interface {
    Capabilities() Capabilities
    Policy() Policy
    Presenter() Presenter
    Prompter() Prompter
    Progress() ProgressSink
}
```

Avoid a universal `UI` interface with dozens of methods. Commands should request only what they consume:

```go
type Prompter interface {
    Input(context.Context, InputRequest) (string, error)
    Secret(context.Context, SecretRequest) (Secret, error)
    Confirm(context.Context, ConfirmRequest) (bool, error)
    Select(context.Context, SelectRequest) (OptionID, error)
    MultiSelect(context.Context, MultiSelectRequest) ([]OptionID, error)
}

type InputRequest struct {
    ID          string
    Label       string
    Description string
    Placeholder string
    Default     *string
    Optional    bool
    Suggestions []Suggestion
    Validate    func(string) error
}

type SelectRequest struct {
    ID          string
    Label       string
    Description string
    Options     []Option
    Default     *OptionID
    Search      SearchPolicy
    Validate    func(OptionID) error
}
```

`OptionID` is a stable opaque identifier; displayed labels are localizable and must never become business values. `Secret` should avoid accidental formatting/logging and be zeroed best-effort after conversion where practical. The adapter maps these requests to Huh fields or accessible numbered-line prompts.

Output is semantic:

```go
type Presenter interface {
    Message(context.Context, Message) error
    Table(context.Context, Table) error
    Tree(context.Context, Tree) error
    Document(context.Context, Document) error
    Result(context.Context, Result) error
    Diagnostic(context.Context, Diagnostic) error
}

type Message struct {
    Kind MessageKind // info, success, warning, note
    Text string
}

type Column struct {
    ID         string
    Heading    string
    Alignment Alignment
    Priority   int
    Sensitive bool
}
```

The human adapter may use color and icons; plain mode uses labels; JSON serializes stable fields. Application code does not call `style.Render`, build table borders, or decide whether output is a TTY.

Progress is an event protocol:

```go
type ProgressSink interface {
    Emit(context.Context, TaskEvent) error
}

type TaskEvent struct {
    TaskID, ParentID string
    Kind             TaskEventKind
    Label            string
    Current, Total   int64
    Unit             string
    At               time.Time
    Detail           map[string]string
}
```

The event sequence must be valid independently of a renderer. The interactive adapter converts it to a task tree; the durable adapter writes sparse lines; the JSONL adapter writes versioned events; tests record it in memory. Backpressure policy must prevent rendering from blocking business work indefinitely while guaranteeing terminal states are retained.

File and clipboard interactions remain separate capabilities:

```go
type PathChooser interface {
    ChoosePath(context.Context, PathRequest) (string, error)
}

type Clipboard interface {
    WriteText(context.Context, string) error
}
```

`PathChooser` has a typed-path prompt and validation in accessible/non-interactive modes. A visual file picker is an optional adapter, not the only route. Clipboard support is opt-in (`--copy` or an explicit key), never automatic, never used for secrets by default, and returns the content normally if unavailable. OSC 52 is not a safe universal default: tmux documents inconsistent terminal support and security implications. Prefer platform APIs/tools behind the interface only after a dedicated threat and compatibility spike.

### Error isolation

Adapters map vendor conditions to Aruo errors:

```go
var (
    ErrCancelled   = errors.New("interaction cancelled")
    ErrEndOfInput  = errors.New("end of input")
    ErrUnavailable = errors.New("interaction unavailable")
)
```

Vendor errors are wrapped internally for debug diagnostics but are not matched by application code. Terminal restoration happens in an idempotent session boundary even when an adapter panics. An adapter must not call `os.Exit`, install process-global signal handlers, or close injected streams.

## Logging and error presentation

Use the Go 1.26 standard library's `log/slog`. It provides structured records, levels, text and JSON handlers, attribute replacement/redaction, and `MultiHandler`; another logging facade would add vocabulary without solving an Aruo problem.

Logging and user presentation are different channels:

- `Presenter.Diagnostic` renders one actionable user-facing error.
- injected `*slog.Logger` records causality, adapter selection, timings, retries, and debug metadata.
- default logs go to stderr only when verbose/debug policy permits them.
- `--log-format=json` produces JSONL diagnostics on stderr and does not alter stdout's requested result schema.
- secrets, tokens, authorization headers, user file content, and environment dumps are redacted before handlers.
- stack traces appear only for Aruo defects in debug output, never for expected validation, cancellation, findings, or trust refusal.

Errors should carry stable code, summary, cause, effect state, safe next action, relevant source/target, and optional documentation link. Suggestions must be deterministic and must not recommend destructive recovery without explicit review. Retrying belongs to the operation policy, not the renderer; progress events describe retries while the error retains attempt history.

## Shell integration

Retain Cobra's generated Bash, zsh, fish, and PowerShell completions. It already supports static scripts, flag completion functions, filename directives, and shell directives. Completion handlers must:

- return within 100 ms for local data;
- never prompt, animate, page, authenticate, mutate state, or print diagnostics to stdout;
- avoid network access by default and use bounded caches when remote completion is explicitly enabled;
- escape through Cobra/shell APIs rather than hand-building shell syntax;
- remain useful without configuration or network access;
- carry descriptions where the shell supports them;
- be tested in supported shell versions, including PowerShell quoting and paths with spaces.

Do not confuse shell Tab completion with Huh field suggestions. They are separate adapters over reusable catalog/search services.

## Cross-platform and terminal strategy

### Tier-one platforms

- Linux amd64/arm64 in common terminals, SSH, containers, tmux, and screen;
- macOS amd64/arm64 in Terminal.app and a modern emulator, locally and over SSH;
- Windows amd64/arm64 in Windows Terminal and PowerShell through ConPTY, with graceful legacy-console degradation.

Platform adapters own process groups, suspension/resume, resize signals, terminal handles, path display, executable-bit semantics, and console control events. Do not emulate Unix signals on Windows in application code. Ctrl+Z suspension is Unix-only. Windows Ctrl+Z console EOF is not a substitute for an on-screen cancel path.

Terminal probes can themselves corrupt input in multiplexers or remote sessions. Favor library capability detection and environment hints; run active queries only when the UI runtime owns input, with timeouts and conservative fallback. Restore raw mode, bracketed paste, mouse, cursor visibility, title, and alternate screen on normal return, error, panic, first cancellation, and suspension.

Mouse support is not part of Aruo's required interaction contract. Bubble Tea may expose it for a future full-screen view, but every operation must remain keyboard-accessible and mouse must be disabled unless the view explicitly owns it.

## Accessibility requirements for the stack

Huh's accessible mode is necessary but not sufficient. The Aruo abstraction enforces:

- `--accessible` and `ARUO_ACCESSIBLE=1` select the plain prompt adapter or Huh accessible adapter before any cursor UI starts;
- every selection has a numbered line-mode equivalent with explicit selected/default states;
- no semantic state depends only on color, icon, cursor position, animation, or rewritten text;
- `NO_COLOR`, `--color=never`, `TERM=dumb`, redirected output, and machine mode remove ANSI;
- reduced-motion mode uses durable phase changes, not spinners or rapidly rewritten bars;
- the accessible palette avoids red/green-only distinctions and background fills;
- focus order and shortcuts are stable, visible, and completely keyboard operable;
- validation repeats the prompt context and does not erase the rejected answer without explanation;
- screen-reader tests cover NVDA/Narrator, VoiceOver, and practical Linux line-mode review;
- Unicode grapheme/display-width behavior is tested, but users can force ASCII/Unicode off;
- prose is localizable; flags, option IDs, JSON keys, and machine enum values remain stable.

The plain adapter is not a degraded afterthought. It is the semantic reference implementation used by screen readers, CI-like terminal limitations, golden tests, and emergency fallback if the rich adapter fails.

## Performance and dependency budgets

No third-party README establishes Aruo performance. Before adoption, create an isolated benchmark branch and measure release binaries with and without each adapter on clean module/build caches.

Required measurements:

| Measurement | Acceptance target |
| --- | ---: |
| `aruo version` and cached help p95 | no regression beyond noise; remain <= 50 ms |
| First interactive paint p95 | <= 100 ms |
| Keypress-to-render p95 | <= 50 ms |
| Local filtering at 100/1,000/100,000 options p95 | <= 100 ms for supported size or explicit paging/limit |
| Ctrl+C acknowledgement p95 | <= 100 ms |
| Idle UI CPU | approximately 0%; event-driven |
| Render rate | capped at 10 fps for progress; coalesced resize/progress events |
| Ordinary CLI RSS | target <= 50 MiB |
| Compressed release binary | record per-platform delta; >5 MiB added by TUX requires review |

Benchmark direct imports, linked binary size (`go build`, `go tool nm`, platform artifact compression), cold/warm startup, allocations, and RSS. Measure Glamour separately because its parser, sanitizer, syntax highlighter, and language data are optional costs. Lazy initialization protects runtime startup but cannot remove linked binary size; build tags or a separate helper binary are future options only if evidence justifies their distribution complexity.

Large lists should not render every row. Filter incrementally, cap visible rows, debounce remote work, and cancel superseded searches. Progress updates should coalesce by task ID. Do not poll while idle. Never optimize away correctness, restoration, or accessible output.

## Dependency risk analysis

| Risk | Consequence | Mitigation |
| --- | --- | --- |
| Charm ecosystem concentration | One maintainer/org event affects prompts, runtime, style, and Markdown | Aruo ports, plain reference adapter, vendorable licenses, contract fixtures, annual exit drill |
| v2 ecosystem churn | APIs and import paths may continue settling | Pin exact stable tags, Renovate/Dependabot PRs, upgrade in isolated adapter commits, no vendor types outward |
| Large transitive graph | CVEs, build time, binary size, supply-chain surface | SBOM, govulncheck, license policy, module graph budget, remove unused direct imports |
| Terminal edge cases | Corrupted screen/input or platform-only bugs | One owner for terminal state, PTY matrix, panic/signal restoration, conservative fallback |
| Accessibility regressions | Rich mode becomes unusable to screen-reader users | Plain adapter parity, golden transcripts, manual assistive-technology release gate |
| Abandoned config/UI library | Forced rewrite | Aruo models/interfaces, adapter conformance suite, no persistent vendor serialization |
| OSC/clipboard security | Untrusted output changes clipboard | Explicit opt-in, content limits, no secrets, capability/threat spike |
| Markdown input complexity | Expensive or unsafe rendering | finite trusted input, size/time bounds, sanitization, plain fallback |
| Global state or competing I/O | Races, hangs, broken redirection | constructor injection; exactly one active UI runtime; prohibit library globals |

Every selected dependency needs an owner, purpose, allowed import paths, license record, current version, update cadence, binary-size delta, and removal plan in a dependency register. Review quarterly for vulnerabilities and annually for maintenance signals. A library is considered at-risk when releases/security responses stop, core compatibility issues remain unresolved, bus factor materially worsens, or it blocks supported Go/platform upgrades. Popularity is not a substitute for these signals.

### Exit strategy

At least once per major release cycle, run the adapter conformance suite against the plain implementation and one rich implementation. Every two years, perform an “exit drill”: build a spike adapter for one representative prompt, table, progress tree, cancellation flow, and accessible transcript without using the current rich package. The objective is not to ship the spike; it proves the ports remain honest.

Do not fork a dependency at the first maintenance concern. Freeze a safe version, assess exposure, upstream fixes, and compare replacement cost. Fork only for bounded security/correctness fixes with a documented sunset. If replacement is required, application and domain packages should not change.

## Integration sequence

This document authorizes no implementation. A future RFC/ADR and measured spike should proceed in reversible stages:

1. Define Aruo semantic models, ports, policy precedence, and plain/reference adapter using only the standard library and `x/term`.
2. Add contract tests for streams, cancellation, EOF, accessibility transcripts, table degradation, progress ordering, and vendor-error mapping.
3. Benchmark the current binary, then add Huh v2/Lip Gloss v2 through the Charm adapter and measure the delta.
4. Implement one representative form—not a product command—to qualify password input, validation, filtering, resize, paste, Ctrl+C, and Windows restoration.
5. Add the progress event model and durable renderer; qualify a Bubble Tea live renderer against concurrent and nested fixtures.
6. Add Glamour only when the first approved Markdown interaction exists, with separate cost and security evidence.
7. Add shell completion from Cobra independently of the interactive adapter.
8. Record the final selection in an ADR with exact versions and benchmark artifacts; keep this report as research evidence.

Each stage is a separate focused commit. No stage may combine library adoption with changes to create/doctor business behavior.

## Rejected architecture choices

### Full Bubble Tea application shell

Rejected. It would make redirection, scripting, CI, screen readers, simple invocations, and plugin output harder while providing little value to most Aruo commands. Bubble Tea remains available behind a selective adapter.

### Huh types as command APIs

Rejected. It would be quick initially but force every command, test, plugin bridge, and later GUI/agent surface to understand a terminal vendor. Commands submit Aruo requests and receive domain identifiers.

### One library for everything

Rejected. Broad suites such as pterm or go-pretty reduce initial choices but couple prompts, logs, tables, progress, and styling to one presentation vocabulary. Cohesion should come from Aruo's semantic model and one render family, not a giant public facade.

### Independent best-of-breed libraries per widget

Rejected. Huh + mpb + go-pretty + termenv + a separate picker would create multiple render loops, width/color interpretations, key models, and restoration responsibilities. Add a separate specialist only when evidence exceeds the integration cost.

### Viper because Cobra commonly pairs with it

Rejected. Cobra compatibility is not an architectural reason. Aruo requires deterministic typed layers, source provenance, unknown-key rejection, explainability, and precise merge rules; Viper's broader registry, remote, alias, watcher, and case-insensitive behavior are unnecessary.

### Automatic clipboard integration

Rejected. Local, SSH, tmux, container, and remote-shell semantics differ; OSC 52 has both inconsistent support and security implications. Explicit copy may be added later behind a capability and consent boundary.

## Final recommendation

Adopt the Charm v2 family as a **replaceable rich-terminal adapter**, not as Aruo's architecture. Keep Cobra for command grammar, Go contexts/signals for lifecycle, `x/term` for low-level primitives, `slog` for diagnostics, and a standard-library plain adapter as the source of semantic truth. Use Huh for forms, Bubble Tea only for interactions that earn a stateful UI, Lip Gloss for rendering, and Glamour/Chroma only on lazy bounded Markdown paths.

This gives Aruo a polished 2026 terminal experience without betting its ten-year product model on the survival of any one UI library. The enduring asset is the Aruo interaction protocol: typed intent in, semantic events out, capability-aware adapters at the edge.

## Primary evidence

- [Huh v2 forms, dynamic fields, accessibility, and Bubble Tea integration](https://github.com/charmbracelet/huh)
- [Huh v2 input implementation: suggestions, validation, password/no-echo, accessible input](https://github.com/charmbracelet/huh/blob/main/field_input.go)
- [Bubble Tea v2 releases and renderer/keyboard/I/O changes](https://github.com/charmbracelet/bubbletea/releases)
- [Bubble Tea current module manifest](https://github.com/charmbracelet/bubbletea/blob/main/go.mod)
- [Lip Gloss v2 styling, layout, tables, trees, profiles, and downsampling](https://github.com/charmbracelet/lipgloss)
- [Huh current module manifest and transitive stack](https://github.com/charmbracelet/huh/blob/main/go.mod)
- [Glamour v2 terminal Markdown and capability guidance](https://github.com/charmbracelet/glamour)
- [Glamour current module manifest](https://github.com/charmbracelet/glamour/blob/main/go.mod)
- [`golang.org/x/term` terminal, size, raw-mode, paste, and password primitives](https://pkg.go.dev/golang.org/x/term)
- [`log/slog` structured logging, handlers, redaction hooks, and Go 1.26 `MultiHandler`](https://pkg.go.dev/log/slog)
- [mpb v8 concurrent progress features and release activity](https://github.com/vbauerster/mpb)
- [Lip Gloss table/tree alternatives](https://github.com/charmbracelet/lipgloss) and [go-pretty module manifest](https://github.com/jedib0t/go-pretty/blob/main/go.mod)
- [Koanf v2 release activity](https://github.com/knadh/koanf/releases) and [Viper feature/precedence model](https://github.com/spf13/viper)
- [Cobra shell completion for Bash, zsh, fish, and PowerShell](https://cobra.dev/docs/how-to-guides/shell-completion/)
- [Windows Terminal and ConPTY platform context](https://github.com/microsoft/terminal)
- [tmux clipboard behavior, terminal gaps, and OSC 52 security concerns](https://github.com/tmux/tmux/wiki/Clipboard)
