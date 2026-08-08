# Changelog

All notable changes to Aruo will be documented in this file.

The format follows [Keep a Changelog 2.0.0](https://keepachangelog.com/en/2.0.0/), and releases follow [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) across the documented CLI, configuration, plan/finding schemas, plugin protocol, and supported public Go APIs.

## [Unreleased]

### Added

- Phase 0 repository design, governance, architecture, research archive, ADR system, RFC process, and documentation structure.
- Go 1.26 module bootstrap, Make task facade, lint/test configuration, cross-platform GitHub Actions, Dependabot, Conventional Commit validation, development container, and draft release automation with SBOMs and provenance attestations.
- `aruo create`'s second catalog entry, `js-library`: a dependency-free JavaScript library template, surfacing the interactive template picker for the first time now that more than one entry is registered.
- `aruo create`'s third catalog entry, `python-library`: a dependency-free Python library template (stdlib `unittest`, `src/` layout).
- `aruo create`'s fourth catalog entry, `ts-library`: a TypeScript library template with strict type-checking. Unlike the others, its `npm install` needs the network for the real `typescript` compiler; Aruo's own test suite verifies its file plan only, and its generated CI does the real `npm ci`/`npm test` run.
- `aruo create`'s fifth catalog entry and first `kind: app` template, `react-app`: React 19 + Vite 8 + Vitest, strict TypeScript, requires Node >=22.22.2 (jsdom's real current engine floor). Same network-dependent-install tradeoff as `ts-library`.
- `aruo create`'s sixth catalog entry, `nuxt-app`: Nuxt 4 with server-side rendering (Nitro), tested with `@nuxt/test-utils`/Vitest/happy-dom. Same network-dependent-install tradeoff as `ts-library`/`react-app`.
- `aruo create`'s seventh catalog entry, `vue-library`: a Vue 3 component library in Vite library mode (dual ESM/UMD output, generated `.d.ts` via vite-plugin-dts, `vue` as a peer dependency). Same network-dependent-install tradeoff as the other npm-based entries.
- `aruo create`'s eighth catalog entry, `next-app`: Next.js 16 (App Router, Turbopack), tested with Vitest/@testing-library/react/jsdom. Pins TypeScript 5.9.3 rather than the 6.x/7.x the Vite-based entries use, matching what `create-next-app` itself currently pins. Same network-dependent-install tradeoff as the other app entries.

### Changed

- Accepted Go for the core CLI, superseding the earlier Rust-first research recommendation.

## Release process

Release intent is captured with the change and reviewed in a release PR. At release, maintainers move entries from `Unreleased` to a dated version, link migrations for breaking changes, build once from the protected tag, verify/sign artifacts, publish through trusted identities, and create curated GitHub release notes from this canonical record. See [the release process](docs/development/release.md).

[Unreleased]: https://github.com/aruodore/aruo-cli/commits/main
