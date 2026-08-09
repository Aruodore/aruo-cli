# Changelog

All notable changes to Aruo will be documented in this file.

The format follows [Keep a Changelog 2.0.0](https://keepachangelog.com/en/2.0.0/), and releases follow [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) across the documented CLI, configuration, plan/finding schemas, plugin protocol, and supported public Go APIs.

## 0.1.0 (2026-08-09)


### ⚠ BREAKING CHANGES

* the Go module path changed from github.com/aruodore/aruo to github.com/aruodore/aruo-cli.

### Features

* **catalog:** add js-library template ([6d82bd2](https://github.com/Aruodore/aruo-cli/commit/6d82bd2f8a3b82748b5c4875a2530060cacbee0e))
* **catalog:** add next-app template ([d543c63](https://github.com/Aruodore/aruo-cli/commit/d543c63652f170352763db6322c1778048c18e14))
* **catalog:** add nuxt-app template ([7fe1b8a](https://github.com/Aruodore/aruo-cli/commit/7fe1b8a9f84072ee68c6c225bb7a69f7d90025ef))
* **catalog:** add python-library template ([84aa7bb](https://github.com/Aruodore/aruo-cli/commit/84aa7bb69c5172f1b4f96ca26dbdc042d1ae6b7c))
* **catalog:** add react-app template ([c584332](https://github.com/Aruodore/aruo-cli/commit/c58433227265ef8cc3a0375fccbddb06b7c03530))
* **catalog:** add ts-library template ([d12d682](https://github.com/Aruodore/aruo-cli/commit/d12d6829d91235ad2824541be534bb8eebc7dcaf))
* **catalog:** add vue-library template ([7d65c92](https://github.com/Aruodore/aruo-cli/commit/7d65c92553b1db31f90a93b898d79d209557c5a3))
* **catalog:** group the template picker by kind ([8c333c1](https://github.com/Aruodore/aruo-cli/commit/8c333c18ff8bcba82992be6ada158f7716b4f918))
* **cli:** add Cobra shell completion ([7821aab](https://github.com/Aruodore/aruo-cli/commit/7821aabb135ef7d7ccb83f19f697d51ef008ebda))
* **cli:** allow --description and --author to be left blank ([7c558b9](https://github.com/Aruodore/aruo-cli/commit/7c558b9a99bfa8e7832d064f0d848cd4133d4c3f))
* **cli:** ask app vs. library before showing the full template list ([e24be71](https://github.com/Aruodore/aruo-cli/commit/e24be7127d9f3b9e218c3947cb0f599315f83639))
* **cli:** integrate lifecycle.Manager into cmd/aruo ([1332974](https://github.com/Aruodore/aruo-cli/commit/13329740ba6801184de8bdfadcfef1182be3c1a0))
* **cli:** integrate terminal session at the composition root ([8711d2a](https://github.com/Aruodore/aruo-cli/commit/8711d2aa68ac17444dbfe1d7249f4ac38f89fe2a))
* **create:** add Nuxt UI to the nuxt-app template ([8f50881](https://github.com/Aruodore/aruo-cli/commit/8f50881dd47d3d9d60f9e498e28559605b67b56c))
* **create:** add vue-app template ([13d1a45](https://github.com/Aruodore/aruo-cli/commit/13d1a45c94fc5d21bd2b4c92dbbb4c66651cf52a))
* **create:** remove the module/package-name prompt entirely ([d7bcd56](https://github.com/Aruodore/aruo-cli/commit/d7bcd56041ecad5cb20ad38c1653abee023ff6eb))
* rename module to github.com/aruodore/aruo-cli ([70167fa](https://github.com/Aruodore/aruo-cli/commit/70167fa673180e3c1f0603666f62962c3fd8d796))
* **tux:** add accessible prompt adapters ([d422741](https://github.com/Aruodore/aruo-cli/commit/d4227414e5347e569554d08ea3c25003531906f9))
* **tux:** add terminal interaction contracts ([e0a684e](https://github.com/Aruodore/aruo-cli/commit/e0a684e223722e02e4a9c35ace93b774c7667b93))
* **tux:** detect terminal capabilities ([1c69a10](https://github.com/Aruodore/aruo-cli/commit/1c69a10987d59a376c81f62ff08024bc95db4191))
* **tux:** let the user go back across create's guided screens ([1fef29c](https://github.com/Aruodore/aruo-cli/commit/1fef29c1e71a5668e79e8e880f9c735bec1904d9))
* **tux:** manage signal lifecycle ([3dccb50](https://github.com/Aruodore/aruo-cli/commit/3dccb50e346206ad0e359e7b07f1daacbf440e37))
* **tux:** render live task progress ([1b05ac6](https://github.com/Aruodore/aruo-cli/commit/1b05ac65c91b1b7672e34e530f8a1fbcff8ed87a))
* **tux:** render semantic terminal output ([7d9e4d4](https://github.com/Aruodore/aruo-cli/commit/7d9e4d4c4f4b898795ef14e071f51b1472c9442d))


### Bug Fixes

* **ci:** fix Windows GOCACHE path and remove a dead session-format param ([1a06959](https://github.com/Aruodore/aruo-cli/commit/1a06959c0c483acb818ad51f3bed252cb7587f49))
* **cli:** drop the long description line from the template picker ([a8b0664](https://github.com/Aruodore/aruo-cli/commit/a8b0664b2d068fffd49ef4152680deb5a38245e3))
* **cli:** script the template picker answer in TestRunCreateInteractive ([1c3389a](https://github.com/Aruodore/aruo-cli/commit/1c3389adecf1c5cc1de69796c378fbb037c03803))
* **cli:** stop calling --module a Go module path for every template ([96f1ee0](https://github.com/Aruodore/aruo-cli/commit/96f1ee0b62ff0d4b0130f7a6d7d6ecfdd95a8650))
* **cli:** stop repeating the template description on the module screen ([d5b7873](https://github.com/Aruodore/aruo-cli/commit/d5b7873422269942d3a24820b40758e10a30fc9a))
* **cli:** stop shadowing err in runCreate ([5cfd729](https://github.com/Aruodore/aruo-cli/commit/5cfd72971f51c717a0c5fc3cf4d37902f001f13d))
* **cli:** stop showing TODO placeholders, derive real defaults instead ([a1c27fa](https://github.com/Aruodore/aruo-cli/commit/a1c27fa6bfeccc893db5be2a6391df1a8b982f54))
* **create:** allow writing into an already-existing empty directory ([25d0180](https://github.com/Aruodore/aruo-cli/commit/25d018071929b9274e9fee6f43641005f3855fa2))
* **doctor:** recognize node --test, unittest, and test_*.py ([5ed4f88](https://github.com/Aruodore/aruo-cli/commit/5ed4f8895202780ab76d48f7ec563488fe19a617))
* **release:** auto-publish releases and flag them as pre-release ([1d45799](https://github.com/Aruodore/aruo-cli/commit/1d45799d3a0d417ed94ef935715180bd137303a4))
* **release:** keep the first release at 0.x, not 1.0.0 ([8590776](https://github.com/Aruodore/aruo-cli/commit/85907767d08417ecf48990bb3160b003386621f7))
* **release:** pin the initial release version to 0.1.0 ([81fe89f](https://github.com/Aruodore/aruo-cli/commit/81fe89fdb1b5f67f1d2da9e2cb74054fa3e70d67))
* **templateengine:** move dependabot.yml and release-please-config.json out of foundation/ ([0c6f5c7](https://github.com/Aruodore/aruo-cli/commit/0c6f5c701169913c24acc4f54e7cb5bc2783e778))
* **tux:** fall back to plain prompter under TERM=dumb ([cda21a1](https://github.com/Aruodore/aruo-cli/commit/cda21a1108dd2b867f934387a5f728f2c51c826b))
* **tux:** guard Progress.Emit against a Close race ([440e6de](https://github.com/Aruodore/aruo-cli/commit/440e6de5a9dfe00bbd54ede9f50338d21c654b4b))
* **tux:** never pre-fill a default into the editable input ([7973486](https://github.com/Aruodore/aruo-cli/commit/7973486024c99d57e9322d8e214f0cd8f62250d4))

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
