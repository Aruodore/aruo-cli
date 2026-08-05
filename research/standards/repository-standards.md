# Repository standards

Aruo implements the existing handbook as policy packs. It MUST preserve the native variations in the ecosystem [repository template](../../../REPOSITORY_TEMPLATE.md) and SHALL NOT generate empty ceremonial directories.

## Universal baseline and reason

| Artifact | Decision |
|---|---|
| `README.md` | Promise, audience, status, five-minute result, limits, support; the public entrypoint |
| `LICENSE` | Machine-identifiable legal permission; selected explicitly |
| `CHANGELOG.md` | User-visible history and migration links; generated from reviewed release intent |
| `ROADMAP.md` | Themes with confidence, owners, and non-commitment language |
| `SECURITY.md` | Private reporting, support window, response expectations; avoids public vulnerability leakage |
| `CONTRIBUTING.md` | Exact local loop, review criteria, DCO/CLA policy, architecture links |
| `CODE_OF_CONDUCT.md` | Community behavior and real enforcement contacts |
| governance/maintainers/support | Ownership, decision rights, succession, supported venues |
| `.github/` | Typed issue forms, PR template, CODEOWNERS, least-privilege pinned workflows, dependency updates |
| `docs/` | Diátaxis-oriented versioned source and ADRs |
| `examples/` | Executable supported use cases tested in CI |
| tests/benchmarks | Native test layout plus governed performance evidence |
| scripts/task facade | Short documented entrypoints; reusable logic stays tested in packages |

## Generated GitHub automation

The baseline includes formatting/lint/type checks, unit tests, build/package validation, docs build/link check, secret/dependency review, and release verification. Integration/E2E, coverage, fuzzing, benchmarks, provenance, and CodeQL are enabled based on risk and language support. Actions are pinned by full commit digest with update automation. Permissions default to read-only and elevate per job.

Templates MUST contain contextual TODOs or omit a file; they must never ship fake contacts, security SLAs, benchmark claims, badges, or contributor identities. `aruo check` evaluates substance where feasible: valid license identity, reachable security contact, runnable quick start, workflows triggered on protected paths, and examples exercised.

## Branch and contribution model

Default: protected `main`, short-lived branches, required reviewed PRs, linear or merge history chosen once, signed/tagged releases, no long-lived `develop` branch. Projects needing release branches document support and backport policy. Issue labels are a small controlled taxonomy (`type`, `area`, `status`, `priority`, `good first issue`) with automation that does not overwrite human judgment.

## Conformance levels

- **Foundation:** identity, license, ownership, build/test, security reporting.
- **Maintained:** docs/examples, release automation, dependency policy, CI matrix.
- **Production:** operational/SLO evidence where applicable, compatibility, provenance, incident/migration practice.
- **Research:** dataset/model cards, reproducibility, evaluation and failure analysis.

Lifecycle and conformance are separate: an experimental project can conform to its declared level without pretending to be stable.
