# Terminal UX specification

**Status:** Normative design specification  
**Research baseline:** 5 August 2026  
**Applies to:** Every Aruo command, prompt, progress renderer, pager, and future plugin-hosted interaction

## Purpose

This document defines how Aruo behaves as a process inside a terminal. It is the implementation contract for signals, keyboard input, prompts, progress, output, errors, accessibility, terminal compatibility, performance, and testing.

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** are normative.

Aruo is a command-line application first. It may use inline terminal UI where that improves a decision, but it must remain predictable when redirected, automated, accessed through SSH, read with assistive technology, or run in a limited terminal. Rich presentation is an enhancement over a complete line-oriented contract.

## Research conclusions

### What modern CLIs converge on

Current professional tools differ visually but converge on several interaction contracts:

- GitHub CLI separates plain/JSON output, paging, color, TTY forcing, prompt disabling, accessible colors, accessible prompts, and spinner disabling into distinct controls.
- Docker BuildKit chooses TTY progress only for a terminal and otherwise provides plain, quiet, or raw JSON progress.
- `kubectl` tells scripts to request an explicit machine format and avoid implicit context; it offers JSON, YAML, names, JSONPath, templates, and header-free output.
- Terraform separates plan from apply, supports `-input=false`, and makes JSON UI imply no input.
- Cargo independently controls quiet/verbose output, color, Unicode, hyperlinks, progress, and progress width.
- `uv`, Bun, npm, and pnpm expose explicit progress suppression, reporter, and verbosity controls. pnpm selects cursor-updating output for a TTY and append-only output otherwise, while also offering NDJSON.
- Railway switches search and scaling commands between an interactive picker/TUI and structured non-TTY behavior; missing non-interactive inputs fail rather than opening prompts.
- Wrangler avoids interactive deployment choices when they could overwrite remote settings and offers JSON/NDJSON for automation.
- Supabase, Azure, AWS, and Google Cloud CLIs provide explicit output formats, debug controls, non-interactive authentication paths, and output suppression.
- GitHub CLI and Google Cloud CLI now expose explicit screen-reader modes; modern Go prompt libraries provide a line-oriented accessible mode instead of assuming cursor-addressed TUIs work with assistive technology.

The recurring design is not “detect a TTY and turn everything fancy.” It is a capability matrix: input, cursor control, width, color, Unicode, motion, paging, hyperlinks, and machine output are independent.

### What Aruo should learn from each family

| Family | Strong terminal pattern | Pattern Aruo should avoid copying blindly |
| --- | --- | --- |
| GitHub CLI | Separate output fields, templates, paging, accessible modes, environment controls | Many overlapping color environment variables |
| Docker | `auto`, `tty`, `plain`, `quiet`, and raw progress modes | PTY attachment merging stderr into stdout |
| kubectl | Explicit stable script output; dry-run; visible context | Huge inherited flag surfaces and implicit active context |
| Terraform | Plan before effects; saved plans; machine UI disables input | Machine JSON streams whose compatibility is unclear to casual users |
| Helm | Debug, color policy, dry-run, atomic/wait options | Dry-run output that may expose secrets without an extra safeguard |
| Cargo | Independent color, Unicode, hyperlink, verbosity, and progress controls | Configuration knobs before a demonstrated need |
| uv | Fast startup, repeated verbosity, no-progress/offline modes, long help via pager | A very large global option surface |
| Bun/Deno | Fast feedback and clear permission/security choices | Runtime-specific conventions outside Aruo's scope |
| npm/pnpm | TTY-aware progress and multiple log levels | Noisy install output and vocabulary that varies by command |
| Vite | Short interactive creation path with complete non-interactive parity | Prompt-first behavior where automation inputs are incomplete |
| Fly/Railway | Browserless authentication, streaming logs, detached/CI modes | Inferring permission for remote effects from interactivity |
| Supabase/Wrangler | Debug to stderr, structured formats, workdir/cwd, defensive non-interactive behavior | Emoji-heavy status as the only semantic signal |
| Azure/AWS/gcloud | Rich query/format controls, pager controls, output suppression, debug modes | JSON as a human default and inconsistent error schemas |

## Foundational interaction model

### One command, three presentation modes

Aruo has three modes derived after argument parsing:

1. **Interactive human:** stdin and stderr are usable terminals, human output is selected, CI is not detected, and input has not been disabled.
2. **Non-interactive human:** line-oriented human/plain output without prompts or cursor movement.
3. **Machine:** a structured format such as JSON, JSONL, or SARIF; prompts, paging, color, icons, animation, and prose decoration are disabled.

Mode selection MUST be deterministic and explainable. `aruo config explain terminal.mode` SHOULD report why a mode was selected when that command exists.

Interactive eligibility requires all of:

