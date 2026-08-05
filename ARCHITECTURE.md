# Architecture

## System context

```text
Maintainer / CI / future IDE
            │
          Aruo CLI
            │
   local repository + Git ───── native tools
            │                    (go, uv, cargo, npm)
   optional forge / artifact registries
```

Aruo is a local-first modular monolith. Domain packages know nothing about terminal rendering, GitHub, YAML libraries, or individual package managers. External systems enter through narrow adapters.

## Internal architecture

```text
cmd/aruo
   │
internal/cli ───────────── output (human / JSON / SARIF)
   │
internal/workflow
   │
├── project/config ── artifact resolver/cache/trust
├── inspect ───────── observations
├── policy ────────── findings + evidence
├── plan ──────────── ordered typed operations
├── reconcile ─────── provenance + conflict model
└── execute ───────── transaction / bounded processes
             │
     internal/platform adapters
```

## Core data flow

The executable-specific boundaries, dependency rules, configuration lifecycle, output model, error contract, and test seams are specified in [CLI Application Architecture](docs/architecture/cli-application.md). ADR-0008 records the choice to contain Cobra in a constructor-built presentation layer.

```text
repository + layered config + locked artifacts
                 │
              inspect
                 ▼
            observations
                 │ policy packs
                 ▼
              findings
                 │ workflow intent
                 ▼
        deterministic plan/diff
                 │ approval
                 ▼
       staged apply → verify → evidence
```

Planning does not execute repository code. External effects such as publishing are distinct plan nodes and require explicit targets. File changes stage in repository-local temporary storage, validate preconditions, and commit atomically where supported. Recovery records incomplete work.

## Module responsibilities

- `internal/domain`: versioned project, observation, finding, plan, artifact, and evidence types.
- `internal/config`: layered resolution, validation, migrations, and provenance of values.
- `internal/workspace`: root discovery, ignore semantics, safe filesystem transactions.
- `internal/inspect`: normalized facts from read-only adapters.
- `internal/policy`: pure evaluation and remediation proposals.
- `internal/plan`: dependency ordering, conflicts, risk, previews, rollback metadata.
- `internal/reconcile`: old blueprint/current repository/new blueprint comparison.
- `internal/execute`: cancellation, bounded subprocesses, journaling, and verification.
- `internal/artifact`: resolve, verify, cache, lock, and compatibility.
- `internal/templateengine`: bounded, deterministic `fs.FS` bundle rendering into caller-owned file plans.
- `internal/doctor`: read-only repository observations, versioned checks, scoring, and remediation evidence.
- `internal/plugin`: process protocol, capabilities, permissions, lifecycle.
- `internal/platform`: Git, forge, OS, network, and native-tool adapters.
- `internal/report`: stable machine formats and accessible human presentation.
- `pkg/`: reserved for intentionally supported Go libraries; empty until an API earns stability.

## Plugin architecture

Plugins are child processes speaking version-negotiated JSON Lines over stdin/stdout. Manifests declare publisher, protocol range, artifact digest, requested permissions, and contributed capabilities. Plugins return typed observations/findings/operations; they do not write files or print directly into core output. Repository-declared plugins do not activate in an untrusted workspace.

## Template engine

Blueprints compose foundation, language, workload, capability, and organization layers. Restricted templates render new text; semantic adapters change JSON, YAML, TOML, XML, and language structures. Executable hooks are plugins or explicit process operations, never hidden template behavior. Provenance is recorded at semantic-key or managed-region granularity.

## Configuration

Committed `aruo.yaml` expresses intent; ignored local overrides express developer preference; `.aruo/lock.yaml` pins artifact versions and digests. Precedence is flag → environment → local → project → workspace → organization → default. Enforced policy cannot be overridden silently. Secrets are references only.

## Future expansion

An IDE, local index daemon, or hosted fleet service may call the same workflow contracts. They are added only after measurements justify them. The local CLI and open schemas remain complete without hosted infrastructure.

## Invariants

Deterministic plans; no hidden network or code execution; no secrets in logs/state; user ownership of generated files; explicit external effects; stable machine output; native ecosystem authority; accepted ADRs preserved.
