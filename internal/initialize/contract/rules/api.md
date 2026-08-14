# Interfaces and compatibility

These rules apply to externally consumed interfaces, including network APIs, command interfaces, file formats, events, and public library APIs.

- Before releasing a new or incompatible interface, MUST define its consumers, trust boundary, compatibility expectations, input/output contract, and failure behavior.
- Authorization, idempotency, retry, pagination, and concurrency semantics MUST be defined when the interface exposes those concerns.
- MUST use stable machine-consumable error identifiers where callers need programmatic recovery.
- MUST preserve compatibility or provide an explicit migration path when changing a published interface.
- Interface tests SHOULD cover the success and failure cases that actually exist, including authorization and retry behavior when applicable.
