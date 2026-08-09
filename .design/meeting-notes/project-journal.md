# Aruo project journal

This page preserves the product, architecture, implementation, and release decisions made while creating Aruo. Early entries were reconstructed on 2026-08-09 from accepted ADRs, research, documentation, implementation, and Git history. They record the evidence that remains in the repository and label unknown discussion details instead of inventing a transcript.

---

## 2026-08-04 — Foundation decisions

- Date: 2026-08-04
- Recorded: 2026-08-09 (retrospective reconstruction)
- Participants: Aruo project lead; no further attendance record survives
- Sources: ADR-0001 through ADR-0009, architecture and product research, initial repository baseline

### Context

Aruo needed a durable product boundary before implementation began. The central
question was whether it should be another project generator or a broader system
for understanding and safely evolving repositories. The bootstrap also needed to
choose a core language, extension boundary, template model, license, contributor
toolchain, CLI framework, and rendering contract without locking the domain model
to any one framework.

### Evidence reviewed

- Existing generators are strongest at creation but do not preserve intent,
  provenance, drift explanations, or safe upgrades.
- Repository files, native ecosystem tools, and CI remain the practical sources
  of truth; a mandatory service would add availability, privacy, and lock-in risk.
- Go and Rust both fit a portable systems CLI. Go reduced delivery and contributor
  complexity while retaining fast startup, static distribution, cross-compilation,
  concurrency, and strong standard tooling.
- In-process extensions would inherit the CLI's filesystem, memory, and credential
  authority and couple third parties to Go APIs.
- Monolithic templates make combinations and upgrades difficult to reason about;
  unrestricted executable generators make effects difficult to inspect.
- Ecosystems share desired outcomes but not idiomatic layouts or toolchains.
- A constructor-built Cobra tree could provide mature parsing while keeping the
  framework out of application and domain packages.
- A pure `text/template` + `io/fs` renderer could produce a deterministic plan
  without gaining write, network, environment, clock, or random access.

### Decisions

1. Define Aruo as a local-first repository lifecycle control plane. Project
   creation is the first useful slice, not the long-term product boundary.
2. Use Go for the core modular monolith. Keep schemas language-neutral and public
   packages deliberate; most implementation remains under `internal/`.
3. Run third-party plugins as child processes over a versioned JSON Lines protocol
   with explicit capabilities, cancellation, and limits. Revisit WASM only when a
   stricter portable sandbox is justified.
4. Compose constrained foundation, language, workload, capability, and
   organization blueprints. Native ecosystem tools remain authoritative for
   ecosystem structure.
5. Standardize observable outcomes and evidence, not cosmetic repository trees.
   Exceptions must be scoped, owned, justified, and time-bounded.
6. License Aruo under Apache-2.0. Generated repositories choose their own reviewed
   license; using Aruo does not impose Aruo's license on user projects.
7. Bootstrap with Go 1.26, Make, Go-native checks, GitHub Actions, Dependabot,
   Conventional Commit PR titles, release-please, and GoReleaser. Avoid a second
   task runner and mandatory local hooks until evidence supports them.
8. Contain Cobra inside the presentation layer, build a fresh command tree for
   every invocation, inject dependencies and streams, and render returned errors
   once at the CLI boundary.
9. Keep built-in template rendering pure and deterministic. It returns a sorted
   in-memory file plan; writing, conflict policy, discovery, and structured edits
   belong to separate layers.

### Tradeoffs accepted

- Local completeness defers fleet-wide history and collaboration.
- Go trades Rust's stronger type and memory guarantees for simpler delivery and
  contribution.
- Process isolation adds serialization and lifecycle overhead.
- Constrained composition and a small renderer helper vocabulary require more
  deliberate core design than arbitrary generator code.
- Make is not equally native on every Windows setup; CI and the dev container
  provide the reproducible fallback.

### Unresolved questions and gates

- Prove startup, binary-size, cross-compilation, deterministic-plan, and schema
  acceptance criteria before expanding the implementation.
- Define plugin signing, permission brokerage, and conformance before accepting
  third-party execution.
