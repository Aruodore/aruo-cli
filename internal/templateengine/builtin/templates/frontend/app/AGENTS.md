# AI development contract

Read this before changing code. This repository is a browser application with a static production artifact. Its backend contract is external and must be treated as an untrusted network boundary.

## Boundaries and flow

- UI and feature code belongs in `src/`. Keep code close to the feature that owns it.
- Browser-safe configuration uses the framework's public environment prefix. It is public by definition and must never contain secrets.
- Data flows `user input -> validation -> API client -> external service -> validated response -> UI`.
- Centralize transport behavior only when it consistently handles a real concern such as base URL, timeouts, error translation, or credentials. Do not create generic repositories/services for simple local state.

## Required behavior

- Validate untrusted form input and external API responses at their boundaries.
- Model loading, empty, failure, retry, and offline states. Never swallow a rejected promise.
- Authentication UI is not authentication. Authorization must be enforced by the backend on every protected operation.
- Never store long-lived secrets or sensitive tokens in source, public environment variables, local storage, logs, analytics, or error reports.
- Dependencies require a concrete purpose, maintenance and license review, a lockfile update, and an advisory review.
- Preserve semantic HTML, keyboard access, visible focus, reduced motion, and meaningful accessible names.
- Consider cache behavior and rollback whenever deployment changes assets or API compatibility.

## Completion

Add tests for behavior and failure paths, not implementation details. Run `npm run check`, review the diff, confirm documentation and `aruo.yaml` still match reality, and report any check not run. Never disable validation, accessibility rules, security controls, or tests to make work pass. Document unresolved production concerns explicitly.
