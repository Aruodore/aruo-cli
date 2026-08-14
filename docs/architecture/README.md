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

This diagram is the target shape of the `internal/workflow` pipeline, none of which is implemented yet; see [Module responsibilities](#module-responsibilities) below for what actually exists today (`internal/cli`, `internal/initialize`, `internal/create`, `internal/doctor`, `internal/templateengine`, `internal/catalog`, `internal/tux`).

```text
cmd/aruo
   │
internal/cli ───────────── output (human / JSON / SARIF)
   │
internal/workflow (planned)
   │
├── project/config ── artifact resolver/cache/trust
├── inspect ───────── observations
├── policy ────────── findings + evidence
├── plan ──────────── ordered typed operations
├── reconcile ─────── provenance + conflict model
└── execute ───────── transaction / bounded processes
             │
     internal/platform adapters (planned)
```

## Core data flow

This is the target flow once `inspect`/`policy`/`plan`/`execute` exist; today only `doctor`'s read-only inspect-and-score path and `create`'s direct template-to-write path are implemented. The executable-specific boundaries, dependency rules, configuration lifecycle, output model, error contract, and test seams are specified in [CLI Application Architecture](cli-application.md). ADR-0008 records the choice to contain Cobra in a constructor-built presentation layer.

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

### Implemented

- `internal/cli`, `internal/cli/command`, `internal/cli/iostreams`, `internal/clierror`: Cobra command wiring, output streams, and error presentation.
- `internal/create`: destination validation and staged, atomic writes for `aruo create`.
- `internal/templateengine`, `internal/templateengine/builtin`: bounded, deterministic `fs.FS` bundle rendering into caller-owned file plans.
- `internal/catalog`, `internal/catalog/builtin`: the embedded template catalog `aruo create` selects from.
- `internal/doctor`: read-only repository observations, versioned checks, scoring, and remediation evidence.
- `internal/initialize`: read-only stack detection, embedded contract planning, managed provenance, and no-overwrite repository adoption.
- `internal/tux` and its `charm`, `plain`, `session`, `policy`, `term`, and `lifecycle` subpackages: terminal capability detection, prompt/progress adapters, and signal lifecycle. See [the terminal UX specification](../cli/terminal-ux.md) for what is wired today versus documented as a gap.
- `internal/buildinfo`: build-time version metadata.
- `pkg/`: reserved for intentionally supported Go libraries; empty until an API earns stability.

### Planned (not yet implemented)

These package names describe the target shape of the `internal/workflow` pipeline in [Internal architecture](#internal-architecture) above. None of them exist in the tree yet; they are recorded here as a design target, not a current module list, and each requires an accepted implementation RFC before it lands.

- `internal/domain`: versioned project, observation, finding, plan, artifact, and evidence types.
- `internal/config`: layered resolution, validation, migrations, and provenance of values.
- `internal/workspace`: root discovery, ignore semantics, safe filesystem transactions.
- `internal/inspect`: normalized facts from read-only adapters.
- `internal/policy`: pure evaluation and remediation proposals.
- `internal/plan`: dependency ordering, conflicts, risk, previews, rollback metadata.
- `internal/reconcile`: old blueprint/current repository/new blueprint comparison.
- `internal/execute`: cancellation, bounded subprocesses, journaling, and verification.
- `internal/artifact`: resolve, verify, cache, lock, and compatibility.
- `internal/plugin`: process protocol, capabilities, permissions, lifecycle.
- `internal/platform`: Git, forge, OS, network, and native-tool adapters.
- `internal/report`: stable machine formats and accessible human presentation.

## Plugin architecture

Not yet implemented; no `internal/plugin` package exists. This section records the design target. Plugins are child processes speaking version-negotiated JSON Lines over stdin/stdout. Manifests declare publisher, protocol range, artifact digest, requested permissions, and contributed capabilities. Plugins return typed observations/findings/operations; they do not write files or print directly into core output. Repository-declared plugins do not activate in an untrusted workspace.

## Template engine

`internal/templateengine` today renders bounded, deterministic `fs.FS` bundles (the built-in catalog entries under `internal/catalog/builtin`) into a caller-owned file plan; layered blueprint composition and semantic JSON/YAML/TOML/XML adapters below are the design target, not current behavior. Blueprints compose foundation, language, workload, capability, and organization layers. Restricted templates render new text; semantic adapters change JSON, YAML, TOML, XML, and language structures. Executable hooks are plugins or explicit process operations, never hidden template behavior. Provenance is recorded at semantic-key or managed-region granularity.

## Configuration

Not yet implemented; no `internal/config` package or `aruo.yaml`/`.aruo/lock.yaml` reader exists. This section records the design target. Committed `aruo.yaml` expresses intent; ignored local overrides express developer preference; `.aruo/lock.yaml` pins artifact versions and digests. Precedence is flag → environment → local → project → workspace → organization → default. Enforced policy cannot be overridden silently. Secrets are references only.

## Future expansion

An IDE, local index daemon, or hosted fleet service may call the same workflow contracts. They are added only after measurements justify them. The local CLI and open schemas remain complete without hosted infrastructure.

## Invariants

Deterministic plans; no hidden network or code execution; no secrets in logs/state; user ownership of generated files; explicit external effects; stable machine output; native ecosystem authority; accepted ADRs preserved.

## Related pages

- [Project structure](project-structure.md)
- [CLI application architecture](cli-application.md)
- [Template engine architecture](template-engine.md)
- [Create command architecture](create-command.md)
- [Repository doctor architecture](doctor.md)
- [Architecture decisions](../../decisions/README.md)
