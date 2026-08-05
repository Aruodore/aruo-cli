# Testing strategy

Testing follows risk rather than a universal coverage number. The purpose is confidence in plans, mutations, compatibility, recovery, and user contracts.

## Layers

- **Unit:** pure config resolution, policy evaluation, graph ordering, path rules, redaction, and type invariants.
- **Integration:** filesystem/Git/process adapters in temporary repositories; native tool integration with pinned fixtures.
- **Golden:** help, diagnostics, JSON/SARIF schemas, plans, diffs, and generated repositories. Updates are explicit and reviewed.
- **Property/fuzz:** parsers, path normalization, merge/reconciliation, graph cycles, idempotency, and arbitrary Unicode/config inputs.
- **Conformance:** external plugins, policy packs, artifact manifests, and machine-output consumers.
- **End-to-end:** released binary in clean Linux/macOS/Windows environments, offline and non-interactive modes.
- **Fault/recovery:** cancellation, disk-full/permission failure, killed process, stale precondition, conflicting edit, and network interruption.
- **Security:** malicious archives/templates/plugins, traversal/symlink escape, secret redaction, signature failure, command injection, and untrusted workspace.

## Coverage expectations

Coverage is reported for trend and changed-code review, not used alone as quality proof. Critical planner, transaction, trust, and schema code requires branch/negative-path evidence. Mutation testing is used selectively for policy and planner logic. Any untested high-risk path needs an owned exception.

## Fixtures and utilities

Fixtures are minimal, licensed, immutable, and describe the behavior they represent. Test builders expose domain concepts instead of filesystem trivia. Network tests use recorded or local deterministic services; no default test depends on external availability. Clocks, randomness, environment, terminal properties, and process runners are injected.

## CI shape

Fast deterministic checks gate every PR. Platform and native-tool matrices run in parallel. Expensive fuzzing, race, benchmark, and ecosystem compatibility suites run on schedule and before release. Flaky tests are defects: quarantine requires an issue, owner, expiry, and maintained signal.

