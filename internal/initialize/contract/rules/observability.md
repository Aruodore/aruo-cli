# Errors and observability

- Fail explicitly and preserve error causes internally while returning stable, non-sensitive public errors.
- Emit structured logs with level, event name, timestamp, request/correlation identifier, and useful bounded context.
- Never log credentials, tokens, session identifiers, raw personal data, or unbounded request bodies.
- Define liveness as process health and readiness as dependency readiness. Keep detailed dependency failures out of public responses.
- Add metrics/traces only for a stated operational question. Document cardinality, cost, privacy, retention, and provider ownership.
- Background work must define retry limits, backoff, idempotency, poison-message handling, and observable terminal failure.
