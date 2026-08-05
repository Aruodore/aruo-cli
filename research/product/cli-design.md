# CLI design

## Hierarchy

```text
aruo create [blueprint] [path]       create a new repository
aruo adopt [path]                    add Aruo model to an existing repository
aruo inspect                         show normalized repository facts
aruo check [policy|capability]       verify and report; replaces overlapping lint/doctor/audit
aruo fix [finding...]                plan safe local remediations
aruo plan <workflow>                 render operations without applying
aruo upgrade [artifacts...]          reconcile blueprint/policy evolution
aruo add <capability>                add docs, release, benchmark, etc.
aruo remove <capability>             remove only proven managed contributions
aruo run <task>                      invoke normalized native task
aruo release plan|prepare|verify     local release orchestration
aruo publish <target>                explicit external package/release effect
aruo migrate plan|apply              named repository/ecosystem migrations
aruo template list|show|verify       artifact discovery and authoring support
aruo plugin list|inspect|install|... extensions and trust
aruo config get|set|explain|validate configuration
aruo auth login|status|logout        future remote/forge auth; absent until needed
aruo completion <shell>
aruo version
```

`doctor` is reserved for the Aruo installation/runtime, while `check` covers repository conformance. `docs`, `benchmark`, `lint`, and `test` are capabilities/tasks (`aruo run docs`, `aruo run benchmark`) rather than an ever-expanding top level. `release` prepares and verifies; `publish` is deliberately separate because it changes external state.

## Global contract

`--cwd`, `--config`, `--profile`, `--format human|json|sarif`, `--color auto|always|never`, `--quiet`, `--verbose`, `--offline`, `--non-interactive`, `--yes`, and `--no-input`. `NO_COLOR`, dumb terminals, redirected output, CI, and accessibility preferences are respected. JSON goes to stdout; logs/progress to stderr.

## Interaction design

- Interactive create asks high-impact questions first, shows recommended defaults and consequences, and ends with a reviewable plan.
- Non-interactive mode requires complete inputs or fails with exact missing keys; never silently chooses identity/license/publish targets.
- Spinners appear only on TTYs for indeterminate work. Determinate multi-step work uses stable `3/8` progress. Plain mode emits line-oriented updates.
- Color and emoji reinforce but never carry meaning. Status uses text/symbol plus color; hyperlinks have visible URLs in non-TTY output.
- Errors use: summary, cause, location/evidence, recovery command, and debug reference. Expected findings do not print stack traces.
- Suggestions are edit-distance constrained and never execute automatically.

## Safety and scriptability

Every mutation previews a diff unless `--yes` is supplied with a deterministic plan. Destructive/external operations require explicit target and support `--dry-run` when the provider permits. Stable exit codes and versioned JSON schemas are part of compatibility. Shell completion is generated from the command model and supports bash, zsh, fish, and PowerShell.
