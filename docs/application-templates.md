# Application templates

Aruo exposes one focused baseline for each supported application framework.

IDs: `react`, `vue`, `next`, and `nuxt`.

Each template is the smallest coherent framework application: strict types, accessible UI foundations, deterministic tests, linting, formatting, production-build verification, CI, a security policy, the managed engineering contract, application-owned stack guidance, and an intent manifest. Backend and operator concerns are not disguised as solved; `aruo.yaml` records them explicitly.

The templates are foundations, not production certification. Identity, authorization policy, data storage, deployment, recovery, observability, provider configuration, and other product-specific capabilities remain explicit work until repository or runtime evidence proves otherwise.

## Status vocabulary

- `SOLVED`: implemented and backed by inspectable evidence.
- `REQUIRED`: must be implemented or configured before the relevant use.
- `OPTIONAL`: necessary only when the application has that requirement.
- `DEFERRED`: intentionally outside the template's scope.
- `UNKNOWN`: cannot be established from available evidence.

`aruo doctor` verifies supported repository evidence and distinguishes it from runtime or provider claims.