- Measure blueprint conflict rates on repositories that have realistic user edits.
- Keep hosted collaboration, WASM plugins, and semantic structured-file adapters
  behind later evidence and explicit decisions.

### Resulting records

- [ADR-0001](../../decisions/0001-local-first-repository-control-plane.md)
  through [ADR-0007](../../decisions/0007-repository-toolchain-and-automation.md)
- [ADR-0008](../../decisions/0008-constructor-built-cobra-cli.md)
- [ADR-0009](../../decisions/0009-pure-standard-library-template-renderer.md)
- [Product thinking](../product.md) and [working vision](../vision.md)

---

## 2026-08-05 to 2026-08-06 — Terminal experience implementation

- Date: 2026-08-05 to 2026-08-06
- Recorded: 2026-08-09 (retrospective reconstruction)
- Participants: Aruo maintainers; individual attendance was not recorded
- Sources: terminal-stack research, ADR-0010, terminal UX benchmark, commits implementing `internal/tux`

### Context

The CLI needed to feel deliberate in an interactive terminal while remaining
predictable in pipes, CI, `TERM=dumb`, narrow terminals, screen-reader workflows,
SSH/tmux, and Windows. The implementation also needed one owner for cancellation
and terminal restoration. Adopting a full-screen application model would have made
ordinary CLI composition and machine use harder.

### Evidence reviewed

- The 2026 Go terminal ecosystem evaluation favored Charm v2 for rich interaction,
  but only behind an Aruo-owned adapter boundary.
- Huh does not provide an in-place degraded interaction path when a terminal cannot
  address the cursor.
- Raw terminal mode can consume Ctrl+C before it reaches an OS signal handler.
- The measured Charm-family binary cost was about 2.93 MiB stripped, below the
  5 MiB review threshold; version execution remained well below the 50 ms budget.
- CI, `NO_COLOR`, redirected streams, accessible mode, and narrow terminals require
  behavior based on capabilities rather than assumptions about an operating system.

### Decisions

1. Define Aruo-owned semantic contracts for prompts, presentation, progress, and
   sessions. Command code never depends directly on Charm types.
2. Use Huh, Bubble Tea, Bubbles, and Lip Gloss v2 as replaceable rich adapters,
   with `x/term` for terminal primitives.
3. Maintain a dependency-free plain adapter as the semantic reference path, not a
   second-class fallback. Select it for non-interactive/machine modes, accessible
   mode, and terminals without cursor addressing.
4. Detect terminal capabilities once per session and ensure color and icons never
   carry meaning alone.
5. Let Aruo's lifecycle manager own signals and cancellation. Disable Bubble Tea's
   signal handler and normalize UI cancellation with context cancellation.
6. Model progress as task events and cap rich rendering at 10 FPS. Durable plain
   output remains sparse and readable rather than simulating animation.
7. Keep help, completion, and version paths independent of terminal-session
   construction. Add Cobra shell completion as a conventional CLI surface.

### Consequences

- Rich and plain behavior can be tested against the same semantic contracts.
- The CLI has a deterministic accessibility and automation path even if the rich
  dependency family changes.
- Dependency concentration in Charm is accepted, with exact pins, a plain adapter,
  conformance tests, benchmarks, and an explicit replacement exercise as controls.
- Structured JSONL progress, Markdown rendering, and an explicit `/dev/tty`
  interactive mode remain future work rather than implied current capability.

### Follow-ups

- Re-run adapter conformance at least once per major release cycle.
- Review normalized benchmark regressions above 20% before growing the rich stack.
- Periodically prove the adapter boundary by spiking a representative replacement.
- Keep documentation explicit about implemented versus planned terminal features.

### Resulting records

- [ADR-0010](../../decisions/0010-charm-v2-terminal-ux-stack.md)
- [Terminal UX](../../docs/cli/terminal-ux.md)
- [Technology evaluation](../../research/technology/2026-go-terminal-ux-stack.md)
- [Terminal benchmark](../../benchmarks/results/2026-08-06-terminal-ux-baseline.md)

---

