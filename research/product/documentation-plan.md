# Documentation plan

## Site decision

Publish `aruo.aruodore.com` from versioned Markdown in the repository using VitePress initially: fast static output, strong accessible defaults, local fuzzy search, code groups, dark mode, and straightforward customization. Keep content framework-neutral so migration is possible. Use immutable release deployments and stable aliases before introducing copied multi-version trees; Docusaurus itself warns that versioning adds build and contributor complexity. [VitePress overview](https://vitepress.dev/guide/what-is-vitepress), [VitePress search](https://vitepress.dev/reference/default-theme-search), [Docusaurus versioning](https://docusaurus.io/docs/versioning).

## Information architecture

```text
Home / Getting started / Installation
Tutorials
How-to: create, adopt, audit, upgrade, release, migrate
Reference: commands, config, schemas, exit codes
Templates / Plugins / Policy packs
Architecture / Design principles
Examples / Benchmarks
Contributing / Governance / Security
Roadmap / FAQ / Changelog
```

This applies the ecosystem [Documentation Standard](../../../DOCUMENTATION_STANDARD.md): tutorials teach, how-tos solve tasks, reference states exact contracts, and concepts explain tradeoffs.

## Experience requirements

- Home gives audience, problem, honest lifecycle status, install, and first verified result above the fold.
- Installation covers binary packages, checksum/signature verification, supported platforms, upgrades, and uninstall.
- Every command page is generated from the same command model used by `--help`, then enriched with examples and failure modes.
- Config/API reference is generated from schemas; examples are validated.
- Code tabs share semantic example IDs so variants do not drift; copy controls are keyboard/screen-reader accessible.
- Search indexes active version by default, marks old results, records privacy-preserving failed-search terms, and works locally for the current site.
- Dark/light/system themes, visible focus, reduced motion, semantic landmarks, descriptive links, and WCAG 2.2 AA are release gates.
- Interactive examples run in a sandbox with no credentials/network by default; static transcript fallback is mandatory.

## Delivery and ownership

PR previews, link/schema/spell/build checks, tested snippets, performance/accessibility budgets, redirects, sitemap/canonical metadata, and release snapshot are automated. Each section has an owner. Quarterly review validates top tasks and search failures. Documentation version support follows product support; unsupported archives remain visible and clearly labeled.
