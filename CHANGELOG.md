# Changelog

All notable changes to Aruo will be documented in this file.

The format follows [Keep a Changelog 2.0.0](https://keepachangelog.com/en/2.0.0/), and releases follow [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) across the documented CLI, configuration, plan/finding schemas, plugin protocol, and supported public Go APIs.

## [Unreleased]

### Added

- Phase 0 repository design, governance, architecture, research archive, ADR system, RFC process, and documentation structure.
- Go 1.26 module bootstrap, Make task facade, lint/test configuration, cross-platform GitHub Actions, Dependabot, Conventional Commit validation, development container, and draft release automation with SBOMs and provenance attestations.
- `aruo create`'s second catalog entry, `js-library`: a dependency-free JavaScript library template, surfacing the interactive template picker for the first time now that more than one entry is registered.

### Changed

- Accepted Go for the core CLI, superseding the earlier Rust-first research recommendation.

## Release process

Release intent is captured with the change and reviewed in a release PR. At release, maintainers move entries from `Unreleased` to a dated version, link migrations for breaking changes, build once from the protected tag, verify/sign artifacts, publish through trusted identities, and create curated GitHub release notes from this canonical record. See [the release process](docs/development/release.md).

[Unreleased]: https://github.com/aruodore/aruo/commits/main
