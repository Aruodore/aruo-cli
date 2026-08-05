# Style guide

## Naming and layout

Use short descriptive lowercase Go package names; avoid `util`, `common`, `helpers`, and stuttering names. Export only stable consumer concepts. Files group cohesive behavior, not arbitrary type categories. Commands and config keys use lowercase kebab-case; environment variables use `ARUO_` plus uppercase snake case.

## Formatting and static checks

Go source is formatted with `gofmt`/`goimports` and passes `go vet` plus the repository’s pinned curated golangci-lint policy. Markdown, YAML, JSON, shell, and website files use their pinned formatters. Generated files identify source and generation command and are checked for drift.

## APIs and comments

Public APIs are small, context-aware, deterministic where promised, and accept interfaces at the consumer boundary rather than preemptively. Exported Go identifiers have doc comments beginning with the identifier. Comments explain constraints, invariants, tradeoffs, and surprising behavior—not syntax. TODOs include an issue/RFC reference and owner or removal condition.

## Errors and logging

Errors add operation and subject context while preserving causes with `%w`. Expected findings are values, not errors. Error strings are lowercase and have no trailing punctuation. Do not both log and return the same error. Logs are structured, secret-redacted, sent to stderr, and quiet by default; machine output never mixes with diagnostics.

## Documentation

Use direct, inclusive language, sentence-case headings, defined acronyms, explicit units, and tested commands. Avoid “simple,” “obvious,” unsupported superlatives, and fake examples. Every task page states prerequisites, expected output, failure modes, and cleanup.

## Tests and examples

Name tests for observable behavior. Prefer table tests only when cases genuinely share setup/assertions. Avoid sleeps and global mutable state. Golden updates require reviewable diffs. Examples use safe realistic names and no live credentials.

## Repository changes

Keep PRs cohesive. Public behavior changes include docs, tests, changelog/release intent, compatibility notes, and schema migrations as applicable. Suppressions and policy exceptions include reason, scope, owner, and expiry.

