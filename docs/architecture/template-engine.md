# Template Engine Architecture

## Purpose

The template engine converts a validated blueprint bundle and explicit project data into a deterministic, caller-owned in-memory file plan. It does not discover templates, resolve trust, execute hooks, merge structured documents, or write a repository. Those responsibilities belong to artifact, blueprint, plan, semantic-edit, and apply layers.

```text
embedded / local / plugin bundle
       │ fs.FS + typed blueprint
       ▼
validate sources, destinations, metadata, values, limits
       ▼
parse restricted text/template programs
       ▼
render or copy → normalize ordering → immutable file plan
       │
       ▼
future repository planner (conflicts, diff, approval, atomic apply)
```

This pure boundary is essential: generation must be previewable and testable without touching a working tree.

## Technology evaluation

### `text/template`

Go's standard `text/template` supports substitution, `if`/`with`/`range`, named templates, strict missing-map-key errors, custom function maps, and concurrent execution after construction. It adds no runtime dependency and is familiar across Go tooling. It does not sandbox arbitrary objects, cancel a template mid-execution, or understand repository semantics.

Aruo selects it with `missingkey=error`, a small deterministic function allowlist, JSON-like variable validation, source/output byte limits, and trusted bundle resolution outside the engine. Template data exposes plain metadata and values, never filesystem, network, environment, clock, random source, process runner, secrets, or plugin handles.

### `embed`

The standard `embed.FS` implements `fs.FS`, is immutable and concurrency-safe, and fails the build when declared patterns do not resolve. It is the source for first-party templates shipped in the binary. The engine depends only on `fs.FS`, so development fixtures, organization bundles, and future plugin bundles use the same renderer without special cases.

Embedding increases binary size and couples first-party template updates to an Aruo release. External versioned bundles remain necessary for faster catalog evolution; their resolver, signatures, digests, licenses, and locks are outside this engine.

### Third-party template libraries

| Option | Benefit | Why it is not the core engine |
| --- | --- | --- |
| Sprig | Large familiar function catalog | Includes nondeterministic, environment, network, crypto, and broad conversion behavior; capability review becomes difficult |
| Pongo2/Gonja | Jinja/Django-like syntax familiar beyond Go | Adds a second language, dependency, compatibility surface, and sandbox review without improving Aruo's core repository model |
| Liquid | Deliberately constrained, portable syntax | Weaker alignment with Go contributors and still adds an external language/runtime contract |
| gomplate | Rich data sources and production templating features | Network, environment, secret, and datasource access conflict with deterministic offline planning |
| fast/compiled HTML engines | High throughput for request rendering | Repository generation is latency-insensitive compared with correctness; these engines optimize the wrong workload |

The engine does not import a general function library. A helper is accepted only if it is deterministic, bounded, broadly required by repository templates, documented, and independently tested.

### Filesystem abstraction

Template sources use the standard `io/fs` interfaces and slash-separated valid paths. `embed.FS`, `fstest.MapFS`, `os.DirFS`, archive filesystems, and verified plugin bundle adapters can all satisfy this read boundary. Aruo does not adopt Afero: its broader mutable filesystem API is unnecessary for pure rendering and would duplicate standard-library contracts.

`os.DirFS` is not a security boundary because symlinks can escape its root. The future artifact resolver must open untrusted directories through a containment-safe mechanism, verify the resulting file inventory and digest, and then provide the engine a bounded filesystem.

## Data model

A bundle has two pieces:

- a typed `Blueprint` with stable ID, language, and file specifications;
- an `fs.FS` containing its source files.

Each file specification declares source path, destination path template, whether content is rendered or copied verbatim, and portable permission bits. Project data has first-class name, module, description, author, license, and language fields plus validated JSON-like variables for blueprint-specific choices.

Destinations may substitute metadata but must resolve to unique, relative, slash-separated `fs.ValidPath` values. Empty paths, absolute paths, dot traversal, backslashes, and collisions fail before a result is returned. Results are sorted by destination path regardless of manifest order.

## Rendering contract

- Undefined map keys fail; no `<no value>` output is permitted.
- Conditional content uses standard `if`, `with`, and `range` actions.
- Binary or syntax-conflicting content is copied with rendering disabled.
- Source and rendered file sizes are bounded.
- The caller supplies all values; time, randomness, environment, and network are unavailable.
- Cancellation is checked before and between files. Templates themselves must be trusted and bounded because `text/template` cannot interrupt an executing action.
- The engine returns no partial plan after an error; successful plan bytes belong to the caller.
- File order and bytes are deterministic for identical bundle and data inputs.

Structured JSON, YAML, TOML, XML, workflow, and language manifests should use semantic operations when merging or upgrading existing repositories. Text templates are appropriate for new bounded artifacts; they are not an excuse for whole-file ownership forever.

## Language-specific templates

Language specialization is expressed by composable bundles, not conditionals spanning every ecosystem. The initial embedded example is a Go library bundle. Future language bundles share foundation fragments through blueprint composition at a higher layer; the renderer sees the resolved list of files and one language-consistent metadata object.

A blueprint language must match project metadata. Multi-language repositories will be represented by a higher-level composition containing multiple language components rather than weakening that invariant.

## Embedded catalog

First-party sources live below `internal/templateengine/builtin/templates/` and are compiled into the executable. `builtin.GoLibrary()` returns the embedded filesystem and typed blueprint; callers do not depend on physical source paths. Every built-in bundle requires render tests and native ecosystem qualification before becoming a supported catalog item.

## Plugin extension point

The renderer's extension boundary is data, not in-process code. A future out-of-process plugin may contribute a manifest and file bundle through the versioned plugin protocol. Core resolves identity and permissions, validates the inventory and digest, converts the manifest to `Blueprint`, and renders it using the same engine.

Plugins cannot register Go template functions, hand core callable values, write destinations, or execute hooks during rendering. A proposed helper or semantic operation requires a versioned core/protocol capability. This preserves determinism and language-independent plugins.

## Error and observability model

Errors identify the blueprint, source or destination, and stage (`validate`, `read`, `parse`, or `execute`) while wrapping the cause. Template content and variable values are not copied into errors because future inputs may be sensitive even though secrets are prohibited. The engine emits no logs or terminal output; the application layer decides how to present a failed plan.

## Testing and performance

Unit tests cover substitution, conditions, missing values, raw copies, modes, path traversal, collisions, unsafe values, limits, cancellation, deterministic ordering, embedded assets, and concurrent use. Built-in bundles receive golden/native validation at their owning layer.

Benchmarks report whole-bundle render time, allocations, and bytes. Performance is guarded against regression, but clarity and safety outrank micro-optimizations. Parsed-template caching will be considered only after profiling a realistic multi-layer catalog and defining cache invalidation for mutable external filesystems.

## Extension checklist

Before adding a template function, metadata field, source type, or rendering feature:

1. Prove it cannot be represented as explicit input, blueprint composition, raw copy, or semantic operation.
2. Define determinism, resource limits, error behavior, and cross-platform path semantics.
3. Confirm it exposes no ambient authority or secret value.
4. Add unit, concurrency, malicious-input, and benchmark coverage.
5. Record compatibility impact; template syntax and helper behavior are versioned artifact contracts.
