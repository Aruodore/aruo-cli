# AI development contract

Read this file before changing code. The application is a modular monolith: Nuxt handles HTTP and rendering; PostgreSQL owns durable application data.

## Boundaries

- `app/` contains pages and browser UI. It never imports `server/` code or reads secrets.
- `server/api/` contains thin HTTP handlers: validate input, call application/database code, translate results.
- `server/utils/` contains small cross-cutting server concerns. Do not create generic utility dumping grounds.
- `server/db/` owns schema, database construction, and committed SQL migrations.
- `shared/` is only for contracts safe in both browser and server bundles.
- Data flows `request -> validation -> handler/use case -> database -> response`. Add a service layer only when logic is reused or the handler stops being readable.

## Required behavior

- Validate all untrusted input at the boundary. Reject unknown object fields unless the API intentionally permits them.
- Return the stable `{ error: { code, message, requestId } }` shape. Never expose stack traces or database errors.
- Treat authentication and authorization separately. This starter has neither: do not add user or private data until both are implemented and tested server-side. Client middleware is never an authorization boundary.
- Read server configuration through `useServerEnvironment`. Secrets never enter `runtimeConfig.public`, logs, source, fixtures, or error responses.
- Use structured logging and request IDs. Never log authorization headers, cookies, tokens, passwords, connection strings, or request bodies.
- Change the Drizzle schema, run `npm run db:generate`, inspect the SQL, and commit both schema and migration. Never use schema push in production or edit an applied migration.
- Dependencies require a concrete reason, active maintenance, acceptable license, lockfile update after the initial install, and `npm audit` review. Prefer platform and existing-library capabilities.
- Keep operations bounded. Add timeouts, limits, cleanup, and failure behavior for network, file, database, email, or job work.

## Tests and completion

- Put fast behavior tests in `tests/`. Use real PostgreSQL integration tests when query behavior matters; do not mock SQL into false confidence.
- Cover success, invalid input, unauthorized/forbidden behavior when applicable, and dependency failure.
- Before declaring work complete, run `npm run check`, review the diff, confirm docs and `.env.example` match behavior, and report any check not run.
- Never disable validation, security controls, lint rules, or tests to make a change pass. Never invent production guarantees. Document an unresolved limitation instead.
