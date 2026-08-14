# Data and migrations

These rules apply when the repository stores or transforms durable data.

- MUST enforce ownership and authorization at a trusted boundary using trusted identity and the selected resource; never trust client-supplied ownership claims.
- MUST define invariants at the narrowest reliable layer. Use storage constraints and transaction boundaries when the storage system owns multi-write invariants.
- When a deployed or shared schema changes, MUST use the repository's migration mechanism and preserve compatibility with supported application versions. Do not rewrite a migration known to have been applied unless the migration system explicitly makes that safe.
- Production or shared-data migrations MUST consider locking, backfill cost, partial failure, roll-forward or rollback, and mixed-version deployment.
- If loss of non-reconstructible data would harm users or operations, MUST record recovery ownership, retention, protection, and a restoration procedure; restoration evidence SHOULD be tested proportionate to the risk.
- Tests MUST isolate mutable test data and SHOULD cover constraints, concurrency, and failure behavior relevant to the change.