- stdin is a terminal;
- stderr is a terminal suitable for prompts;
- output format is `human`;
- neither `--no-input` nor `ARUO_NO_INPUT` disables input;
- CI detection does not disable input, unless `--interactive` explicitly opts in.

Redirecting stdout alone does not forbid prompts because requested results and prompts use separate streams. Machine format always forbids prompts. `TERM=dumb`, inaccessible cursor control, or an explicit accessibility mode selects line-oriented prompts even if stdin/stderr are terminals.

### Stream ownership

- **stdout:** requested result only.
- **stderr:** prompts, progress, warnings, diagnostics, and the single rendered error.
- **stdin:** user answers or requested input data, never both ambiguously.

A command that consumes data from stdin MUST NOT also prompt from stdin. It must require flags/configuration for missing decisions or open `/dev/tty` only under an explicit, documented `--interactive` request. Plugins MUST return typed events; they MUST NOT write directly to Aruo's terminal streams.

### Interaction pipeline

Every consequential command follows:

```text
parse -> discover -> resolve -> validate -> plan -> approve -> apply -> verify -> summarize
```

Interactive mode gathers unresolved intent before planning. `--yes` approves a complete deterministic plan; it MUST NOT invent required identity, credentials, targets, license choices, or conflict policy. Machine and no-input modes fail early with every missing value listed together.

## Process lifecycle and cancellation

### Signal contract

The process composition root owns signals. Application and domain packages receive cancellation through `context.Context`; they do not subscribe to operating-system signals.

| Event | Required behavior | Exit status |
| --- | --- | ---: |
| First Ctrl+C / SIGINT | Cancel root context, stop animation, acknowledge cancellation, stop starting work, begin safe cleanup | `130` unless operation completed before cancellation was observed |
| Second Ctrl+C during cleanup | Restore terminal best-effort and force immediate process termination | `130` |
| SIGTERM | Cancel without prompting, perform bounded cleanup, then terminate | `143` |
| Terminal close/logoff/shutdown on Windows | Best-effort cancellation and terminal restoration within OS deadline | non-zero; `143` where representable |
| SIGQUIT on Unix | Preserve Go/system diagnostic behavior; do not turn it into ordinary cancellation | platform default |
| SIGKILL | Cannot be handled; transactional design must make interruption recoverable | platform default |
| SIGTSTP / Ctrl+Z on Unix | Restore terminal modes, allow shell suspension, re-detect and redraw after SIGCONT | shell-managed |

The first interrupt MUST produce at most one concise acknowledgement on stderr:

```text
Cancelling… Press Ctrl+C again to exit immediately.
```

Accessible/reduced-motion mode uses `Cancelling...` and no animated glyph. Repeated signals MUST NOT produce repeated lines.

### Cancellation propagation

- The root context MUST reach every network request, filesystem plan/application stage, lock wait, plugin RPC, and child process.
- New tasks MUST NOT start after cancellation.
- Concurrent workers MUST drain or stop through a shared cancellation cause.
- Child processes MUST receive an appropriate graceful signal as a process group where Aruo owns the group.
- Aruo MUST continue reading child stdout/stderr until the child exits or the force deadline expires, preventing pipe deadlock.
- Cleanup MUST use a separate bounded context; it must not reuse an already-cancelled context for essential rollback.
- Cancellation errors MUST preserve the initiating cause so SIGINT, timeout, parent cancellation, and internal failure are distinguishable.

### Cleanup and rollback

Cancellation is not failure recovery by wishful thinking. Each mutating engine MUST define a commit boundary.

- Before the commit boundary, staging artifacts and locks are removed.
- After an atomic commit, Aruo reports that the operation completed even if presentation was interrupted.
- For multi-system effects that cannot be atomic, Aruo records completed steps and prints or stores a recovery plan.
- Rollback MUST be idempotent and MUST NOT destroy pre-existing user state.
- Temporary paths MUST be recoverable or recognizable after SIGKILL/power loss.
- On cancellation, secrets and partially rendered sensitive values MUST NOT be printed.

The final cancellation message distinguishes outcomes:

```text
Cancelled. No repository changes were applied.
Cancelled after 3 of 5 operations. Run `aruo recover <id>` for the recovery plan.
```

### Deadlines

- Local cancellation acknowledgement: at most 100 ms at p95.
- Cooperative local cleanup target: 2 seconds.
- Network/subprocess graceful cleanup target: 10 seconds unless the operation documents a stronger constraint.
- SIGTERM cleanup MUST be bounded and non-interactive.
- A second Ctrl+C bypasses all remaining deadlines.

### Exit-code contract

```text
0    success / policy conformant
1    operational failure
2    invocation or configuration error
3    completed assessment with findings at/above threshold
4    unresolved plan or merge conflict
5    trust or security refusal
124  Aruo-owned timeout
130  interrupted by SIGINT / Ctrl+C
143  terminated by SIGTERM
```

