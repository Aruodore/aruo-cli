# The state of developer tooling in 2026

**Research date:** 5 August 2026  
**Status:** Evidence snapshot; architectural input, not an accepted decision  
**Scope:** Tools used to create, operate, test, document, release, extend, and govern software repositories

## Executive summary

Developer tooling is stronger than it has ever been at individual jobs and still weak at preserving intent across a repository's life.

Fast, integrated tools such as `uv`, Ruff, Bun, Deno, Cargo, and the Go toolchain reduce installation and coordination costs. Vite and native initializers make day-zero creation pleasant. GitHub Actions, Dependabot, release-please, Changesets, and GoReleaser can automate most of a delivery pipeline. Documentation systems provide excellent search, navigation, API generation, and interactive content. Mature plugin systems make tools adaptable.

Yet developers still assemble these capabilities into a lifecycle by hand. A generator rarely knows how to update the repository it created. An auditor usually knows that a file exists, but not why it differs. Configuration layers hide provenance. Release automation infers user impact from commit syntax. CI policy is copied as YAML and drifts. Plugins often execute with the user's full authority. A passing badge can conceal an unowned or unreleasable project.

The central opportunity for Aruo is therefore not a larger template catalog. It is a durable, inspectable lifecycle model connecting:

```text
intent -> plan -> create/adopt -> verify -> maintain -> migrate -> release
            |          |             |          |             |
            +---------- provenance and evidence --------------+
```

Aruo should orchestrate ecosystem-native tools, preserve why each managed artifact exists, calculate changes before applying them, and emit the same typed evidence to humans, CI, editors, and automation. It should not become a universal package manager, build system, or unrestricted hook runner.

## Method and confidence

This study reviewed current official documentation, specifications, and source repositories. Recent behavior was preferred over historical reputation. Successful projects were treated as evidence of patterns, not proof that every choice transfers to Aruo.

The report distinguishes:

- **Observation:** directly supported by a tool's current behavior or documentation.
- **Inference:** a conclusion drawn across observations.
- **Implication:** a constraint or opportunity for Aruo to validate through an ADR, RFC, prototype, or benchmark.

The 2026 ecosystem will continue to move. Exact tool versions and vendor features have a short half-life; architectural lessons should outlive them. Market share and maintainer sentiment were not measured quantitatively, so adoption claims are intentionally qualitative.

## Forces shaping the 2026 ecosystem

### Integrated toolchains are winning attention

Cargo established the value of one coherent path for build, test, dependencies, documentation, and publishing. `uv` now applies that pattern to Python with project initialization, a universal lockfile, environments, workspaces, tool execution, and an append-only concurrent cache. Ruff consolidates linting and formatting. Bun and Deno combine runtimes with package, test, and development capabilities.

The advantage is not merely speed. A consolidated tool owns fewer hand-off boundaries, provides one configuration vocabulary, and can make stronger reproducibility guarantees. The risk is ecosystem capture: an integrated tool may become a second platform that fights native conventions.

**Aruo implication:** integrate at the lifecycle and evidence layers. Delegate compilation, resolution, formatting, and publishing to authoritative native tools.

### Creation is becoming a maintenance problem

Vite offers both prompts and deterministic template selection. `uv init` distinguishes applications, libraries, packaged applications, extension modules, and deliberately minimal projects. Copier preserves answers and supports later updates. These are advances over one-shot file copying.

But updating a generated repository remains difficult because users legitimately edit generated files, templates evolve, and hooks introduce effects that are not represented in the rendered tree. Textual three-way merging cannot reliably distinguish customization from drift.

**Aruo implication:** store blueprint identity, version, inputs, ownership boundaries, and content fingerprints. Upgrade from a semantic plan, never by silently regenerating a directory.

### Supply-chain evidence is moving into the default path

GitHub Actions supports short-lived OIDC credentials and artifact attestations; public attestations use Sigstore transparency infrastructure. GoReleaser supports checksums, SBOMs, signing, and attestations. Dependency review can reject newly introduced vulnerable dependencies before merge. Dependabot now supports cooldowns as well as grouped updates.

Automation configuration is itself executable supply-chain code. Mutable action references, broad tokens, untrusted pull-request data, inherited secrets, and opaque installers remain common failure modes.

**Aruo implication:** secure workflows are generated policy, not decorative files. Plans should expose permissions, external actions, immutable identities, network effects, and provenance.

