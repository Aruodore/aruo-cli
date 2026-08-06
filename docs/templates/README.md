# Templates and blueprints

Planned documentation covers blueprint selection, input reference, composition/compatibility, customization, authoring, qualification, upgrade/conflict handling, and provenance examples. Templates never hide executable hooks.

The implemented rendering boundary is specified in [Template Engine Architecture](../architecture/template-engine.md). It deliberately returns a file plan rather than writing a repository. Discovery, composition, semantic edits, and transactional application remain separate layers.

## Model

A template renders a bounded text artifact. A blueprint composes foundation, language, workload, capability, and organization layers into a repository plan. This distinction prevents one giant conditional template and makes combinations testable.

## Discovery and locking

Artifacts are discovered from built-in, configured local, and approved registry sources in that order. Remote discovery is explicit and disabled offline. Selection resolves compatibility constraints and records exact version, digest, source, license, protocol, and signature in `.aruo/lock.yaml`. Names are globally stable publisher/artifact identifiers.

## Variables

Inputs have type, validation, default, help, sensitivity, and derivation rules. High-consequence values—license, organization, visibility, publish target—require explicit selection in non-interactive use. Derived values remain inspectable. Secret values are prohibited.

## Rendering and customization

Templates use a restricted expression language with strict undefined-variable errors, deterministic functions, and controlled line-ending behavior. They cannot access network, environment, clock, random values, or processes unless a value is supplied explicitly. Structured files use semantic adapters instead of text replacement.

Generated files become user-owned. Aruo records managed semantic keys/regions and prior artifact provenance, not blanket ownership. Customization occurs through blueprint inputs, organization overlays, or normal user edits.

## Hooks

There are no hidden executable hooks. Required commands become visible typed plan operations with declared executable, arguments, working directory, environment allowlist, network need, rollback, and approval. Reusable behavior belongs in permissioned plugins.

## Extension points

- First-party language bundles are embedded below `internal/templateengine/builtin/templates/` and exposed by typed bundle constructors.
- Verified local and registry artifacts will provide the same `fs.FS` plus `Blueprint` contract after resolution and trust checks.
- Out-of-process plugins may contribute manifests and file bundles through the future versioned protocol; they cannot register Go functions or write repositories.
- New deterministic helpers require core review, tests, documentation, and compatibility analysis.
- Raw-copy files handle binaries and content whose syntax conflicts with Go template actions.
- Structured-file creation and evolution belongs to semantic adapters rather than increasingly complex text functions.

See the [architecture extension checklist](../architecture/template-engine.md#extension-checklist) before expanding the renderer.

## Validation

Each blueprint tests minimal/default/maximal supported compositions, determinism, native build/test/package/docs, platform paths/line endings, upgrades from supported versions, edited-file conflicts, and license/security rules. Unsupported combinations fail during planning with explanation and alternatives.
