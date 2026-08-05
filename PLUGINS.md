# Plugin architecture

## Contract

Third-party plugins run as isolated child processes and exchange versioned JSON Lines messages over stdin/stdout. A handshake negotiates protocol/features. Requests include operation ID, deadline, cancellation, workspace handle, and granted capabilities. Responses contain typed observations, findings, or plan operations—not direct writes or terminal presentation.

## Lifecycle

```text
discover → inspect → install pinned artifact → verify → grant → lock
         → activate → health check → update/review → disable/revoke
```

Plugins activate lazily. Repository-declared plugins stay disabled until the workspace and artifact are trusted. Failure, timeout, invalid messages, oversized output, and cancellation are isolated and reported without crashing the core.

## API and extension points

Plugins may inspect supported facts, evaluate checks, propose remediations, contribute blueprint capabilities, integrate native tools/providers, or add namespaced commands under `aruo x`. The core owns configuration resolution, planning, mutations, output, credential brokerage, and trust decisions.

## Permissions and isolation

Permissions cover repository reads, proposed write globs, named executable invocation with argument constraints, declared network hosts, credential aliases, and forge operations. Installation displays requests; grants are recorded separately; permission widening on update requires approval. Secrets are passed only through scoped handles when possible.

## Versioning and compatibility

The manifest declares plugin version, publisher, core/protocol ranges, executable targets and digests, license, source, permissions, and schemas. Protocol changes follow SemVer and a published support window. Conformance fixtures test handshake, schemas, ordering, conflicts, cancellation, errors, and security boundaries.

## Distribution

v1 supports bundled first-party artifacts, local paths, and explicitly configured pinned sources. A public marketplace is prohibited until publisher verification, signing/provenance, permissions review, scanning, reporting, revocation, and incident response exist.