### Humans are no longer the only consumers

CI, IDEs, repository dashboards, policy engines, and coding agents all consume tool output. A colored terminal transcript is not an API. Conversely, raw JSON is not humane guidance.

**Aruo implication:** one typed result model should feed terminal output, stable versioned JSON, SARIF where appropriate, and documentation. Machine output must not be scraped from prose.

## Domain findings

### Go CLI development

Go remains a strong implementation language for a cross-platform engineering CLI: static distribution is straightforward, startup is fast, concurrency is built in, and the standard library covers HTTP, archives, templates, structured logging, testing, fuzzing, embedding, and filesystem abstractions.

Go 1.26, released in February 2026, strengthens the toolchain rather than changing its character. The redesigned `go fix` uses the analysis framework shared with `go vet`; Green Tea GC is enabled by default; and `go mod init` intentionally defaults newly created modules to the preceding Go language version to encourage a supported compatibility floor. The Go project supports a major release until two newer major releases exist.

Cobra remains the conventional command-tree framework and is used by large CLIs such as GitHub CLI, Hugo, Helm, and Kubernetes-adjacent tools. Its maturity, completion support, and documentation are valuable. Its package-global examples and generator conventions are not an architectural mandate. Large maintainable CLIs generally construct a root command around explicit dependencies, streams, context, and services.

Patterns worth adopting:

- a thin executable entry point and a command layer that translates CLI inputs into application requests;
- constructor injection for filesystem, clock, environment, process execution, terminal capabilities, network clients, and build metadata;
- `context.Context` across cancellable I/O and subprocess boundaries;
- separate stdout for requested results and stderr for progress, diagnostics, and errors;
- terminal capability detection with `NO_COLOR`, non-TTY, reduced-motion, and plain-output behavior;
- stable exit categories, causal error chains, and remediation attached to errors;
- lazy initialization so `help` and `version` never load configuration, plugins, Git, or the network;
- no process termination, ambient `os.Stdout`, or global logger below the outer command boundary.

Recurring weaknesses:

- command frameworks leak into domain packages;
- validation occurs during side effects rather than before them;
- human output is later treated as a machine interface;
- hidden environment and working-directory dependencies make tests brittle;
- completion, accessibility, and Windows behavior arrive late.

**Direction for Aruo:** Cobra is an acceptable adapter, not the architecture. Define a framework-neutral application core and performance budgets for cold `help`, `version`, and ordinary local commands.

### Repository generators

The ecosystem divides into four useful models:

| Model | Examples | Strength | Structural weakness |
|---|---|---|---|
| Language-native initializer | Cargo, `go mod init`, `uv init`, Poetry | Least surprising ecosystem output | Narrow repository lifecycle scope |
| Framework initializer | Vite, create-next-app, create-t3-app | Excellent contextual prompts and immediate runnable result | Option combinations and framework churn |
| General template renderer | Cookiecutter, Hygen, Plop, Yeoman | Flexible and easy to author | Weak provenance and update semantics |
| Update-aware copier | Copier | Answers and template version enable updates | Text merge conflicts and hook trust remain |

Strong creation experiences share progressive disclosure, a small default path, explicit non-interactive flags, early validation, a summary before destructive work, and a useful next command. `uv init` is especially instructive because it models project kinds rather than pretending every Python repository has one layout. Vite exposes deterministic template selection and a no-prompt mode rather than making automation drive an interactive UI.

The most serious weakness is unrestricted hooks. Cookiecutter pre/post hooks may execute Python or shell. This is powerful, but installing a template can become equivalent to running unreviewed code with workstation authority.

**Direction for Aruo:** templates should be declarative content plus a typed manifest. Any effectful extension must be separately declared, reviewable, capability-limited, and disabled or explicitly allowed in unattended execution. Language packs should invoke native initializers where they are authoritative, then layer repository standards through a plan.

### Documentation systems

Docusaurus, VitePress, MkDocs Material, Nextra, Fumadocs, and hosted products compete mostly on authoring model and integration, not on whether polished navigation and dark mode are possible.

Current patterns:

- Markdown/MDX remains the portable source format;
- local search is viable for modest sites (VitePress uses MiniSearch; MkDocs Material has an integrated search plugin);
- OpenAPI-driven reference generation is a first-class feature in Fumadocs and similar systems;
- code tabs, copy controls, heading links, and responsive navigation are baseline expectations;
- static output and incremental content loading protect performance;
- executable examples and generated references reduce drift better than editorial discipline alone.