## 2026-08-06 — Documentation and first catalog expansion

- Date: 2026-08-06
- Recorded: 2026-08-09 (retrospective reconstruction)
- Participants: Aruo maintainers; individual attendance was not recorded
- Sources: documentation reorganization commits, architecture corrections, JavaScript library template and tests

### Context

The initial repository had many root-level documents, several empty placeholder
directories, and architecture text that mixed implemented code with the intended
lifecycle pipeline. At the same time, the built-in catalog needed to prove that
the renderer and create flow were not accidentally specific to Go.

### Decisions

1. Organize public documentation by reader topic under `docs/` and use README files
   as section entry points. Keep ADRs, RFCs, research, benchmarks, and root project
   contracts in their distinct repositories of record.
2. Remove empty stubs instead of preserving an aspirational tree. Add a directory
   when it has maintained content or a concrete ownership need.
3. Label architecture as implemented, partial, or planned. Do not let diagrams or
   prose imply that roadmap modules already exist.
4. Use `--no-input` as the non-interactive public vocabulary; retire the ambiguous
   `--non-interactive` wording in documentation.
5. Document only live terminal capabilities. Secret input, multi-select, and rich
   progress were not to be claimed until their actual path and tests existed.
6. Add a JavaScript library as the second built-in catalog entry to exercise a
   non-Go ecosystem through the same renderer, planner, create service, and doctor
   contracts.
7. Treat template fixtures and generated output as executable product surfaces:
   test rendering, project checks, and doctor detection instead of relying on file
   presence alone.

### Tradeoffs accepted

- Topic-owned docs create more directories but give each subject one canonical
  entry point and reduce duplicate root documents.
- Deleting empty stubs makes the roadmap less visually prominent but keeps the
  repository honest about what is maintained.
- A second template increases qualification work immediately, which is useful
  pressure on generic catalog and doctor behavior.

### Unresolved questions and follow-ups

- Continue expanding the catalog only where each template can be generated and
  tested as a real project.
- Keep machine-output and progress contracts versioned before presenting them as
  stable automation APIs.
- Revisit documentation navigation when the corpus is large enough to justify a
  generated site; Markdown remains the source of truth.

### Resulting records

- [Documentation map](../../docs/README.md)
- [Architecture overview](../../docs/architecture/README.md)
- [CLI contract](../../docs/cli/README.md)
- [Template documentation](../../docs/templates/README.md)
- [JavaScript library template](../../internal/templateengine/builtin/templates/js)

---

## 2026-08-07 — Guided create experience

- Date: 2026-08-07
- Recorded: 2026-08-09 (retrospective reconstruction)
- Participants: Aruo maintainers; individual attendance was not recorded
- Sources: catalog expansion, guided-flow commits, CLI copy guide, interaction tests

### Context

The catalog expanded from two library templates to applications and libraries
across JavaScript, TypeScript, Python, React, Nuxt, Vue, Next.js, and Go. A flat
picker and a shared “module” question no longer matched how users think about
starting a project. Early prompt copy also repeated descriptions and placed
placeholder or default text inside editable inputs, creating extra deletion and
decision work.

### Decisions

1. Ask whether the user is creating an application or a library before showing
   templates. Group catalog entries by that stable user intent rather than exposing
   one growing flat list.
2. Give every catalog entry a concise, parallel description. Show that description
   once at the choice point; do not repeat it on later screens.
3. Keep editable inputs empty. Defaults and examples may be shown as guidance but
   must not be inserted as text the user has to erase.
4. Derive useful values from stronger inputs instead of presenting `TODO`
   placeholders. Description and author remain optional and may be blank.
5. Stop calling a cross-ecosystem identifier a “Go module path.” Prompt and help
   copy must use vocabulary appropriate to the chosen template.
6. Assign restrained per-entry colors in the rich picker for scanning, while text,
   order, and selection state continue to carry meaning without color.
7. Allow users to move backward across guided screens without discarding the whole
   session. Back navigation is part of the prompt contract, not an ad hoc command
   workaround.
8. Split durable wording rules into a CLI copy style guide so new commands inherit
   the same concise, calm language.