Commands MUST document deviations. Expected findings are not operational errors. A failure after partial external effects retains the appropriate operational/trust code and includes recovery state.

## Keyboard interaction specification

### Layering rule

Keyboard behavior has three owners:

1. the shell/terminal driver before Aruo starts or while it is suspended;
2. Aruo's line editor inside a text prompt;
3. an Aruo picker or full-screen view.

Aruo MUST NOT claim shell features it does not own. In particular, Ctrl+R shell history is not automatically an Aruo feature, and Ctrl+Z is Unix job control rather than a generic “undo” key.

### Global keys while Aruo is interactive

| Key | Aruo behavior | Scope and caveats |
| --- | --- | --- |
| Ctrl+C | Cancel the entire invocation; second press forces exit | Everywhere; never merely clears a field |
| Ctrl+D | End input; at an empty single-line prompt, cancel with an EOF explanation | Unix terminals; Windows commonly uses Ctrl+Z then Enter for console EOF, so never require this shortcut |
| Ctrl+Z | Suspend on Unix after restoring terminal; redraw on resume | Not advertised on native Windows; MUST NOT mean undo |
| Ctrl+L | Redraw current interaction; clear only Aruo's rendered region where possible | Full-screen view may clear and redraw; preserve entered data |
| Escape | Close help/search overlay or go back one level | Must not silently cancel the whole operation on first press |
| `?` | Show contextual key help | Pickers/views only; inserts `?` in text fields |
| Ctrl+G | Cancel current search/overlay and restore prior state | Follows Readline search convention where supported |

### Text fields

Text input SHOULD provide familiar Readline/Emacs-style editing without persistent history of sensitive answers.

| Key | Expected behavior |
| --- | --- |
| Left / Right | Move one grapheme, not one byte |
| Ctrl+B / Ctrl+F | Move one grapheme backward/forward |
| Home / Ctrl+A | Move to start of line |
| End / Ctrl+E | Move to end of line |
| Ctrl+Left / Alt+B | Move one word backward when distinguishable |
| Ctrl+Right / Alt+F | Move one word forward when distinguishable |
| Backspace | Delete previous grapheme |
| Delete / Ctrl+D with non-empty buffer | Delete next grapheme |
| Ctrl+U | Delete from cursor to start of line |
| Ctrl+K | Delete from cursor to end of line |
| Ctrl+W | Delete previous word/whitespace-delimited unit |
| Ctrl+Y | Restore the most recently deleted text within this prompt session |
| Ctrl+R | No global project-answer history; MAY search choices only when the prompt explicitly says so |
| Tab | Accept a unique completion or open completion choices; otherwise insert nothing |
| Shift+Tab | Move to previous field only in multi-field forms; otherwise reverse completion |
| Up / Down | Move through wrapped lines or ephemeral values entered in the current form; never secret history |
| Enter | Validate and submit the complete field |
| Space | Insert a space |

Inputs MUST support paste, including multiline paste into a single-line field by rejecting or normalizing newlines visibly. Bracketed paste SHOULD be used when available. Unicode editing MUST operate on grapheme clusters and display width, not bytes or rune count alone.

Password/token fields MUST disable echo, selection previews, history, completion, and debug logging. Pasted secrets remain allowed. The prompt MUST say when input is hidden. Reveal toggles are discouraged; if offered, they require an explicit key and visible state.

### Single-select and searchable lists

| Key | Expected behavior |
| --- | --- |
| Up / `k` | Previous option when focus is not a text search field |
| Down / `j` | Next option when focus is not a text search field |
| Home | First visible/available option |
| End | Last visible/available option |
| PageUp / PageDown | Move one viewport |
| Type characters | Filter when the list is declared searchable |
| `/` | Focus search in read-only pickers; insert `/` in ordinary text fields |
| Enter | Select highlighted option and continue |
| Escape | Clear search first, then go back on a second press |
| Tab / Shift+Tab | Move between search, list, details, and actions |

`j`, `k`, `q`, `/`, and `?` shortcuts MUST be active only outside editable text. Every letter shortcut requires an arrow/key alternative and visible help.

### Multi-select

- Space toggles the highlighted item.
- Enter accepts the current selection.
- Tab/Shift+Tab changes focus, not selection.
- `a` MAY select all and `n` clear all only when shown in help and safe for the list size.
- The UI always displays selected count and any min/max constraint.
- Disabled options include a textual reason.
- Screen-reader mode uses numbered choices and explicit `selected`/`not selected` announcements.

### Full-screen views

Full-screen alternate-buffer UI is reserved for workflows that materially benefit from persistent spatial context: large search, logs, diffs, or multi-pane inspection. Creation questions and confirmations MUST remain inline.

