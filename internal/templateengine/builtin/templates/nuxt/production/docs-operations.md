# Operations

The Docker image runs the Nuxt Node server as a non-root user. Configuration arrives at runtime. Migrations are a separate, explicit release operation; the web process does not mutate schema on startup.

Use `/api/health/live` only to decide whether the process is responsive. Use `/api/health/ready` to decide whether it should receive traffic. Readiness deliberately hides dependency details from callers.

Before production, the operator must provide TLS, a trusted reverse proxy, secret storage, database connection limits, automated backups, restore drills, monitoring and alerts, a migration/rollback procedure, and enough replicas/capacity for expected load. The starter provides structured stdout logs but no telemetry exporter or log retention.

The server applies conservative framing, MIME-sniffing, referrer, opener, and cross-domain policy headers. A Content Security Policy is intentionally not guessed: define and test one after the product's script, image, font, frame, and connection origins are known.

Deploy immutable images. Apply backward-compatible migrations before new code where possible. Never rely on rolling application rollback to reverse a destructive database migration.
