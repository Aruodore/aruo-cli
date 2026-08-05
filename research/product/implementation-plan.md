# Implementation plan

This is sequencing, not authorization to write implementation code. Phase 0 ends with review and explicit approval.

## Workstreams and gates

### 0. Validate architecture

- interview 8–12 maintainers across Go/Python/Rust/TS and audit 20 edited generated repositories;
- prototype only on disposable fixtures: config round-trip, plan determinism, three-way semantic update, process protocol;
- benchmark Rust versus Go decision criteria;
- threat-model repository code, templates, plugins, credentials, external publishing;
- decide naming/trademark and governance.

Gate: ADRs accepted; user jobs and success metrics validated; no unresolved critical trust boundary.

### 1. Establish contracts

- version Project/Observation/Finding/Plan/Artifact/Evidence schemas;
- define command/exit/output/config compatibility policy;
- create conformance fixtures for paths, line endings, symlinks, conflicts, cancellation, and dirty Git;
- extract handbook requirements into traceable policy IDs with rationale/source.

Gate: an independent implementation could consume fixtures and produce conforming output.

### 2. Build read-only vertical slice

- workspace discovery, configuration resolution/explanation, Go and Python inspectors;
- policy evaluation and human/JSON/SARIF reports;
- performance/security/accessibility harnesses and docs.

Gate: `inspect` and `check` run offline, deterministically, cross-platform, without executing repository code by default.

### 3. Add safe mutation

- planner, semantic adapters, staged filesystem transaction, diff, journal/recovery;
- create/adopt/fix for foundation + Go/Python;
- artifact manifest, digest, lock, local cache; upgrade fixtures.

Gate: fault injection proves no silent partial mutation; edited-repository upgrade targets are met.

### 4. Broaden deliberately

- Rust/TS, CLI and library profiles, then service/frontend;
- GitHub, docs, release and publishing adapters;
- signed catalog and compatibility automation.

Gate: each profile meets its native build/test/package/docs/release evidence matrix before “supported.”

### 5. Stabilize v1

- out-of-process plugin protocol/SDK; permissions and trust UX;
- compatibility/migration policy, security review, reproducible releases;
- governance, contributor onboarding, support rotation, public benchmarks.

## Team topology

Assign owners for core/planner, language adapters, artifact supply chain, developer experience/docs, and security/release. Use one repository initially unless independently released SDKs justify separation. Every workstream maintains ADRs, fixtures, benchmarks, and upgrade notes alongside implementation.

## Definition of done

A feature is done when its contract/schema, human and machine UX, threat considerations, cross-platform tests, failure/recovery behavior, docs/examples, performance impact, migration/compatibility story, and owner are present. Generated output itself passes the applicable Aruodore standards.

## First decisions after approval

1. Approve the product boundary and terminology (`blueprint`, `capability`, `policy pack`).
2. Approve the Rust-versus-Go evaluation rubric, not a language by enthusiasm.
3. Select two design-partner repositories and representative historical revisions.
4. Convert handbook clauses into a policy traceability catalog.
5. Publish ADR-0001 and the v0.1 acceptance suite before building CLI behavior.