Full-screen views MUST:

- display a persistent one-line key legend or `?` help;
- restore cursor, echo, paste mode, mouse mode, and alternate screen on every exit path;
- preserve content on resize and resume;
- offer a line-oriented accessible fallback with equal outcomes;
- avoid mouse-only operations;
- use `q` to quit only in read-only views, never while editing text;
- never require enhanced keyboard protocols for core actions.

## Prompt behavior

### General rules

- Ask only questions whose answers cannot be safely derived.
- Use ecosystem-native words, not Aruo model names.
- Show a realistic example or concise consequence for unfamiliar inputs.
- Put the recommended choice first and label it `Recommended` when alternatives have tradeoffs.
- Pre-fill defaults only when accepting them is safe and unsurprising.
- Validate locally on submission; do not flash errors on every keystroke unless validation is cheap and non-disruptive.
- Preserve the user's input after an error and place the cursor near the problem where possible.
- Group related questions; show progress such as `Step 2 of 4` only when the count is stable.
- Allow review/back navigation before effects.
- End consequential flows with a target-specific summary or plan.

### Confirmation classes

| Risk | Interaction |
| --- | --- |
| Read-only or trivially reversible | No confirmation |
| Creates new state at an empty target | Summary plus `[Y/n]` only if consequences are material; otherwise proceed |
| Mutates existing local state | Display plan/diff, default `[y/N]` |
| Remote, destructive, security-sensitive, or hard to reverse | Name exact target and consequences; default no; typed target may be required |
| Batch/fleet effect | Display count, scope, samples, and saved plan ID; explicit approval required |

Confirmation syntax is always unambiguous:

```text
Apply these 4 changes to github.com/acme/api? [y/N]
```

Uppercase marks the default. Enter accepts the default. Localized yes/no tokens MAY be supported later, but `y`, `yes`, `n`, and `no` remain stable script-independent inputs in English mode.

`--yes` skips only confirmation prompts. It does not bypass trust refusals or provide missing choices. Destructive commands MAY require a separate explicit `--force` whose help names the lost safeguard.

### Selectors and fuzzy search

- Use a single-select list for 2–9 concise choices.
- Use fuzzy/searchable selection for longer or remotely loaded collections.
- Show stable identifiers alongside ambiguous names.
- Search results debounce network work and cancel stale requests.
- Never reorder under the cursor without preserving selected identity.
- Paginated lists expose loading and total/partial result status.
- An empty result explains how to broaden search or provide an explicit identifier.

### Optional questions

Optional inputs say `Optional` and accept Enter to skip. Do not encode “none” as an unexplained blank. Defaults show their source when not built in, for example `main (from git config)`.

### Prompt cancellation and EOF

Ctrl+C cancels the entire invocation. Escape goes back where possible. EOF is rendered as `Input ended; project creation was cancelled` rather than a generic scanner error. Cancellation is not printed as `Error:` unless cleanup or recovery failed.

## Progress reporting

### Choosing a progress form

| Work shape | Interactive terminal | Non-interactive human | Machine |
| --- | --- | --- | --- |
| <200 ms | No progress | No progress | Events only if part of schema |
| Indeterminate | Spinner plus current verb | Start/completion lines if useful | Structured lifecycle events |
| Known units | Progress bar with count | Periodic count or final summary | Structured progress events |
| Several sequential tasks | Task list/tree | One line per completed task | Task events |
| Concurrent tasks | Stable rows, bounded count | Timestamped completion lines | Events with task IDs |
| Streaming logs | Follow view | Append-only lines | JSONL when requested |

Progress describes useful work: `Downloading blueprint`, `Checking 142 files`, not `Please wait`. The final state replaces animation with a durable summary.

### Rendering rules

- Progress goes to stderr.
- Cursor rewriting occurs only when stderr is a terminal and motion is enabled.
- Rendering is capped at 10 frames/second; 4–8 is preferred for spinners.
- Resize events are coalesced; rows are redrawn without leaving stale text.
- Terminal width truncates the middle of paths/identifiers while preserving distinguishing suffixes.
- Concurrent progress has a bounded visible row count and a deterministic overflow summary.
- Completed tasks stop animating.
- Elapsed time appears after 2 seconds.
- ETA appears only with a defensible rate and enough samples; unstable ETA is omitted rather than misleading.
- `--quiet` suppresses progress and ordinary success prose, not errors or requested results.
- `--verbose` adds decisions and subprocess summaries; repeated `-v` MAY increase detail.
- `--debug` adds diagnostic events and a reference ID, with secrets redacted.
- `--progress auto|tty|plain|json|never` controls progress independently of result format.

