# Tooling decisions

Research review: 2026-08-04. Versions are pinned in `.go-version`, `.tool-versions`, workflows, and release configuration. Dependabot proposes updates; maintainers review release notes and compatibility rather than auto-merging.

## Go 1.26

Go 1.26.5 is the current security-patched stable toolchain. `go.mod` declares language compatibility at 1.26.0 and recommends toolchain 1.26.5. CI reads the module file instead of repeating a drifting version. Go’s supported-release policy covers the two latest major releases; Aruo initially tests the current release and will add oldstable when a public compatibility promise exists.

## Make, not a second task graph

Make is the portable discoverability facade for local and CI-equivalent tasks. It delegates to Go and pinned tools and contains no business logic. Adding Task, Just, Mage, or another runner now would duplicate names and maintenance without solving a demonstrated portability problem.

## Formatting and linting

`gofmt` remains canonical. `goimports` organizes imports. golangci-lint 2.11.4 provides one curated, schema-versioned policy and parallel execution. The standard set is supplemented with correctness, security, error, spelling, and suppression hygiene checks. “Enable all” is deliberately rejected because unstable/noisy linters create churn and suppressions rather than confidence.

Formatting is deterministic and automatic through `make fmt`; CI uses `fmt-check`. `go vet` runs separately so its baseline remains visible even if the aggregator changes.

## Tests and security

`go test` is the native runner. Pull requests test Linux, macOS, and Windows; Linux also runs the race detector and atomic coverage. Coverage is evidence for review, not a universal pass percentage. GitHub dependency review blocks newly introduced moderate-or-higher known vulnerabilities and explicitly reviewed incompatible licenses.

`govulncheck` will be added when implementation dependencies exist; adding it before a meaningful call graph would create configuration ceremony without signal.

## Conventional commits and SemVer

Pull request titles follow Conventional Commits and become squash commit messages. This provides readable intent and release-please input without forcing every local work-in-progress commit through a Node-based commit hook. SemVer covers documented CLI behavior, exit codes, configuration, schemas, plugin protocol, artifacts, and supported public Go APIs.

## Release automation

release-please 5 maintains a reviewable release PR and canonical changelog from Conventional Commit history. GoReleaser 2.17.1 builds cross-platform binary archives and checksums from the protected tag; pinned Syft 1.44.0 produces SBOMs. GitHub attestations bind artifacts to workflow/repository/commit provenance. Publishing is protected by the `release` environment. GoReleaser source archives are disabled so internal design material is not copied into Aruo-managed release packages.

## Git hooks

No mandatory pre-commit framework is configured. The Python-based framework is capable and cross-platform, but it would introduce a second runtime solely to call checks already provided by Make and CI. Hooks are bypassable and must never be the enforcement boundary. Editors may run formatting on save; contributors run `make check`; CI is authoritative. Reconsider a Go-native hook manager only if measured contributor feedback shows the task facade is insufficient.

## GitHub Actions security

Workflows default to no/read-only permissions, elevate per job, avoid `pull_request_target` checkout, set timeouts/concurrency, disable checkout credential persistence, and pin actions to full commit SHAs. Dependabot keeps SHA pins current through reviewed PRs. Release credentials use a protected environment and a fine-grained token/GitHub App token when configured.

## Primary sources

- [Go release history and support policy](https://go.dev/doc/devel/release)
- [Go tool dependencies](https://go.dev/doc/modules/managing-dependencies#tools)
- [golangci-lint configuration](https://golangci-lint.run/docs/configuration/file/)
- [golangci-lint official action](https://github.com/golangci/golangci-lint-action)
- [GitHub Actions secure-use guidance](https://docs.github.com/en/actions/reference/security/secure-use)
- [GitHub artifact attestations](https://docs.github.com/en/actions/concepts/security/artifact-attestations)
- [Dependabot configuration reference](https://docs.github.com/en/code-security/reference/supply-chain-security/dependabot-options-reference)
- [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
- [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
- [release-please action](https://github.com/googleapis/release-please-action)
- [GoReleaser documentation](https://goreleaser.com/)
- [Syft SBOM tooling](https://github.com/anchore/syft)