Docusaurus explicitly warns that versioned documentation increases build complexity and contributor burden, recommends keeping few versions, and notes that most projects do not need it. This is a valuable correction to “version everything.”

Recurring weaknesses:

- selecting a theme is mistaken for designing an information architecture;
- API reference, conceptual explanation, and tutorials are mixed together;
- copied examples rot;
- docs versioning is enabled without a support policy;
- search hides weak navigation rather than complementing it;
- generated sites are accessible to browsers but not structured for other consumers.

**Direction for Aruo:** define content contracts—tutorial, how-to, explanation, reference, troubleshooting—and test links, snippets, examples, and command reference generation. Recommend versioned docs only when supported versions genuinely differ. Documentation output should remain portable and useful without Aruo.

### Release tooling and semantic versioning

The strongest release tools divide responsibility clearly:

- **release-please** derives proposed versions and changelogs from Conventional Commits, keeps a reviewable release PR current, then tags and creates a GitHub Release after merge. It deliberately does not publish packages.
- **Changesets** captures release intent and human-readable impact while the change is being made, then composes package bumps and internal dependency updates—especially effective in monorepos.
- **semantic-release** favors fully automatic commit-derived releases; this reduces ceremony but moves correctness pressure to commit history and configuration.
- **GoReleaser** specializes in reproducible multi-platform packaging and distribution, with checksums, SBOMs, signing, attestations, archives, package-manager metadata, and container artifacts.
- **GitHub Releases** is a distribution and communication surface, not by itself a release policy.

SemVer 2.0.0 is useful only after a project defines its public API. For Aruo that API is broader than Go symbols: CLI flags and exit behavior, configuration schema, JSON/SARIF output, plugin protocol, template manifest, managed repository semantics, and migration guarantees can all break consumers. A version number cannot substitute for a support matrix, deprecation window, or migration guide.

Commit syntax is also an imperfect proxy for user impact. A small internal commit can require a config migration; a large feature can be backward compatible. Reviewable release intent is stronger than pure inference.

**Direction for Aruo:** separate “prepare release,” “build artifacts,” and “publish.” Capture change intent at pull-request time, produce a human-reviewed release proposal, build only from an immutable tag, and attach checksums, SBOMs, attestations, and verification instructions. Define versioned compatibility surfaces before 1.0.

### GitHub Actions

GitHub Actions is ubiquitous and capable, but YAML convenience can obscure a distributed execution and authorization system.

Modern baseline practices include:

- explicit least-privilege `permissions` at workflow/job level;
- immutable full-SHA references for third-party actions, updated by automation;
- OIDC federation instead of long-lived cloud credentials;
- protected environments for publication and deployment;
- concurrency groups to cancel obsolete validation while serializing releases;
- dependency review on pull requests and CodeQL where language support and risk justify it;
- artifact attestations and SBOMs for released executables;
- reusable workflows for organization policy, with awareness that called workflows execute in the caller's context;
- separate untrusted pull-request validation from privileged release paths;
- timeouts, deterministic caches, and an intentionally small matrix.

GitHub notes that attestations provide benefit only when consumers verify them. This illustrates a broader problem: producing evidence is not the same as making it operational.

Recurring weaknesses:

- copied workflows drift across repositories;
- action tags look readable but are mutable;
- cache keys and permissions are cargo-culted;
- release workflows mix building, signing, and publishing under one token;
- fork pull requests are handled as trusted inputs;
- generated YAML has no clear owner or upgrade path.

**Direction for Aruo:** generate minimal repository-local callers plus inspectable policy, or fully local workflows when external coupling is undesirable. `doctor` should explain effective permissions, mutable dependencies, trust boundaries, stale runtimes, and missing verification—not merely check that `.github/workflows` exists.

### Testing frameworks and strategy

Go's standard `testing` package now covers unit tests, subtests, examples, benchmarks, fuzzing, and testable filesystem implementations. New benchmarks should use `B.Loop`, which the current documentation describes as more robust and efficient than the older `b.N` style. Native fuzzing retains failures as regression corpus entries. The race detector and fuzzing address different classes of defects and belong in different cadence tiers.

For web-facing generated projects, Playwright demonstrates a useful integrated model: runner, assertions, isolation, parallel execution, multi-browser support, traces, and rich reports. Its fresh browser context per test makes isolation a default rather than a convention.