In CI or redirected output, plain progress MUST be append-only and reasonably sparse. It MUST NOT emit carriage-return animation or ANSI cursor movement.

## Output design

### Formats

| Format | Audience | Contract |
| --- | --- | --- |
| `human` | Interactive people | Adaptive layout, restrained color/icons, may page |
| `plain` | Logs and simple shell reading | Line-oriented UTF-8/ASCII-safe, no ANSI, no paging or animation |
| `json` | Finite machine result | One schema-versioned JSON document |
| `jsonl` | Streams/events | One independently valid schema-versioned object per line |
| `yaml` | Human inspection/interchange | Same typed result as JSON; not the primary automation promise |
| `sarif` | Static-analysis integrations | Findings only, mapped from the canonical finding model |
| `markdown` | Reports/issues/docs | Explicit opt-in; no terminal ANSI |

Machine formats MUST write nothing else to stdout. Warnings and diagnostics stay on stderr, or appear in a documented structured event when the consumer requests an all-structured stream. Keys, enums, nullability, timestamps, paths, and ordering guarantees are versioned.

### Color, icons, Unicode, and emoji

- Color policy is `--color auto|always|never`.
- Non-empty `NO_COLOR` disables automatic color. An explicit command-line `--color=always` MAY override it; config does not silently override the environment preference.
- `TERM=dumb`, non-TTY output, plain/machine formats, and accessibility settings disable color unless explicitly forced.
- Color never carries the only indication of status; include words, shapes, or labels.
- Default palette must remain distinguishable under common color-vision deficiencies and in 4-bit terminals.
- Avoid faint text for required information.
- Unicode policy is `--unicode auto|always|never`; ASCII fallbacks exist for every symbol.
- Emoji are not structural UI. They MAY appear in celebratory human output only when Unicode is enabled, never in errors, tables, machine fields, or progress alignment.
- Icons are followed by text or have a textual accessible alternative.
- Hyperlinks are `auto|always|never`, with the visible URL available in plain/accessibility mode.

### Tables and wrapping

- Tables are for compact homogeneous records, not nested data.
- Human tables select columns by available width; omitted columns are disclosed with `--wide` or a hint.
- Plain tables use tabs only when explicitly requested; default plain records use stable key/value or one-record-per-line output.
- Do not parse human tables in scripts; help points to JSON or explicit templates.
- Numeric columns align right, text left, and headers remain textual.
- Cells wrap only when the layout remains readable; otherwise truncate with a visible marker.
- Paths and URLs are never silently changed.
- Paragraphs wrap to the terminal width capped at 100–120 columns; machine/plain data does not reflow.

### Paging

Paging is permitted only for finite human output sent to a terminal. It is disabled for prompts, progress, machine formats, CI, redirection, and errors. Resolve pager in documented precedence (`ARUO_PAGER`, then `PAGER`), and support `--pager`/`--no-pager`. Failure to start a pager falls back to direct output with a warning only in verbose mode.

## Terminal capability model

Aruo MUST detect and retain individual capabilities rather than a single `isTTY` field:

- stdin terminal status;
- stdout terminal status;
- stderr terminal status;
- width and height for the UI stream;
- color depth;
- Unicode preference/support;
- cursor addressing and erase support;
- alternate screen support;
- hyperlink support;
- interactive input eligibility;
- motion preference;
- screen-reader/accessibility mode;
- CI and dumb-terminal state;
- operating system and shell-relevant limitations.

Flags override Aruo-specific environment, which overrides generic environment hints, which overrides detection. Overrides are independently testable.

### Resize

On Unix, react to `SIGWINCH` or equivalent library events. On Windows, use console/terminal resize events supported by the UI library. Coalesce bursts within roughly 50 ms. If size is unavailable, assume 80 columns and avoid full-screen layout. A resize MUST NOT lose input, reset selection, duplicate output, or move focus.

### Non-interactive and CI

- Never prompt when stdin or stderr is not a usable terminal unless `--interactive` explicitly establishes a safe prompt channel.
- Machine format implies no input.
- CI disables animation, paging, browser opening, and update notices by default.
- Missing choices fail with one diagnostic listing flags/config keys needed.
- Authentication uses tokens/device-code instructions appropriate to the environment; no browser opening over headless SSH without an explicit request.
- Output remains append-only and timestamp-free by default so logs are diffable; streaming events carry timestamps in fields when semantically needed.

### SSH, tmux, screen, remote shells

Treat `$TERM` and terminfo as hints, not proof. Avoid terminal probing that can leak escape responses into input. Use conservative capabilities through unknown multiplexers. Device-code/browserless authentication must always be available. Latency-sensitive search debounces remote requests and shows cancellation. Never assume the local clipboard, GUI browser, or filesystem is available on the remote host.

### Windows

