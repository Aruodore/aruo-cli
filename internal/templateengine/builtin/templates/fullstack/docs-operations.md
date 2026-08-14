# Operations

Build and run the container with runtime configuration supplied by the deployment platform. Run database migrations as a separate, observable release step before shifting traffic. Never run schema push during startup.

Route liveness to `/api/health/live` and readiness to `/api/health/ready`. Do not expose detailed dependency errors in either response. Terminate TLS at a trusted ingress, keep PostgreSQL private, set resource limits, and verify graceful termination within the platform's shutdown window.

The operator owns secret rotation, database backups, restore drills, rollback, migration coordination, capacity, alert routing, incident response, retention, privacy, and regional requirements. A successful container build does not prove any of these.

The committed in-memory limiter is valid for one process. Replace it with a shared atomic implementation before horizontal scaling or using it for authentication, billing, expensive work, or other abuse-sensitive routes.

The dependency audit may report moderate development-only findings for Drizzle Kit's legacy build chain. There is currently no non-breaking upstream resolution. The gate fails on high and critical findings; reassess this exception on every dependency update rather than suppressing audit output.
