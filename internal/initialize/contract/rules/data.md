# Data and migrations

- Application code accesses persistence through an explicit server-side boundary. Never trust client-supplied ownership or authorization fields.
- Commit forward migrations with schema changes. Do not edit an applied migration; add a new one.
- Use database constraints for invariants the database must preserve. Define transaction boundaries around multi-write invariants.
- Review migration locking, backfill cost, rollback/roll-forward strategy, and compatibility during deployments.
- Treat backups as `REQUIRED` once production data matters. Document ownership, retention, encryption, restoration procedure, and evidence from restore drills.
- Tests must isolate data and cover constraints, concurrency-sensitive behavior, and failure paths relevant to the change.