Tier-one support includes Windows Terminal and PowerShell. Aruo MUST also degrade cleanly in legacy console hosts:

- use Go/Windows console APIs or a proven library for VT mode and width;
- handle Ctrl+C and Ctrl+Break through the Go signal abstraction;
- accept CRLF and preserve path semantics;
- avoid relying on Unix `/dev/tty`, SIGWINCH, SIGTSTP, or executable bits;
- document PowerShell quoting where it differs from POSIX shells;
- use ASCII fallback if width/Unicode behavior is unreliable;
- test ConPTY behavior, redirected handles, and process-tree cancellation.

## Error and diagnostic design

### User-facing error anatomy

An error is rendered once and includes only applicable parts:

```text
Error: could not create my-library

The target already exists:
  /work/my-library

Nothing was changed.
Choose another name, or inspect the existing directory before retrying.

Debug reference: ARUO-7F3A
```

Required qualities:

- say what failed in user language;
- name the exact target and relevant source location;
- distinguish cause from consequence;
- state whether anything changed;
- give a safe next action, not generic “try again”;
- quote invalid values and show accepted shape/examples;
- group independent validation errors in one run;
- preserve the underlying causal chain for debug logs;
- never print secrets, authorization headers, environment dumps, or unredacted remote payloads.

Expected cancellation is `Cancelled`, not a red error. Expected policy findings use the finding renderer and exit code 3.

### Verbose and debug modes

- Default: concise result and actionable failure.
- `--verbose` / `-v`: configuration sources, selected adapters, cache decisions, external command summaries, retry notices, timings.
- repeated `-v`: progressively more detail only when documented.
- `--debug`: structured diagnostics, versions, stack trace for Aruo defects, HTTP metadata with bodies off/redacted, and debug reference.
- `ARUO_LOG` MAY set log level; command flags take precedence.
- Stack traces never appear by default and never replace the user-facing explanation.

Debug files use user-cache/state locations, restrictive permissions, bounded retention, and explicit paths. Users can disable diagnostic file writing.

### Retry behavior

- Retry only transient, idempotent operations.
- Show retry cause and delay after the first meaningful delay, not every fast attempt.
- Honor server retry guidance and cancellation.
- Use bounded exponential backoff with jitter.
- Never retry authentication denial, validation failure, trust refusal, or non-idempotent publication automatically.
- After exhaustion, show attempt count and exact safe retry command.

## Accessibility requirements

Accessibility is a supported mode, not a theme.

### Screen readers and braille displays

- `ARUO_ACCESSIBLE=1` and `--accessible` select line-oriented prompts and output.
- Accessible mode disables cursor-addressed redraw, spinners, transient status rows, alternate screen, ambiguous icons, and automatic paging.
- Each prompt is emitted once with its label, description, default, and valid input form.
- Selectors become numbered lists with explicit selected state.
- Validation errors are printed after submission and the prompt is repeated with context.
- Progress becomes sparse durable lines.
- Dynamic updates are not communicated solely by rewriting the same line.

### Keyboard-only operation

Every outcome is reachable with standard keyboard keys. Mouse support, if ever added, is optional. Focus order is stable and visible. No timeout applies to reading or answering prompts. Shortcut help is available without leaving the workflow.

### Motion and cognitive load

- `ARUO_MOTION=never` and `--motion=never` disable all animation.
- Accessible mode implies no motion unless explicitly overridden.
- There is no reliable universal terminal equivalent of web `prefers-reduced-motion`; Aruo therefore provides its own explicit control and SHOULD honor known ecosystem hints conservatively.
- Avoid flashing, high-frequency color changes, and indefinite animation.
- Use consistent verbs, positions, and confirmation grammar.
- Ask one conceptual question at a time.

### Color and contrast

Meet WCAG 2.2 AA contrast where colors are under Aruo's control, while acknowledging user terminal themes. Provide a 4-bit accessible palette and monochrome fallback. Never use red/green alone. Respect high-contrast and reverse-video environments by minimizing background fills.

### Internationalization

- All user-facing prose is localizable without concatenating grammatical fragments.
- Layout handles wider translations and double-width characters.
- Identifiers, flags, JSON keys, yes/no automation tokens, and debug codes remain stable.
- Error locations use Unicode-safe column semantics documented by format.
- Do not assume Latin sorting, word boundaries, or case folding for user content.

## Performance budgets

These are initial release budgets to validate on defined reference machines, not marketing claims.

| Interaction | Budget |
| --- | ---: |
| `aruo version` / cached help p95 | <= 50 ms process time on reference machine |
| First visible response for local interactive command p95 | <= 100 ms |
| Keypress-to-render p95 | <= 50 ms |
| Search/filter response for local list p95 | <= 100 ms |
| Ctrl+C acknowledgement p95 | <= 100 ms |
| Resize-to-stable-render p95 | <= 100 ms |
| Animation rendering | <= 10 frames/s |
| Idle interactive CPU | approximately 0%; no polling loop |
| Baseline resident memory for ordinary CLI paths | target <= 50 MiB; regression budget enforced after measurement |

