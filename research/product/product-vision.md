# Product vision

## North star

In five years, Aruo should be the local-first engineering control plane that can explain what a repository is, whether it is healthy, why its practices differ from policy, and how to evolve it safely. It succeeds when a maintainer can move from idea to trustworthy release—and a newcomer from clone to first contribution—without reconstructing undocumented institutional knowledge.

## Users and jobs

- **Maintainer:** create, release, secure, evolve, and eventually hand over a project.
- **Contributor:** understand architecture, run checks, make a valid change, and see why a gate exists.
- **Platform team:** publish standards once, apply them across languages, inspect exceptions, and roll out migrations safely.
- **Auditor/researcher:** reproduce evidence and trace artifacts to source, inputs, and policy versions.

## Product pillars

1. **Compose:** create or adopt repositories from capabilities, not monolithic templates.
2. **Understand:** inventory code, contracts, tools, ownership, dependencies, risks, and lifecycle state.
3. **Assure:** run risk-layered checks and emit durable, machine-readable evidence.
4. **Evolve:** preview and reconcile standards, templates, dependencies, APIs, and repository migrations.
5. **Ship:** coordinate version intent, release gates, provenance, publishing, docs, and rollback guidance.
6. **Learn:** derive onboarding, architecture maps, explanations, and improvement proposals from verified repository facts.

## Five-year product shape

- Years 0–1: excellent offline CLI; create/adopt/audit; small signed template catalog; deterministic plans.
- Years 1–2: safe upgrades, release orchestration, cross-repository policy packs, plugin SDK.
- Years 2–3: repository graph and organization fleet view; migration campaigns; historical health trends.
- Years 3–5: evidence-grounded engineering assistant proposing architecture/docs/tests and verifying claims; optional hosted collaboration, never mandatory for core workflows.

## Boundaries

Aruo is not a universal build system, package manager, CI provider, Git forge, IDE, or autonomous code author. It does not claim “production-ready” based on file presence. It verifies declared outcomes and reports evidence and uncertainty. AI may propose; deterministic engines plan and validate; humans approve consequential mutations and releases.

## Measures

- median time to first verified project and first external contribution;
- percentage of recommendations with an explanation and reproducible evidence;
- upgrade success/conflict/rollback rates;
- local/CI parity and flaky-check rate;
- release lead time and rollback/migration incidents;
- false-positive and accepted-exception rates;
- retention after 30/180 days, not repository creation count.
