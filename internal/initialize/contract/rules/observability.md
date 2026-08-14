# Errors and observability

- MUST preserve useful internal error context while exposing only stable, non-sensitive information across trust boundaries.
- MUST NOT log credentials, authentication tokens, session secrets, raw sensitive personal data, or unbounded attacker-controlled content.
- Libraries SHOULD return errors or expose diagnostic hooks rather than emit unsolicited logs.
- Long-running services SHOULD emit structured operational events with severity, event identity, and bounded context. Correlation identifiers SHOULD be propagated when work crosses request or asynchronous boundaries.
- Long-running network services MUST distinguish process liveness from dependency readiness when both signals are provided; public health responses MUST NOT disclose sensitive dependency details.
- Metrics and traces SHOULD be added only for a stated reliability or operational question, with cardinality, cost, privacy, retention, and ownership considered.
- Durable asynchronous work MUST define retry limits, backoff, idempotency where needed, terminal-failure handling, and a way for operators or callers to observe failure.
