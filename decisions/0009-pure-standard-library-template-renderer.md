# ADR-0009: Pure standard-library template renderer

- Status: Accepted
- Date: 2026-08-04
- Owners: Aruo maintainers

## Context

Aruo needs embedded, conditional, language-specific repository templates while preserving deterministic plans, plugin isolation, cross-platform behavior, and user ownership. A renderer that writes directly, reads ambient state, or accepts arbitrary helper code would violate those boundaries.

## Decision

Use `text/template` with strict missing keys for rendered text, `embed.FS` for built-in assets, and `io/fs` for all source bundles. Rendering is pure: it validates explicit metadata and JSON-like variables, enforces path and byte limits, and returns a sorted in-memory file plan. It performs no destination writes, discovery, network access, environment reads, time/random access, or hooks.

Do not adopt Sprig, gomplate, Afero, or another template language in core. Future plugins contribute versioned file bundles and manifests through the out-of-process protocol; they cannot register in-process functions.

## Consequences

- Embedded, local, test, and plugin sources share one small interface.
- Plans are deterministic and easy to test before repository mutation.
- The helper vocabulary is intentionally small and may require core review for additions.
- Artifact ingestion must separately enforce containment, signatures, digests, and trust.
- Structured-file evolution requires semantic adapters rather than renderer expansion.
- Large catalogs may eventually need a measured parsed-template cache.

## Alternatives considered

- Broad Go-template functions: convenient but expose nondeterminism and a large security/compatibility surface.
- Jinja/Liquid-like engine: approachable syntax but a second runtime contract and dependency.
- Mutable virtual filesystem: simplifies direct generation but mixes rendering with effects and conflict policy.
- Executable templates: maximally flexible but incompatible with safe, explainable planning.
