# Configuration design

## Files and precedence

Use `aruo.yaml` as committed project intent, `aruo.local.yaml` as ignored developer overrides, and an optional organization policy source resolved to `.aruo/lock.yaml`. Precedence, highest first: CLI flag → purpose-built environment variable → local override → project → workspace → organization policy → built-in default. Organization policy may mark keys enforced; lower layers cannot override them. `aruo config explain <path>` shows value, source, enforcement, and rationale.

Secrets are references (`env:`, OS keychain alias, provider credential alias), never values. Environment variables are reserved for runtime/CI overrides rather than encoding the whole configuration.

## Proposed shape

```yaml
schema: https://aruo.dev/schema/config/v1.json
apiVersion: aruo.dev/v1alpha1
project:
  name: example
  kind: library
  lifecycle: incubating
  description: "..."
identity:
  organization: example-org
  owners: [maintainer-handle]
  license: Apache-2.0
repository:
  forge: github
  visibility: public
  branch: main
languages:
  - id: go
    version: ">=1.25 <1.27"
capabilities:
  documentation: { provider: vitepress, versioning: release-snapshots }
  testing: { profile: library }
  ci: { provider: github-actions }
  release: { strategy: release-pr }
  benchmarks: { profile: library }
policies:
  packs: [aruo/open-source@1, aruo/supply-chain@1]
templates:
  blueprint: aruo/go-library@1
plugins: []
exceptions: docs/project/exceptions.md
```

## Why each domain exists

- `project` selects workload/lifecycle obligations; lifecycle is not inferred from age.
- `identity` drives attribution, ownership, legal defaults, and community files.
- `repository` defines forge/branch integration without hard-coding GitHub into core.
- `languages` selects native adapters and compatibility matrices; ranges avoid accidental toolchain drift.
- `capabilities` states desired outcomes; provider choice remains replaceable.
- `policies` identifies versioned standards and allows reproducible audits.
- `templates` records origin for upgrades; the lock records exact digests.
- `plugins` is explicit activation and permission configuration.
- `exceptions` makes deviations reviewable, owned, and expiring.

## Schema evolution

Unknown keys fail by default with source location and suggested migration. Minor additive schema changes remain compatible; semantic changes require `aruo migrate config`. The CLI preserves comments/order when editing through a concrete-syntax YAML library. Resolved configuration can be printed with secrets redacted. Config validation is offline and editors consume the published JSON Schema.
