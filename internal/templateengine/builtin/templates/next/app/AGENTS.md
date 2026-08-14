# AI development contract

This is a Next.js modular monolith. Server Components, Route Handlers, and Server Actions are server trust boundaries; Client Components are not.

- Keep UI in `app/`. Put server-only application logic in `server/` when introduced. Never import server secrets into client code.
- Validate all request, form, environment, webhook, and external API data at the boundary.
- Authentication establishes identity; authorization is a separate server-side decision on every protected operation. Route visibility is not authorization.
- Define stable, safe error responses. Never expose stack traces, provider payloads, SQL errors, secrets, cookies, or tokens.
- Add a database only when the product needs durable data. Then use committed migrations, bounded connections, and real integration tests. Never mutate schema on web startup.
- Model timeouts, retries, idempotency, limits, cleanup, and partial failure for every external dependency.
- Dependencies need a concrete purpose, maintenance/license review, lockfile update, and advisory review.
- Before completion run `npm run check`, inspect the diff, reconcile docs and `aruo.yaml`, and report checks not run. Never disable safety checks to make work pass.