The right testing portfolio follows risk boundaries:

- unit tests for pure policy and transformation logic;
- black-box command tests for public CLI behavior;
- golden tests for plans and rendering, with explicit review of updates;
- contract tests for schemas and plugin protocols;
- property/fuzz tests for paths, merges, parsers, and hostile template inputs;
- integration tests for Git, subprocesses, and real filesystems;
- end-to-end tests for installed release artifacts on tier-one operating systems;
- mutation or fault-injection experiments selectively for high-risk policy logic.

Coverage is diagnostic evidence, not the target. A high percentage cannot establish meaningful assertions, platform behavior, or release usability.

**Direction for Aruo:** every engine should expose deterministic seams and test doubles, but filesystem security and cross-platform semantics must also be exercised against real operating systems. Tests should verify explanations and remediation as public behavior.

### Benchmarking

Fast CLIs earn trust through latency consistency, not a single best-case number. Go's benchmark format and `benchstat` support statistically robust A/B comparison. Benchmarks can report allocations as well as time.

A credible program records tool version, commit, OS, architecture, CPU, filesystem, cache state, corpus, sample count, and variance. It separates:

- microbenchmarks for parsers, matching, rendering, and hashing;
- component benchmarks for repository discovery and plan construction;
- end-to-end cold and warm command latency;
- scalability corpora covering repository size, file count, plugin count, and template count;
- memory and binary-size budgets;
- comparison with the native workflow Aruo claims to improve.

Noisy shared CI should usually report statistically significant changes rather than fail on tiny percentage thresholds. Stable dedicated runners may enforce broad regression budgets.

**Direction for Aruo:** publish methodology and raw results, not marketing-only numbers. Prioritize cold no-op commands, bounded memory during scans, cancellation, and avoiding unnecessary network access. Cache only derived data with explicit keys and safe invalidation.

### Dependency management

Package managers increasingly treat resolution state as a versioned interface. Go modules use minimal version selection and record toolchain requirements. `uv.lock` is universal across platforms, has a versioned schema, and is expected in version control; its cache is append-only and concurrency-safe. Dependabot's cooldown capability acknowledges that immediate adoption of every new version is not always the safest policy.

Useful distinctions:

- applications and development tools need repeatable exact resolution;
- libraries need a deliberately broad, tested compatibility claim;
- generated repositories should not introduce a second dependency authority;
- tool dependencies should be source-controlled using the ecosystem's supported mechanism;
- security updates need a faster path than routine churn;
- update batching should preserve reviewability and bisectability.

Recurring weaknesses are PR floods, unreviewed transitive changes, lockfile conflicts, unsupported runtime drift, and automation that updates manifests without testing released artifacts.

**Direction for Aruo:** express update policy—cadence, cooldown, grouping, compatibility floor, security SLA—and render it into native mechanisms. Audit both declared and resolved dependencies, including GitHub Actions. Never invent an Aruo lockfile for language packages.

### Developer experience

The best tools make the safe path short while keeping automation explicit.

Common successful patterns:

- progressive disclosure: one strong default, advanced choices discoverable later;
- exact parity between prompt answers and non-interactive flags/configuration;
- preflight validation before network or filesystem mutation;
- plan/diff/dry-run for consequential work;
- idempotence where the operation semantics permit it;
- errors that name the failed operation, affected resource, cause, and next action;
- quiet success, useful summaries, and opt-in verbosity;
- offline or cache-aware behavior with clear network disclosure;
- shell completion, examples in help, and nearby next steps;
- accessible color, Unicode, motion, and prompt fallbacks;
- stable machine output with a schema version.

Recurring pain comes from hidden work, surprise authentication, prompts inside CI, spinners in redirected logs, swallowed subprocess output, and failures that recommend deleting state.

**Direction for Aruo:** every mutating operation should follow inspect -> resolve -> validate -> plan -> approve -> apply -> verify. Interactive mode is a front end to that pipeline, not separate logic.

### Open-source governance

Repository community files are necessary but insufficient. Kubernetes demonstrates governance as an operating model: persistent ownership groups, temporary cross-cutting working groups, explicit committees for sensitive responsibilities, documented membership, and visible subproject ownership. Smaller projects need less ceremony, but the same questions remain.

A sustainable repository makes clear:

- who can triage, review, merge, release, and make security decisions;
- how contributors gain responsibility and how inactive maintainers step back;
- where technical decisions, conduct reports, and vulnerabilities go;
- which versions are supported and how end-of-life is announced;
- how conflicts and deadlocks are resolved;
- how stewardship survives a founder's departure;
- who owns each subsystem, template, and external integration.

Recurring weakness: projects generate GOVERNANCE, CODEOWNERS, and SECURITY files before real people and response capacity exist. This creates promises without operations.

**Direction for Aruo:** governance guidance must scale by project stage and require named roles, review cadence, and succession signals. `doctor` should distinguish document presence from operational evidence and avoid scoring volunteer responsiveness as if it were build output.

### Plugin systems

Successful systems expose different lessons:

- ESLint offers named contribution types—rules, processors, parsers, formatters, and shared configuration—making extension boundaries discoverable.
- Vite uses lifecycle hooks and compatible Rollup concepts; ordering and environment differences make hook semantics powerful but complex.
- Terraform runs plugins as separate executables behind a versioned Protobuf/gRPC protocol. Major protocol versions define compatibility and minor versions are additive.
- VS Code uses extension hosts, lazy activation, declared contribution points, and placement rules to protect startup and support local, web, and remote contexts.

The core trade-off is authority versus interoperability:

| Model | Portability | Isolation | Startup | API evolution |
|---|---:|---:|---:|---:|
| In-process Go module | Low | Low | Best | Coupled to Go ABI/source API |
| Child process + typed RPC | High | Medium | Moderate | Strong version negotiation |
| WASI component | High | Potentially high | Moderate | Promising, ecosystem still evolving |
| Shell hook | High superficially | Very low | Variable | No meaningful contract |

Process separation prevents memory corruption and version-linking conflicts; it does not itself create a sandbox. A child process still inherits filesystem, environment, network, and credentials unless the host mediates them.

**Direction for Aruo:** begin with narrow declarative contribution points. Future executable plugins should negotiate protocol versions, declare capabilities, receive bounded inputs, return typed findings or proposed operations, and never directly own terminal formatting or unrestricted repository writes. Distribution needs identity, checksums, signatures, compatibility metadata, and a lock. A marketplace should follow—not precede—the trust model.

### Configuration systems

Configuration complexity is dominated by precedence and provenance, not serialization syntax. YAML is approachable and supports comments but has edge cases and lossy round trips. TOML is predictable for moderate configuration. JSON is universal and schema-friendly but poor for hand-authored comments. Executable configuration is flexible but undermines deterministic analysis and sandboxing. CUE demonstrates powerful constraint composition and order-independent validation, but asking every user to learn a constraint language would raise adoption cost.

Aruo needs a typed internal model regardless of surface syntax. The external system should provide:

1. built-in defaults;
2. organization policy;
3. user configuration;
4. repository configuration;
5. named profile;
6. explicit environment overrides;
7. command-line flags.

Precedence alone is insufficient. `config explain <path>` should show the effective value, source file or environment variable, overridden candidates, default status, sensitivity, and schema/deprecation information.

Other required properties:

- reject unknown keys by default with spelling suggestions;
- schema-version the document and provide explicit migrations;
- preserve comments/order where migrations can do so safely;
- distinguish absent, inherited, explicit zero, and explicit disable;
- keep secrets as references, not values stored in project configuration;
- resolve paths relative to the file that declares them unless explicitly documented otherwise;
- allow plugins only namespaced typed configuration;
- produce redacted diagnostics and deterministic serialization.

**Direction for Aruo:** use a familiar declarative format backed by a strict Go model and published machine schema. Consider CUE internally or as an advanced policy interchange only after proving that it materially improves composition without becoming a user prerequisite.

### Cross-platform filesystem APIs

Repository tools operate on adversarial inputs: template paths, archives, symlinks, case collisions, reserved names, interrupted writes, and concurrently modified trees.

Go's `io/fs` provides a host-independent, slash-separated, UTF-8 path contract and `testing/fstest` support. `filepath.Localize` converts valid portable paths to OS paths. Since Go 1.24, `os.Root` provides traversal-resistant operations under a directory; by Go 1.26 it includes reading, writing, recursive creation/removal, links, symlinks, and rename operations. This is substantially safer than `Join` followed by a lexical prefix check and avoids common symlink TOCTOU errors on supported platforms.

Portability hazards still include:

