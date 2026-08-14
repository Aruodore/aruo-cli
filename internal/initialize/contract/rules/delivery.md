# Delivery and operations

- Production configuration must fail fast when required values are missing or invalid. Development defaults must not silently become production defaults.
- Build artifacts must be reproducible from committed sources and lockfiles. Run applications as a non-root user where containerized.
- Deployment guidance must state migrations, health checks, graceful shutdown, rollback/roll-forward strategy, and secret injection.
- Keep CI permissions minimal and pin third-party automation to immutable revisions where practical.
- Update `aruo.yaml` using `SOLVED`, `REQUIRED`, `OPTIONAL`, `DEFERRED`, or `UNKNOWN`. A `SOLVED` claim needs evidence; every other status needs a reason.
- The starter establishes a baseline, not certification. Runtime security, infrastructure, backups, monitoring, incident response, and compliance remain operator responsibilities unless proven otherwise.
