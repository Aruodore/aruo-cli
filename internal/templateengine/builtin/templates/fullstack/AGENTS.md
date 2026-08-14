# AI development contract

This repository is a modular monolith. Read `aruo.yaml` and `docs/architecture.md` before changing it.

- Browser code never reads secrets or imports server modules.
- HTTP handlers validate every untrusted value, reject unexpected fields by default, and return the documented error envelope.
- Authentication proves identity; authorization decides access. Enforce authorization on the server for every protected operation. A hidden screen is not access control.
- Add product tables only after defining ownership, lifecycle, validation, authorization, and deletion behavior.
- Change the Drizzle schema, generate a migration, inspect the SQL, and commit both. Never edit an applied migration or use schema push in production.
- Log structured events with request IDs. Never log secrets, credentials, cookies, authorization headers, full request bodies, or connection strings.
- Bound external work with timeouts, size limits, retries only when safe, and cleanup. Do not claim email, file storage, jobs, caching, analytics, telemetry, or backups are solved without a working provider and tests.
- Add a dependency only when its role is explicit, platform capabilities are insufficient, maintenance and license are acceptable, and the lockfile and audit are reviewed.
- Test success, invalid input, forbidden behavior, and dependency failure where relevant. Use real PostgreSQL for query behavior.
- Before declaring completion run `npm run check` and `npm run check:security`, inspect the diff, and align code, `.env.example`, `aruo.yaml`, and documentation.
- Never bypass validation, authorization, migrations, tests, lint rules, or security controls to make a change pass. Record unresolved concerns honestly.