- Windows reserved device names, drive-relative paths, path length, file locks, and rename behavior;
- case-insensitive and case-preserving filesystems;
- Unicode normalization differences;
- symlink availability and privilege differences;
- executable bits and POSIX permissions that Git records differently from host APIs;
- line endings and generated binary detection;
- non-atomic replacement across filesystem boundaries;
- antivirus/indexer interference and delayed deletion;
- network filesystems with weaker atomicity or consistency.

**Direction for Aruo:** represent template paths using `io/fs` semantics, validate collisions under case-folding, stage changes in a sibling directory, fsync/close before replacement where durability matters, and use traversal-resistant rooted operations. Model file mode and symlink intent explicitly. Test released binaries on Windows, macOS, and Linux—mocks cannot establish portability.

## Patterns from successful projects

| Project/tool family | Durable lesson | Caution |
|---|---|---|
| Go toolchain | Cohesive commands, compatibility promise, standard testing and formatting | Aruo cannot assume Go's control over an ecosystem |
| Cargo | One obvious project lifecycle and reproducible metadata | Native ownership is language-specific |
| uv / Ruff | Exceptional speed plus consolidation removes coordination overhead | Fast replacement still needs compatibility discipline |
| GitHub CLI | Task-oriented command tree, browser/API integration, structured output | Service coupling and authentication add ambient state |
| Docker CLI | Stable object/action vocabulary and context-aware operation | Context selection can make effects surprising |
| Terraform | Plan-before-apply and versioned out-of-process plugins | Provider trust, state, and protocol complexity are substantial |
| Vite | Minimal fast default with explicit framework variants | Template/framework combinations multiply quickly |
| Deno / Bun | Integrated distribution and fast feedback | Owning runtime semantics is outside Aruo's role |
| Kubernetes | Explicit ownership and scalable governance | Process can overwhelm small communities |
| Docusaurus / VitePress / MkDocs / Fumadocs | Content as source, strong static delivery, extensible authoring | Theme selection does not solve information architecture |

Across these examples, the durable traits are clear boundaries, native conventions, deterministic state, progressive disclosure, first-class diagnostics, and explicit compatibility. Their failures tend to occur where context, authority, or ownership is implicit.

## Recurring ecosystem pain points

1. **Day-zero/day-two discontinuity.** Creation tools do not retain enough intent to maintain what they created.
2. **Fragmented policy.** The same support or security requirement is separately encoded in docs, workflows, bot config, and release scripts.
3. **Opaque provenance.** Users can see an effective setting or changed file but not which source and rationale produced it.
4. **Unsafe extensibility.** Templates and plugins commonly inherit full user authority.
5. **Semantic drift.** Text templates cannot understand user-owned edits, renamed concepts, or evolving native conventions.
6. **Automation without trust design.** CI often has more credentials and mutability than its job requires.
7. **Human-only interfaces.** Prose output is scraped by automation; machine output lacks remediation and stable schemas.
8. **False quality signals.** File-presence checklists and aggregate scores reward appearance rather than operability.
9. **Dependency noise.** Update bots create volume without encoding compatibility, cooldown, or ownership policy.
10. **Documentation drift.** Examples, command help, config reference, and supported versions diverge.
11. **Unreliable performance claims.** Warm-cache microbenchmarks are presented as end-user experience.
12. **Cross-platform afterthoughts.** Windows, case collisions, symlinks, permissions, and interrupted writes are discovered after release.
13. **Governance theater.** Generated policies outpace actual maintainers and response capacity.
14. **SemVer ambiguity.** Projects version a package but leave CLI, config, output, and protocol compatibility undefined.

## Opportunities for Aruo

### 1. A repository intent model

Represent project kind, ecosystems, lifecycle policy, managed capabilities, ownership, support windows, and exceptions as typed intent—not just template variables. Rendered files become projections of this model.

### 2. Lifecycle continuity

Use the same model for `create`, future `adopt`, `doctor`, `upgrade`, migration, and release preparation. Preserve blueprint/plugin versions and the evidence used to select them.

### 3. Explainable provenance

For any setting, finding, or proposed file change, answer: what produced this, which policy requires it, what evidence was observed, who owns it, and how can it be changed safely?

### 4. Plans as a stable intermediate representation

All mutating features should emit a deterministic typed plan containing preconditions, operations, conflicts, external effects, permissions, rollback limits, and verification. Interactive approval and CI policy consume the same plan.

