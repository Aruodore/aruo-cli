# Plugin architecture

## Decision

Use out-of-process plugins over a versioned JSON-RPC-like protocol on stdio. This follows the portability of Git/Docker-style executable extensions while avoiding the supply-chain blast radius of in-process IDE/package plugins. ESLint’s explicit metadata/config/rule/processor surface and Astro/Vite’s named lifecycle hooks show the value of small contracts; VS Code Workspace Trust shows that code execution requires an explicit trust boundary. [ESLint plugins](https://eslint.org/docs/latest/extend/plugins), [Astro integrations](https://docs.astro.build/en/reference/integrations-reference/), [VS Code Workspace Trust](https://code.visualstudio.com/api/extension-guides/workspace-trust).

## Extension points

- discover observations;
- evaluate checks and propose findings;
- contribute blueprint capabilities and semantic operations;
- provide language/package/docs/release adapters;
- add commands only under `aruo x <plugin>` initially;
- render additional report formats.

Plugins do not receive raw credential stores, unrestricted core objects, or implicit network/filesystem/process access.

## Manifest and protocol

Manifest fields: stable ID, version, publisher, protocol range, executable targets/digests, license, capabilities, requested permissions, schemas, homepage/source, signature/provenance. Handshake negotiates protocol and feature versions. Requests carry workspace ID, operation ID, cancellation, deadline, locale, and explicitly granted handles. Responses are typed observations, findings, or plan operations—not direct terminal output.

## Permission model

Permissions are granular: repository read; proposed writes by glob; execute named binaries with argument constraints; network to declared hosts; credential aliases; Git/forge read/write. Install displays permissions; widening on update requires approval. Untrusted repositories cannot auto-activate project-declared plugins. Organization allowlists and offline deny mode are supported.

## Lifecycle

`search → inspect → install --version → verify signature/digest → grant → lock → enable → update/review → revoke`. Core ships first-party plugins statically or as separately signed artifacts but uses the same conformance suite. Crash, timeout, invalid protocol, oversized response, and cancellation are isolated and actionable.

## Ordering and conflict

Plugins declare constraints (`before`, `after`, capability requirements), never numeric global priority. The planner topologically sorts and rejects cycles. Two operations targeting the same semantic key must merge through the owning adapter or produce a visible conflict. Hook names and payloads are versioned; deprecated protocol versions have published support windows.

## Marketplace posture

No public marketplace before signing, provenance, permissions, publisher verification, malware response, reporting, revocation, and reproducible package validation exist. v1 supports local paths and pinned URLs/registries with explicit trust. Popularity is never a trust signal by itself.