### Why this shape

The flow should progressively disclose decisions. Project kind narrows the catalog;
the selected template then determines which metadata is relevant. This avoids
asking every possible question up front and keeps the common path coherent as the
catalog grows.

### Unresolved questions and follow-ups

- Test catalog grouping and back navigation in plain as well as rich interaction.
- Avoid adding categories until the user distinction is stable and materially
  reduces choice cost.
- Continue removing metadata prompts that do not affect generated output or cannot
  be explained in ecosystem-native language.

### Resulting records

- [CLI copy style guide](../../docs/cli/copy-style-guide.md)
- [Terminal UX contract](../../docs/cli/terminal-ux.md)
- [Catalog implementation](../../internal/catalog)
- [Create architecture](../../docs/architecture/create-command.md)

---

## 2026-08-08 to 2026-08-09 — Public release and template scope

- Date: 2026-08-08 to 2026-08-09
- Recorded: 2026-08-09 (retrospective reconstruction)
- Participants: Aruo maintainers; individual attendance was not recorded
- Sources: release preparation commits, v0.1.0 release, README and roadmap corrections, template additions

### Context

Preparing the first public release exposed decisions that had been tolerable during
bootstrap but confusing as user-facing contracts: the Go module identity differed
from the public CLI repository, project creation asked for a package/module name
that templates could derive, existing empty directories were rejected, public docs
still described Phase 0, and release automation initially interpreted the first
release as 1.0.0.

### Decisions

1. Use `github.com/aruodore/aruo-cli` as the canonical Go module identity while the
   product and command remain Aruo / `aruo`.
2. Remove the module/package-name prompt from guided creation. Derive ecosystem
   identifiers from the destination/project name until a template demonstrates a
   real need for a separate value; advanced overrides can be added deliberately.
3. Permit generation into an existing directory only when it is empty. Continue to
   reject a non-empty target so create cannot silently merge or overwrite user work.
4. Present the actual release scope honestly: project creation is implemented;
   broader audit, planning, upgrade, policy, and plugin ambitions remain roadmap
   work.
5. Ship a qualified nine-template catalog spanning app and library workloads.
   Add Vue as an application template and include Nuxt UI in the Nuxt application
   baseline so those starters represent deliberate supported choices rather than
   empty framework shells.
6. Keep the first public version at `0.1.0`. Configure release-please with an
   explicit initial version and `0.x` strategy instead of accepting an accidental
   1.0 semantic-stability signal.
7. Let an approved release PR create and publish artifacts automatically, but mark
   the v0.x GitHub releases as pre-releases while compatibility is still forming.
8. Treat release-facing prose as a factual interface. Replace stale pre-release
   claims after publication, use real CI/release badges, and remove wording that
   sounds inflated or machine-generated.
9. Remove `.github/README.md` because GitHub could surface it instead of the real
   repository README. There should be one canonical public landing page.

### Tradeoffs accepted

- Derivation makes the common creation path shorter but postpones unusual package
  naming until there is an explicit override contract.
- Empty-directory support is convenient for users who pre-create folders while
  retaining a simple safety boundary; merging into populated repositories belongs
  to a future planned/apply workflow.
- Automatic publication reduces release toil but makes release configuration and
  protected credentials part of the trusted delivery path.
- A broader template catalog raises maintenance cost; every catalog entry therefore
  remains subject to generation and qualification tests.

### Follow-ups

- Define compatibility and graduation criteria before removing the pre-release
  designation or approaching 1.0.
- Keep README, roadmap, changelog, and release metadata synchronized with shipped
  behavior after every release.
- Measure whether users need an explicit package/module override before adding the
  prompt back in another form.
- Keep populated-directory adoption separate from create until preview, conflict,
  provenance, and rollback contracts are implemented.

### Resulting records

- [v0.1.0 changelog](../../CHANGELOG.md)
- [Public README](../../README.md)
- [Roadmap](../../ROADMAP.md)
- [Template catalog](../../docs/templates/README.md)
- [Create command architecture](../../docs/architecture/create-command.md)

