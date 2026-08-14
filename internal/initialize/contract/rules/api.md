# APIs, validation, and errors

- Define each API's caller, trust boundary, authorization rule, input/output schema, idempotency needs, and failure model before exposing it.
- Parse and validate input once at the boundary. Internal code receives validated types rather than raw request data.
- Use stable status codes and machine-readable error identifiers. Do not leak stack traces, SQL details, secrets, or provider responses.
- Propagate request or correlation identifiers through logs and downstream work.
- Test success, validation failure, unauthenticated, unauthorized, dependency failure, and retry/idempotency behavior where applicable.
