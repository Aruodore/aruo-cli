# System architecture

## Shape

```text
CLI / future API
      │
Workflow application layer ─── Reporter (human/JSON/SARIF)
      │
Repository session + planner + policy evaluator
      │
IR: Project / Observation / Finding / Plan / Evidence
      │
Adapters: filesystem · Git · forge · process · language · docs · release
      │
Artifact resolver ─ templates/policies/plugins ─ trust store/cache/lock
```

The dependency rule points inward: domain and planner know no CLI framework, forge API, package manager, or template syntax. Adapters implement narrow ports. Workflows coordinate but do not embed language policy.

## Modules

| Module | Responsibility | Must not own |
|---|---|---|
| CLI | parsing, completion, interaction, rendering | domain decisions |
| Workspace | root discovery, ignore rules, filesystem transaction | policy |
| Inspector | adapters produce normalized observations | remediation |
| Policy | evaluate facts and explain findings | direct writes |
| Planner | dependency/conflict ordering and preview | terminal UI |
| Reconciler | provenance-aware semantic/text merge | template discovery |
| Executor | bounded processes, cancellation, rollback journal | hidden network |
| Artifact registry | resolve, verify, cache, lock | arbitrary activation |
| Evidence store | local run metadata and export | secrets/raw excessive logs |
| Git/forge | status, diff, metadata, optional GitHub operations | core assumptions about GitHub |
| Language adapters | native commands, layouts, package metadata | cross-cutting policy |
| Plugin host | protocol, permissions, lifecycle, isolation | in-process third-party code |

## Repository state

`.aruo/lock.yaml` records selected artifact versions/digests and managed regions; `.aruo/state.json` is derived local cache and ignored; `.aruo/evidence/` is optional bounded evidence; `aruo.yaml` holds human intent. A repository remains usable if all Aruo-specific files are removed.

## Upgrade model

Each prior operation stores provenance at file/managed-region/semantic-key granularity. Upgrade replays the old blueprint, computes the new blueprint, compares both with the current repository, then proposes semantic operations. Conflicts are first-class plan nodes, never silently resolved. This borrows Copier’s three-way insight while avoiding whole-file ownership.

## Deployment evolution

Start as one binary and several internal packages—a modular monolith. Introduce a local daemon only if indexing or IDE latency proves it necessary. Add an optional hosted control plane only for fleet aggregation, collaboration, signed artifact distribution, and policy administration. The CLI and formats remain independently useful.