Help/version MUST NOT initialize configuration, Git, plugins, network clients, update checks, or template catalogs unnecessarily. Prompt libraries and styling dependencies require binary-size and startup benchmarks before adoption. Remote latency is displayed after 2 seconds and is always cancellable.

## Implementation guidance for Go

This section guides implementation but does not select dependencies by itself.

### Process and signals

- Use `signal.NotifyContext` at the process composition root for interrupt/termination cancellation.
- After the first interrupt, unregister/reset interrupt handling so a subsequent interrupt can force the platform default, or implement an equivalent explicit second-signal path.
- Use context cancellation causes.
- Keep terminal restoration in an idempotent top-level defer and in the UI runtime's panic/error boundary.
- Abstract process-group signaling per platform; do not assume `os.Process.Signal(os.Interrupt)` works on Windows.
- Separate operation context from bounded cleanup context.

### Terminal services

Create an injected terminal service that owns stream handles and calculated capabilities. Use maintained primitives such as `golang.org/x/term` for terminal detection/size/raw mode, plus a proven width/grapheme implementation. Do not spread `os.Stdout`, ANSI literals, or environment reads across commands.

Use terminfo or conservative library capability detection for cursor operations. Never infer full capabilities solely from `TERM` substring matching. Cache detection for one invocation and re-read dimensions after resize/resume.

### Prompt/UI library policy

Simple Aruo workflows should use inline prompts, not a full-screen framework. A library is acceptable only if it supports:

- injected streams and contexts;
- accessible line mode;
- terminal restoration on signal/error;
- grapheme-aware editing and width;
- paste and resize handling;
- configurable key maps;
- no package-global state;
- deterministic rendering tests;
- Windows ConPTY;
- no telemetry/network behavior.

Charmbracelet Huh v2 is a candidate because it exposes input/select/multiselect/confirm and a first-class accessible mode. Bubble Tea v2 is a candidate only for future full-screen views; adopting it for basic creation prompts would add unnecessary rendering and signal complexity. Both require a focused dependency, accessibility, cancellation, startup, binary-size, and Windows evaluation before an ADR.

### Typed presentation model

Application services return typed results, findings, plans, progress events, and errors. Renderers transform them into human/plain/JSON/JSONL/YAML/SARIF/Markdown. Terminal styling never leaks into domain values. Structured schemas carry their own version.

## Testing strategy

### Unit and golden tests

- capability precedence across flags, Aruo environment, generic environment, TTY state, CI, and defaults;
- display width, grapheme editing, truncation, wrapping, table column choice, ASCII fallback, and redaction;
- key maps per widget and accessible-mode equivalents;
- prompt validation, defaults, back/review, EOF, and cancellation;
- exit-code mapping and typed error rendering;
- stable JSON/JSONL/SARIF schemas independent from human snapshots;
- progress throttling with a fake clock;
- terminal restoration functions are idempotent.

Golden tests must normalize only genuinely variable fields. Avoid snapshotting ANSI output without semantic assertions.

### PTY integration tests

Run the compiled binary under a pseudo-terminal and verify:

- first and second Ctrl+C behavior;
- terminal echo/cursor/raw mode restored after success, cancellation, error, and panic boundary;
- Ctrl+Z/SIGCONT on Unix;
- resize while entering text, searching, viewing a diff, and showing concurrent progress;
- redirected stdout with terminal stderr;
- piped stdin never triggers a prompt;
- plain output contains no ANSI or carriage-return redraw;
- prompt input and child process output cannot deadlock;
- Unicode, combining marks, emoji width, CJK width, and narrow terminals.

### Cross-platform matrix

Release qualification covers:

- Linux: bash/zsh, common VTE/xterm-compatible terminal, tmux, SSH, `TERM=dumb`;
- macOS: Terminal.app and a modern emulator, zsh, tmux, SSH;
- Windows: Windows Terminal with PowerShell, cmd where supported, ConPTY redirection, Ctrl+C/Ctrl+Break, resize;
- CI: GitHub Actions plus at least one non-GitHub provider or containerized non-TTY fixture;
- accessibility: NVDA or Narrator on Windows, VoiceOver on macOS, and a Linux screen-reader/line-mode review where practical.

### Fault and lifecycle tests

Inject cancellation at every operation boundary: before staging, during rendering, before rename, after local commit, during remote request, during retry delay, during verification, and during cleanup. Assert final filesystem/remote model, recovery record, locks, temp files, exit code, and message.

