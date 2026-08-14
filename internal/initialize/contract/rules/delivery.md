# Delivery and operations

These rules apply to build, CI, release, or deployment capabilities that the repository contains.

- Executables with production configuration MUST reject missing or invalid required values. Development defaults MUST NOT silently authorize or configure production behavior.
- Release artifacts MUST be traceable to committed source and reproducible using the ecosystem's declared dependency-resolution mechanism.
- Containerized applications MUST run with the least privileges practical and SHOULD run as a non-root user.
- Deployment guidance MUST cover only applicable concerns, such as schema changes, health signals, shutdown, rollback or roll-forward, and secret injection.
- CI credentials and permissions MUST be minimal. Third-party automation MUST use immutable revisions or a documented trusted update and integrity mechanism.
- `aruo.yaml` is application-owned. Update only capability entries materially affected by the task, preserve unknown fields, and use `SOLVED` only with repository-verifiable or explicitly supplied runtime evidence. Other statuses require a concise reason.
- The contract is a baseline, not certification. Operators retain responsibility for runtime security, infrastructure, recovery, monitoring, incident response, and compliance unless evidence assigns and verifies that responsibility elsewhere.
