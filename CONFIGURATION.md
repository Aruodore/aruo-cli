# Configuration

## Sources and precedence

Highest wins: CLI flag → `ARUO_*` environment override → ignored `aruo.local.yaml` → project `aruo.yaml` → workspace config → organization policy → built-in default. Enforced organization keys cannot be weakened locally. `aruo config explain <key>` reports the resolved value, origin, enforcement, and rationale.

## Files

- `aruo.yaml`: committed human intent—identity, lifecycle, languages, capabilities, policies, blueprint, plugins.
- `aruo.local.yaml`: ignored personal performance/output/tool-path preferences; never required in CI.
- `.aruo/lock.yaml`: committed exact artifact versions, digests, protocol versions, and managed provenance.
- `.aruo/state.json`: ignored derived cache/state; safe to delete.

Every format has `apiVersion` and a published JSON Schema. Unknown keys fail with source spans and suggestions. Comment/order-preserving migrations are performed by `aruo migrate config`; semantic changes never happen silently.

## Defaults and profiles

Defaults are conservative, offline, non-destructive, and ecosystem-native. Named profiles may select coherent development/CI or organization contexts; profiles cannot contain secrets or obscure the resolved configuration. Inheritance is shallow and cycle-checked.

## Environment and secrets

Environment variables are reserved for CI/runtime overrides and credentials, not full nested configuration. Names begin `ARUO_`. Secrets are references to environment, OS keychain, or provider credential aliases. Resolved output always redacts them.

## Plugins

Plugins are explicitly listed with pinned identity/version and granted capabilities. Manifest-requested permissions and repository grants are distinct; permission widening requires approval. See [PLUGINS.md](PLUGINS.md).

