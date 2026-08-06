# Research archive

Research is a permanent engineering artifact, not disposable pre-work. It records ecosystem evidence, benchmark methodology, technology evaluations, rejected options, and open questions separately from accepted ADRs — immutable evidence for future decisions.

```text
research/
├── ecosystem/       competitor and repository studies
├── architecture/    models, risks, and design explorations
├── product/         vision, UX, roadmap, and implementation studies
├── technology/      language/tool evaluations and benchmarks
├── benchmarks/      investigative benchmark reports
└── future/          framed questions and ideas, not commitments
```

Every study records date, question, method, primary sources, findings, limitations, and decision impact. Accepted conclusions link to ADRs/RFCs. Earlier conclusions are preserved when superseded.

## Current state reports

- [`ecosystem/2026-developer-tooling-state.md`](ecosystem/2026-developer-tooling-state.md) — dated survey of the 2026 tooling landscape and provisional constraints for Aruo. This is research input rather than an accepted architecture decision.
- [`technology/2026-go-terminal-ux-stack.md`](technology/2026-go-terminal-ux-stack.md) — proposed cohesive Go terminal stack, dependency analysis, and vendor-isolation architecture.
- [`ecosystem/competitive-analysis.md`](ecosystem/competitive-analysis.md) — competitor and repository studies behind the findings below.

## Findings

- generators excel at day zero but usually lose provenance and upgrade intent;
- native language initializers provide the least surprising project spine;
- mature CLIs separate human/machine output, document context precedence, and distinguish interactive authentication from CI credentials;
- repository maturity comes from ownership, executable examples, layered tests, compatibility, reproducible releases, and security—not a universal folder tree;
- plugin ecosystems gain reach while expanding workstation supply-chain risk;
- docs/release/lint tools solve individual stages, leaving lifecycle policy fragmented.

## Lessons applied

Aruo delegates native work, owns a versioned lifecycle model, uses deterministic plans and granular provenance, isolates plugins, separates release from publish, and preserves language-native layouts. Go was accepted after a Rust-first research recommendation; both the historical evaluation and decision rationale remain available.

## Open questions

- Can semantic upgrades exceed a 90% conflict-free target on realistically edited repositories?
- Which policy observations remain comparable across ecosystems without erasing nuance?
- What is the smallest useful plugin permission model on all tier-1 platforms?
- Can YAML comments/order be preserved reliably through migrations?
- Which health indicators predict maintainer outcomes rather than vanity activity?

Future investigations begin with a question, decision relevance, method, sources, limitations, date, owner, and archival path. Findings that change architecture produce an ADR or RFC.