### 5. Semantic upgrades

Classify files as Aruo-owned, user-owned, shared/structured, or unmanaged. Prefer format-aware edits and native commands. Never overwrite a shared file merely because a template changed.

### 6. Native-tool orchestration

Discover capabilities and versions, show exact commands, capture results, and preserve native layouts. Aruo should make Cargo, Go, uv, npm-compatible managers, Git, and GitHub easier to use together, not hide them.

### 7. One evidence model, several outputs

Audit checks return typed observations, findings, confidence, severity, ownership, remediation, and evidence locations. Render terminal, JSON, SARIF, and documentation from that model.

### 8. Capability-based extensibility

Plugins contribute schemas, detectors, template components, policies, or proposed operations. The host mediates filesystem, process, network, secrets, and output access. Trust is visible and lockable.

### 9. Repository qualification

Validate not only presence but function: install from the documented path, run examples, build release archives, verify attestations, follow a fresh-contributor workflow, and test migrations against representative repositories.

### 10. Policy with context and exceptions

Standards should be parameterized by project kind and maturity. A small library, CLI, web service, and research repository should not receive identical CI, governance, documentation, or benchmark requirements. Exceptions need rationale and expiry, not silent disablement.

### 11. Teaching repositories

Generated artifacts should explain their purpose near the point of use, expose next steps, and link to maintainable guidance. Aruo should improve team understanding rather than create magical infrastructure only Aruo can operate.

### 12. Fleet evolution, later

Once single-repository provenance and upgrades are trustworthy, organization-wide policy diffing, rollout waves, compatibility dashboards, and migration simulation become defensible. Building fleet control before local correctness would amplify errors.

## What Aruo should not become

- another compiler, package resolver, formatter, or language build tool;
- an unrestricted remote-template and shell-hook marketplace;
- a mandatory hosted service for operating an open-source repository;
- an opaque AI agent that mutates repositories without a deterministic plan;
- a universal folder standard that erases ecosystem conventions;
- a quality score optimized for badges or file presence;
- a configuration megafile mirroring every underlying tool option;
- a plugin API exposing internal Go packages as its long-term compatibility contract;
- a generator that owns user code after creation;
- a release button that combines version choice, build, credentials, and publication invisibly.

## Provisional architectural constraints

These are research-derived constraints to test and formalize separately. They are not accepted ADRs.

### MUST

- Keep domain and application behavior independent of the CLI framework.
- Calculate and expose a deterministic plan before consequential writes or external effects.
- Preserve artifact, configuration, policy, and plugin provenance.
- Use ecosystem-native tools as authorities for native dependency/build behavior.
- Define stable, schema-versioned machine output separately from human presentation.
- Treat template paths, archives, plugins, CI, and release inputs as trust boundaries.
- Support non-interactive operation without hidden prompts.
- Test released artifacts on tier-one operating systems.
- Define compatibility for CLI, config, output, templates, and plugin protocol before 1.0.
- Work without a required Aruo cloud service.

### SHOULD

- Prefer declarative extension points before executable plugins.
- Use strict typed configuration and expose value provenance.
- Separate release preparation, artifact construction, and publication.
- Make cache use observable, bounded, and disposable.
- Preserve user-owned edits and fail clearly on ambiguous ownership.
- Produce actionable audit findings with evidence and confidence rather than a single score alone.
- Keep startup work proportional to the command requested.
- Generate docs/reference from canonical command and schema models where feasible.

### MAY, after evidence

- Use CUE internally for policy composition.
- Adopt an out-of-process RPC or WASI plugin runtime.
- provide centrally distributed reusable workflows;
- add organization fleet management;
- use hosted services for optional indexing, collaboration, or telemetry.

## Questions that require ADRs, RFCs, or experiments

1. What is the canonical repository intent schema, and which parts are user-owned?
2. YAML, TOML, or another surface format: which best survives comments, migrations, and schema tooling?
3. Can structured editing achieve acceptable conflict-free upgrade rates on realistically modified repositories?
4. Which plugin contributions require executable code, and what cross-platform capability enforcement is actually possible?
5. Should the first plugin protocol be subprocess RPC, WASI, or deferred entirely?
6. What cold-start, memory, scan, and binary-size budgets represent meaningful user experience?
7. Which audit findings predict project health without penalizing project size, age, or volunteer capacity?
8. What support contract applies to blueprints and generated CI after an Aruo release reaches end of life?
9. How are organization policy and repository autonomy reconciled, explained, and overridden?
10. Which release-intent model works for a single Go binary today without foreclosing multi-artifact repositories?

