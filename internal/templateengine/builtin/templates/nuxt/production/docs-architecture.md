# Architecture

This is a modular monolith. One deployable process owns browser rendering and HTTP APIs; PostgreSQL owns durable state. This keeps transactions, deployment, and local reasoning simple until the product demonstrates a need for another process.

Requests enter a `server/api` handler, are validated, perform a bounded use case, and return data or the documented error envelope. Handlers may call the database directly for a small single-purpose operation. Extract a service when business rules are reused or cannot be understood in the handler.

The starter does not invent product tables. Add the first domain concept only after its ownership, lifecycle, authorization policy, validation rules, and deletion behavior are understood.

Authentication and authorization are not configured yet. Before storing user or private data, choose an identity mechanism, establish a server-side session, define ownership or role policies, enforce those policies inside every protected server operation, and add negative integration tests. A hidden page or client route guard is not authorization.
