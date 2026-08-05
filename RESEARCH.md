# Research archive

Research is a permanent engineering artifact, not disposable pre-work. It records ecosystem evidence, benchmark methodology, technology evaluations, rejected options, and open questions separately from accepted ADRs.

## Findings

- generators excel at day zero but usually lose provenance and upgrade intent;
- native language initializers provide the least surprising project spine;
- mature CLIs separate human/machine output, document context precedence, and distinguish interactive authentication from CI credentials;
- repository maturity comes from ownership, executable examples, layered tests, compatibility, reproducible releases, and security—not a universal folder tree;
- plugin ecosystems gain reach while expanding workstation supply-chain risk;
- docs/release/lint tools solve individual stages, leaving lifecycle policy fragmented.

See [`research/ecosystem/competitive-analysis.md`](research/ecosystem/competitive-analysis.md) and the indexed archive in [`research/README.md`](research/README.md).

The current ecosystem snapshot is [`research/ecosystem/2026-developer-tooling-state.md`](research/ecosystem/2026-developer-tooling-state.md). It is research input rather than an accepted architecture decision.

## Lessons applied

Aruo delegates native work, owns a versioned lifecycle model, uses deterministic plans and granular provenance, isolates plugins, separates release from publish, and preserves language-native layouts. Go was accepted after a Rust-first research recommendation; both the historical evaluation and decision rationale remain available.

## Open questions

- Can semantic upgrades exceed a 90% conflict-free target on realistically edited repositories?
- Which policy observations remain comparable across ecosystems without erasing nuance?
- What is the smallest useful plugin permission model on all tier-1 platforms?
- Can YAML comments/order be preserved reliably through migrations?
- Which health indicators predict maintainer outcomes rather than vanity activity?

Future investigations begin with a question, decision relevance, method, sources, limitations, date, owner, and archival path. Findings that change architecture produce an ADR or RFC.