## Recommended investigation sequence

1. Build a versioned intent and plan specification on paper; test it against Go CLI, Python library, TypeScript library, and web application examples.
2. Assemble a hostile filesystem corpus covering traversal, links, case collisions, Unicode, Windows names, permission failures, and interrupted writes.
3. Measure native baseline workflows and define performance budgets before optimizing implementation.
4. Prototype provenance-aware adoption and a three-version template upgrade; measure conflicts on edited fixtures.
5. Define audit evidence and output schemas; validate terminal, JSON, and SARIF projections from identical fixtures.
6. Threat-model templates, plugins, GitHub Actions, package installation, and release publication.
7. Run contributor onboarding studies with fresh environments on Windows, macOS, and Linux.
8. Only then propose executable plugin and organization-scale architecture.

## Source register

Primary sources reviewed for this snapshot:

- Go: [Go 1.26 release notes](https://go.dev/doc/go1.26), [release policy and history](https://go.dev/doc/devel/release), [modules reference](https://go.dev/ref/mod), [`testing` package](https://pkg.go.dev/testing), [fuzzing](https://go.dev/doc/security/fuzz/), [`io/fs`](https://pkg.go.dev/io/fs), and [traversal-resistant `os.Root`](https://go.dev/blog/osroot).
- CLI and creation: [Cobra](https://github.com/spf13/cobra), [GitHub CLI](https://github.com/cli/cli), [Vite scaffolding](https://vite.dev/guide/), [`uv init`](https://docs.astral.sh/uv/concepts/projects/init/), [uv project layout](https://docs.astral.sh/uv/concepts/projects/layout/), [Copier configuration](https://copier.readthedocs.io/en/latest/configuring/), and [Cookiecutter hooks](https://cookiecutter.readthedocs.io/en/stable/advanced/hooks.html).
- Documentation: [Docusaurus versioning](https://docusaurus.io/docs/versioning), [VitePress search](https://vitepress.dev/reference/default-theme-search), [MkDocs configuration](https://www.mkdocs.org/user-guide/configuration/), [MkDocs Material search](https://squidfunk.github.io/mkdocs-material/plugins/search/), and [Fumadocs OpenAPI](https://www.fumadocs.dev/docs/integrations/openapi).
- Releases: [SemVer 2.0.0](https://semver.org/), [release-please](https://github.com/googleapis/release-please), [Changesets](https://github.com/changesets/changesets), [GoReleaser](https://goreleaser.com/), [GoReleaser SBOMs](https://www.goreleaser.com/customization/sbom/), and [GoReleaser signing](https://www.goreleaser.com/customization/sign/binary_sign/).
- GitHub delivery and security: [Actions security](https://docs.github.com/en/actions/concepts/security), [reusable workflows](https://docs.github.com/en/actions/concepts/workflows-and-actions/reusing-workflow-configurations), [artifact attestations](https://docs.github.com/en/actions/concepts/security/artifact-attestations), [dependency review](https://docs.github.com/en/code-security/concepts/supply-chain-security/dependency-review), [Dependabot version updates](https://docs.github.com/en/code-security/concepts/supply-chain-security/dependabot-version-updates), and [OIDC](https://docs.github.com/en/actions/reference/security/oidc).
- Plugins and configuration: [Terraform plugin protocol](https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol), [Vite plugin API](https://vite.dev/guide/api-plugin.html), [ESLint extension model](https://eslint.org/docs/latest/extend/), [VS Code extension hosts](https://code.visualstudio.com/api/advanced-topics/extension-host), and [CUE configuration concepts](https://cuelang.org/docs/concept/how-cue-enables-configuration/).
- Governance: [Kubernetes community and governance](https://github.com/kubernetes/community).

## Conclusion

The ecosystem does not need Aruo to replace its best tools. It needs Aruo to make their combined lifecycle coherent, inspectable, secure, and maintainable.

The durable product thesis is: **preserve engineering intent and turn it into reviewable plans, native operations, and verifiable evidence throughout a repository's life.** If Aruo can do that while respecting ecosystem ownership and user edits, it can solve a problem that generators, linters, CI templates, and release bots each address only in fragments.