Test SIGKILL/power-loss recovery with subprocess fixtures rather than expecting cleanup code to run. Test child processes that ignore termination, close pipes slowly, emit large output, and spawn descendants.

### Performance tests

Benchmark cold/warm help and version, first prompt paint, filtering 10/1,000/100,000 choices, progress with 1/10/100 concurrent tasks, narrow/wide rendering, cancellation latency, idle CPU, allocations, memory, and binary-size delta for UI dependencies. Store raw samples and compare statistically; do not gate on noisy single runs.

## Definition of done for an interactive command

An interactive command is not production-ready until it has:

- a complete non-interactive equivalent;
- explicit stdin/stdout/stderr behavior;
- documented cancellation and commit boundary;
- first/second interrupt tests;
- accessible line-mode parity;
- no-color, no-Unicode, no-motion, dumb-terminal, narrow-terminal, and redirected-output tests;
- machine output schema where automation is expected;
- safe defaults and a reviewable plan for effects;
- redaction and debug behavior;
- Windows, macOS, and Linux qualification;
- startup, responsiveness, cancellation, memory, and binary-size evidence;
- help text listing keys, examples, defaults, external effects, and recovery.

## Source register

Primary and authoritative sources reviewed for this 2026 baseline:

- [GitHub CLI environment controls](https://cli.github.com/manual/gh_help_environment) and [formatting](https://cli.github.com/manual/gh_help_formatting)
- [Docker CLI](https://docs.docker.com/reference/cli/docker/), [BuildKit progress modes](https://docs.docker.com/reference/cli/docker/buildx/build/), and [TTY attachment](https://docs.docker.com/reference/cli/docker/container/run)
- [`kubectl` usage conventions](https://kubernetes.io/docs/reference/kubectl/conventions/), [global options](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_options/), and [`get` output](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_get/)
- [Terraform CLI](https://developer.hashicorp.com/terraform/cli/commands) and [plan input/JSON behavior](https://developer.hashicorp.com/terraform/cli/commands/plan)
- [Helm install safety options](https://helm.sh/docs/helm/helm_install/) and [Helm global terminal controls](https://helm.sh/docs/helm/helm_env/)
- [Cargo terminal configuration](https://doc.rust-lang.org/cargo/reference/config.html)
- [`uv` CLI controls](https://docs.astral.sh/uv/reference/cli/) and [help/paging/verbosity](https://docs.astral.sh/uv/getting-started/help/)
- [Bun install output controls](https://bun.sh/docs/pm/cli/install), [npm configuration](https://docs.npmjs.com/cli/using-npm/config/), [pnpm reporters and CI behavior](https://pnpm.io/cli/install), [Deno CLI/config reference](https://docs.deno.com/runtime/reference/), and [Vite creation behavior](https://vite.dev/guide/)
- [Fly CLI](https://fly.io/docs/flyctl/) and [Fly automation/JSON output](https://fly.io/docs/flyctl/integrating/)
- [Railway CLI modes](https://docs.railway.com/cli), [deploy/CI output](https://docs.railway.com/cli/deploying), [TTY search behavior](https://docs.railway.com/cli/templates), and [browserless login](https://docs.railway.com/cli/login)
- [Supabase CLI global flags](https://supabase.com/docs/reference/cli/usage)
- [Wrangler command behavior](https://developers.cloudflare.com/workers/wrangler/commands/workers/) and [structured diagnostic output](https://developers.cloudflare.com/workers/wrangler/system-environment-variables/)
- [Azure CLI output formats](https://learn.microsoft.com/en-us/cli/azure/format-output-azure-cli), [AWS CLI output formats](https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-output-format.html), [AWS paging](https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-pagination.html), [gcloud configurations](https://docs.cloud.google.com/sdk/gcloud/reference/topic/configurations), and [gcloud screen-reader behavior](https://docs.cloud.google.com/sdk/gcloud/reference/topic/accessibility)
- [Go signal behavior](https://pkg.go.dev/os/signal), [Windows Ctrl+C/Ctrl+Break](https://learn.microsoft.com/en-us/windows/console/ctrl-c-and-ctrl-break-signals), and [GNU Readline interaction](https://www.gnu.org/software/bash/manual/html_node/Readline-Interaction.html)
- [NO_COLOR](https://no-color.org/), [terminfo](https://invisible-island.net/ncurses/man/terminfo.5.html), [Huh accessible prompts](https://github.com/charmbracelet/huh), and [Bubble Tea v2 terminal model](https://github.com/charmbracelet/bubbletea/releases)

## Final principle

The terminal is a shared, stateful, imperfect interface. Aruo earns trust by leaving it usable, making every effect interruptible and recoverable, presenting the same semantics in rich and plain modes, and never requiring a visual trick to understand what happened.
