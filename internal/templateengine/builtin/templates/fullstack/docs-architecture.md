# Architecture

The application is a modular monolith: one repository, one server boundary, one PostgreSQL database. React and Vue use a small Hono server because Vite is a frontend build tool; Next and Nuxt use their native server runtimes.

Data flows from request to validation to a small use case to the database and back through a stable response contract. Keep route adapters thin. Extract a service only when rules are reused or a route can no longer be understood locally.

No business entities are included. The first product feature owns its schema, authorization policy, validation, and tests. Authentication is required before private data; the identity method and verification channel are intentionally selected by the application, not guessed by this starter.
