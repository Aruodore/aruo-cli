# Application template profiles

Aruo exposes two profiles for every supported application framework.

## Comprehensive profile

IDs: `react`, `vue`, `next`, and `nuxt`.

This profile supplies broadly reusable production infrastructure: validated configuration, server-side request boundaries, PostgreSQL migrations, authentication infrastructure, authorization policy seams, structured errors and logs, abuse controls, health probes, tests, CI, containers, and operational documentation. It contains no fictional business entities. Email, jobs, cache, files, analytics, and telemetry are included only when the template can provide an honest adapter or explicit integration boundary.

The profile is not a production certificate. Provider credentials, domain policy, backup and restore operations, capacity planning, alert routing, privacy decisions, and application-specific authorization remain the operator's responsibility.

## Lean profile

IDs: `react-lean`, `vue-lean`, `next-lean`, and `nuxt-lean`.

This profile is the smallest coherent framework application: strict types, accessible UI baseline, deterministic tests, linting, formatting, production build verification, CI, security policy, an AI-agent contract, and an intent manifest. Backend concerns are not disguised as solved. They are marked `REQUIRED`, `OPTIONAL`, or `DEFERRED` in `aruo.yaml`.

## Status vocabulary

- `SOLVED`: implemented and backed by inspectable evidence in the generated repository.
- `REQUIRED`: the application developer or operator must implement or configure it before the relevant production use.
- `OPTIONAL`: necessary only for products with that requirement.
- `DEFERRED`: intentionally outside this profile's scope.
- `UNKNOWN`: cannot be established from repository evidence alone.

A future `aruo doctor` should verify evidence and behavior, not merely file names.
