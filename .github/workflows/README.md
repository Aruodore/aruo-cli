# Workflows

- `ci.yml` — module hygiene, vet, cross-platform tests, race detection, and coverage.
- `lint.yml` — pinned golangci-lint policy and format drift.
- `dependency-review.yml` — vulnerable/incompatible dependency changes.
- `pr-title.yml` — Conventional Commit pull request titles without checking out untrusted code.
- `release.yml` — release PR, draft GitHub Release, cross-platform archives, SBOMs, and attestations.

Each workflow defines purpose, triggers, least privilege, trusted/untrusted input behavior, concurrency, caches/artifacts, and timeouts. Third-party actions are pinned to full commit SHAs and updated through reviewed Dependabot PRs.
